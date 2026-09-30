package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/eino-contrib/jsonschema"
)

// Providers expect a self-contained object schema. The default reflector emits a
// root $ref into $defs, which hides type/properties from the model, so it guesses
// and sends nested values as escaped JSON strings.
var structuredToolSchemaReflector = &jsonschema.Reflector{
	ExpandedStruct: true,
	DoNotReference: true,
	Anonymous:      true,
}

func structuredToolSchema(t reflect.Type) *jsonschema.Schema {
	generated := structuredToolSchemaReflector.ReflectFromType(t)
	generated.Version = "" // $schema is metadata the provider does not need
	return generated
}

var jsonUnmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

// decodeStructuredToolInput decodes LLM-supplied arguments into input, repairing
// nested values the model encoded as JSON strings and reporting failures in JSON
// terms rather than Go type names.
func decodeStructuredToolInput(argumentsInJSON string, input any) error {
	target := reflect.ValueOf(input)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return errors.New("arguments must be decoded into a JSON object")
	}
	elemType := target.Elem().Type()

	first := reflect.New(elemType)
	err := strictDecodeJSON(argumentsInJSON, first.Interface())
	if err == nil {
		target.Elem().Set(first.Elem())
		return nil
	}

	if repaired, changed := unwrapStringifiedJSON(json.RawMessage(argumentsInJSON), elemType); changed {
		second := reflect.New(elemType)
		if retryErr := strictDecodeJSON(string(repaired), second.Interface()); retryErr == nil {
			target.Elem().Set(second.Elem())
			return nil
		}
	}

	return describeStructuredInputError(err, elemType)
}

func strictDecodeJSON(argumentsInJSON string, input any) error {
	decoder := json.NewDecoder(strings.NewReader(argumentsInJSON))
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(input); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("expected a single JSON object")
		}
		return err
	}
	return nil
}

// unwrapStringifiedJSON walks raw against the expected type and replaces strings
// that contain encoded JSON with the value they encode.
func unwrapStringifiedJSON(raw json.RawMessage, t reflect.Type) (json.RawMessage, bool) {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == nil || t.Kind() == reflect.Interface || reflect.PointerTo(t).Implements(jsonUnmarshalerType) {
		return raw, false
	}

	value := bytes.TrimSpace(raw)
	changed := false
	if len(value) > 0 && value[0] == '"' && t.Kind() != reflect.String {
		var text string
		if json.Unmarshal(value, &text) == nil {
			inner := json.RawMessage(strings.TrimSpace(text))
			if len(inner) > 0 && (inner[0] == '{' || inner[0] == '[') && json.Valid(inner) {
				value = inner
				changed = true
			}
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		var fields map[string]json.RawMessage
		if json.Unmarshal(value, &fields) != nil {
			return value, changed
		}
		nested := false
		for name, field := range fields {
			fieldType, ok := structFieldTypeByJSONName(t, name)
			if !ok {
				continue
			}
			if repaired, fieldChanged := unwrapStringifiedJSON(field, fieldType); fieldChanged {
				fields[name] = repaired
				nested = true
			}
		}
		if !nested {
			return value, changed
		}
		encoded, err := json.Marshal(fields)
		if err != nil {
			return value, changed
		}
		return encoded, true

	case reflect.Slice, reflect.Array:
		var items []json.RawMessage
		if json.Unmarshal(value, &items) != nil {
			return value, changed
		}
		nested := false
		for i, item := range items {
			if repaired, itemChanged := unwrapStringifiedJSON(item, t.Elem()); itemChanged {
				items[i] = repaired
				nested = true
			}
		}
		if !nested {
			return value, changed
		}
		encoded, err := json.Marshal(items)
		if err != nil {
			return value, changed
		}
		return encoded, true

	case reflect.Map:
		var entries map[string]json.RawMessage
		if json.Unmarshal(value, &entries) != nil {
			return value, changed
		}
		nested := false
		for key, entry := range entries {
			if repaired, entryChanged := unwrapStringifiedJSON(entry, t.Elem()); entryChanged {
				entries[key] = repaired
				nested = true
			}
		}
		if !nested {
			return value, changed
		}
		encoded, err := json.Marshal(entries)
		if err != nil {
			return value, changed
		}
		return encoded, true
	}

	return value, changed
}

func structFieldTypeByJSONName(t reflect.Type, name string) (reflect.Type, bool) {
	if t.Kind() != reflect.Struct {
		return nil, false
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if tag == "-" {
			continue
		}
		if field.Anonymous && tag == "" {
			if embedded, ok := structFieldTypeByJSONName(dereference(field.Type), name); ok {
				return embedded, true
			}
			continue
		}
		if tag == "" {
			tag = field.Name
		}
		if strings.EqualFold(tag, name) {
			return field.Type, true
		}
	}
	return nil, false
}

func dereference(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// describeStructuredInputError rewrites encoding/json failures so they never
// reference Go types the model cannot know about.
func describeStructuredInputError(err error, t reflect.Type) error {
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		expected := jsonTypeName(typeErr.Type)
		if typeErr.Field == "" {
			return fmt.Errorf("arguments must be a JSON %s but received a %s", expected, typeErr.Value)
		}
		return fmt.Errorf("field %q must be a JSON %s but received a %s; send JSON values directly, do not encode them as a string",
			typeErr.Field, expected, typeErr.Value)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("arguments are not valid JSON at byte %d: %s", syntaxErr.Offset, syntaxErr.Error())
	}

	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return errors.New("arguments are not valid JSON: the JSON object is incomplete")
	}

	if name, ok := unknownFieldName(err); ok {
		allowed := jsonFieldNames(t)
		if len(allowed) == 0 {
			return fmt.Errorf("unknown field %q", name)
		}
		return fmt.Errorf("unknown field %q; allowed fields: %s", name, strings.Join(allowed, ", "))
	}

	return fmt.Errorf("arguments must be a JSON object matching the tool schema: %s",
		strings.TrimPrefix(err.Error(), "json: "))
}

func unknownFieldName(err error) (string, bool) {
	const prefix = "json: unknown field "
	message := err.Error()
	if !strings.HasPrefix(message, prefix) {
		return "", false
	}
	var name string
	if json.Unmarshal([]byte(strings.TrimPrefix(message, prefix)), &name) != nil {
		return "", false
	}
	return name, true
}

func jsonFieldNames(t reflect.Type) []string {
	t = dereference(t)
	if t.Kind() != reflect.Struct {
		return nil
	}
	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if tag == "-" {
			continue
		}
		if field.Anonymous && tag == "" {
			names = append(names, jsonFieldNames(field.Type)...)
			continue
		}
		if tag == "" {
			tag = field.Name
		}
		names = append(names, tag)
	}
	sort.Strings(names)
	return names
}

func jsonTypeName(t reflect.Type) string {
	t = dereference(t)
	switch t.Kind() {
	case reflect.Struct, reflect.Map:
		return "object"
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return "string"
		}
		return "array of " + jsonTypeName(t.Elem()) + "s"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "value"
	}
}
