package shell

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sync"

	"github.com/lmorg/ttyphoon/ai/agent"
	"github.com/lmorg/ttyphoon/ai/agent/aitypes"
	"github.com/lmorg/ttyphoon/ai/agent/sessiondb"
	"github.com/lmorg/ttyphoon/utils/markdownutil"
)

type Bash struct {
	agent   aitypes.Agent
	enabled bool
}

func init() {
	agent.ToolsAdd(&Bash{})
}

//go:embed bash_description.md
var bashDescription string

func (t *Bash) New(agt aitypes.Agent) (aitypes.Tool, error) {
	return &Bash{agent: agt, enabled: true}, nil
}

func (t *Bash) Enabled() bool { return t.enabled }
func (t *Bash) Toggle()       { t.enabled = !t.enabled }
func (t *Bash) Name() string  { return "bash" }
func (t *Bash) Path() string  { return "internal" }
func (t *Bash) Description() string {
	return bashDescription
}
func (t *Bash) DefaultPermissions() aitypes.DefaultPermissions {
	return aitypes.DefaultPermissions{Invocation: "askPermission", Subagents: "deny"}
}

// bashInputT is deliberately a script rather than a command line. The script
// is written to a temporary file so Bash receives the exact multiline input.
type bashInputT struct {
	Script string `json:"script"`
}

func (t *Bash) InputType() reflect.Type { return reflect.TypeOf(bashInputT{}) }
func (t *Bash) StreamsOutput() bool     { return true }

func (t *Bash) CallStructured(ctx context.Context, input any) (string, error) {
	request, ok := input.(*bashInputT)
	if !ok || request == nil {
		return "", fmt.Errorf("bash received an invalid structured input %T", input)
	}
	s, err := t.run(ctx, request.Script)
	if err != nil {
		s += fmt.Sprintf("\n%v", err)
	}
	return s, nil
}

func (t *Bash) CallDescription(ctx context.Context, input any) string {
	request, ok := input.(*bashInputT)
	if !ok || request == nil {
		return "Tool received an invalid structured input"
	}
	return markdownutil.MarkdownCodeBlock(request.Script, "bash")
}

func (t *Bash) run(ctx context.Context, script string) (string, error) {
	dir, err := os.MkdirTemp("", "ttyphoon-bash-")
	if err != nil {
		return "", fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	path := dir + "/script.sh"
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		return "", fmt.Errorf("write temporary Bash script: %w", err)
	}

	command := exec.CommandContext(ctx, "bash", path)
	stream := agent.OpenAITypedStreamBlock(ctx, sessiondb.StreamBlockToolOutput, "")
	defer stream.Close()
	output := &bashOutputWriter{stream: stream}
	command.Stdout = output
	command.Stderr = output
	err = command.Run()
	if err != nil {
		return output.String(), fmt.Errorf("bash exited unsuccessfully: %w", err)
	}
	return output.String(), nil
}

func (t *Bash) ObservationStructured(input any, output string, err error) aitypes.ToolObservation {
	request, ok := input.(*bashInputT)
	if !ok || request == nil {
		return aitypes.ToolObservation{Tool: t.Name(), Status: "error", Error: fmt.Sprintf("invalid structured input %T", input)}
	}
	observation := aitypes.ToolObservation{
		Tool:   t.Name(),
		Status: "ok",
	}
	if err != nil {
		observation.Status = "error"
		observation.Error = err.Error()
	}
	command := aitypes.CommandObservation{
		Command: request.Script,
		Status:  observation.Status,
		Summary: firstNonEmptyLine(output),
	}
	observation.CommandsRun = []aitypes.CommandObservation{command}
	return observation
}

type bashOutputWriter struct {
	mu     sync.Mutex
	output bytes.Buffer
	stream *agent.AIStreamBlockEmitter
}

func (w *bashOutputWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.output.Write(p); err != nil {
		return 0, err
	}
	if w.stream != nil {
		w.stream.Emit(string(p))
	}
	return len(p), nil
}

func (w *bashOutputWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.output.String()
}
