package filetools

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/lmorg/ttyphoon/ai/agent"
	"github.com/lmorg/ttyphoon/ai/agent/aitypes"
	"github.com/lmorg/ttyphoon/types"
	"github.com/lmorg/ttyphoon/utils/markdownutil"
	"golang.org/x/tools/txtar"
)

type ReadFiles struct {
	agent   aitypes.Agent
	enabled bool
}

func init() {
	agent.ToolsAdd(&ReadFiles{})
}

//go:embed read_description.md
var readDescription string

//go:embed read_auth_desc.md
var readFileAuthDescription string

var readFileAuthDescriptionTemplate = template.Must(template.New("read-file-auth-description").Funcs(markdownutil.TemplateFuncMap()).Parse(readFileAuthDescription))

func (f *ReadFiles) New(agent aitypes.Agent) (aitypes.Tool, error) {
	return &ReadFiles{agent: agent, enabled: true}, nil
}

func (t *ReadFiles) Enabled() bool { return t.enabled }
func (t *ReadFiles) Toggle()       { t.enabled = !t.enabled }

func (t *ReadFiles) Name() string { return "readFiles" }
func (t *ReadFiles) Path() string { return "internal" }
func (t *ReadFiles) DefaultPermissions() aitypes.DefaultPermissions {
	return aitypes.DefaultPermissions{Invocation: "alwaysAllow", Subagents: "allow"}
}
func (t *ReadFiles) Description() string {
	return readDescription
}

type readFilesInputT struct {
	Files []string `json:"files"`
}

func (t *ReadFiles) InputType() reflect.Type { return reflect.TypeOf(readFilesInputT{}) }

func (t *ReadFiles) CallStructured(ctx context.Context, input any) (string, error) {
	request, ok := input.(*readFilesInputT)
	if !ok || request == nil {
		return "", fmt.Errorf("readFiles received an invalid structured input %T", input)
	}
	return t.callFiles(ctx, request.Files)
}

func (t *ReadFiles) CallDescription(_ context.Context, input any) string {
	request, ok := input.(*readFilesInputT)
	if !ok || request == nil {
		return "Tool received an invalid structured input"
	}

	var description bytes.Buffer
	if err := readFileAuthDescriptionTemplate.Execute(&description, request); err != nil {
		return fmt.Sprintf("Unable to render description: %v", err)
	}
	return description.String()
}

func (t *ReadFiles) callFiles(_ context.Context, files []string) (response string, err error) {
	var archive txtar.Archive

	for i := range files {
		filename, pathErr := resolveWorkspacePath(t.agent, files[i])
		if pathErr != nil {
			archive.Files = append(archive.Files, txtar.File{
				Name: files[i],
				Data: fmt.Appendf(nil, "!!! Cannot open file: %v", pathErr),
			})
			continue
		}

		t.agent.Renderer().DisplayNotification(types.NOTIFY_INFO, t.agent.ServiceName()+" requesting file: "+filename[len(t.agent.GetMeta().Pwd):])

		var b []byte
		info, err := os.Stat(filename)
		if err != nil {
			b = fmt.Appendf(nil, "!!! Cannot open file: %v", err)

		} else if info.Name()[0] == '.' {
			b = fmt.Appendf(nil, "!!! Cannot open file: %s", "file hidden")

		} else if b, err = os.ReadFile(filename); err != nil {
			b = fmt.Appendf(nil, "!!! Cannot open file: %v", err)
		}

		archive.Files = append(archive.Files, txtar.File{
			Name: files[i],
			Data: b,
		})

		err = nil
	}

	response = string(txtar.Format(&archive))

	return response, nil
}

func (t *ReadFiles) Observation(input, output string, err error) aitypes.ToolObservation {
	var files []string
	if json.Unmarshal([]byte(input), &files) != nil {
		return aitypes.ToolObservation{Tool: t.Name(), Status: "ok", Inputs: []string{input}}
	}
	return t.observation(files, err)
}

func (t *ReadFiles) ObservationStructured(input any, _ string, err error) aitypes.ToolObservation {
	request, ok := input.(*readFilesInputT)
	if !ok || request == nil {
		return aitypes.ToolObservation{Tool: t.Name(), Status: "error", Error: fmt.Sprintf("invalid structured input %T", input)}
	}
	return t.observation(request.Files, err)
}

func (t *ReadFiles) observation(files []string, err error) aitypes.ToolObservation {
	observation := aitypes.ToolObservation{Tool: t.Name(), Status: "ok"}
	if err != nil {
		observation.Status = "error"
		observation.Error = err.Error()
		return observation
	}
	observation.FilesRead = files
	observation.Counts = map[string]int{"files": len(files)}
	return observation
}
