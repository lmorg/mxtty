package tmux

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const _SEPARATOR = `|||`

type cmdDefinitionT struct {
	cmd    string
	fields []cmdFieldT
}

type cmdFieldT struct {
	name   string
	format string
}

func (def *cmdDefinitionT) CmdLine(parameters ...string) []byte {
	var fields []string
	for i := range def.fields {
		fields = append(fields, fmt.Sprintf(`%s:#{%s}`, def.fields[i].name, def.fields[i].format))
	}

	s := fmt.Sprintf("%s %s -F '%s'\n", def.cmd, strings.Join(parameters, " "), strings.Join(fields, _SEPARATOR))

	return []byte(s)
}

func (tmux *Tmux) sendCommand(cmd *cmdDefinitionT, t reflect.Type, parameters ...string) (any, error) {
	resp, err := tmux.SendCommand(cmd.CmdLine(parameters...))
	if err != nil {
		return nil, err
	}

	var slice []any

	for i := range resp.Message {
		v := reflect.New(t)
		err = parseMxttyLine(resp.Message[i], v)
		if err != nil {
			return nil, err
		}
		slice = append(slice, v.Interface())
	}

	return slice, nil
}

func parseMxttyLine(b []byte, v reflect.Value) error {
	fields := strings.Split(string(b), _SEPARATOR)

	for i := range fields {
		values := strings.SplitN(fields[i], ":", 2)

		err := setFieldValue(v, values[0], values[1])
		if err != nil {
			return err
		}
	}

	return nil
}

func setFieldValue(v reflect.Value, name string, value string) error {
	// Ensure that we have a pointer to a struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct")
	}

	// Get the actual struct value
	v = v.Elem()

	// Get the field by name
	field := v.FieldByName(name)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", name)
	}

	// Ensure the field is settable
	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", name)
	}

	// Set the value, ensuring the types match
	switch field.Type().String() {
	case reflect.String.String():
		field.Set(reflect.ValueOf(value))

	case reflect.Int.String():
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("cannot convert field to int: %s", name)
		}
		field.Set(reflect.ValueOf(i))

	case reflect.Bool.String():
		if value == "true" {
			field.Set(reflect.ValueOf(true))
		} else {
			field.Set(reflect.ValueOf(false))
		}
	}

	return nil
}
