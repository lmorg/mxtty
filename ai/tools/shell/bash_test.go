package shell

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/lmorg/ttyphoon/ai/agent/aitypes"
)

func TestBashExposesStructuredInputAndStreamsOutput(t *testing.T) {
	tool := &Bash{}
	var structured aitypes.StructuredTool = tool
	if got := structured.InputType(); got != reflect.TypeOf(bashInputT{}) {
		t.Fatalf("InputType() = %v, want bashInputT", got)
	}
	if !tool.StreamsOutput() {
		t.Fatal("StreamsOutput() = false, want true")
	}
}

func TestBashRunsScriptAndReturnsCombinedOutput(t *testing.T) {
	output, err := (&Bash{}).CallStructured(context.Background(), &bashInputT{
		Script: "printf 'stdout\\n'; printf 'stderr\\n' >&2",
	})
	if err != nil {
		t.Fatalf("CallStructured() error = %v", err)
	}
	if !strings.Contains(output, "stdout\n") || !strings.Contains(output, "stderr\n") {
		t.Fatalf("output = %q, want stdout and stderr", output)
	}
}

func TestBashReturnsRecoverableExitFailureAsOutput(t *testing.T) {
	output, err := (&Bash{}).CallStructured(context.Background(), &bashInputT{
		Script: "printf 'before failure\\n'; exit 7",
	})
	if err != nil {
		t.Fatalf("CallStructured() error = %v, want recoverable failure as output", err)
	}
	if !strings.Contains(output, "before failure\n") || !strings.Contains(output, "bash exited unsuccessfully: exit status 7") {
		t.Fatalf("output = %q, want partial output and model-visible exit failure", output)
	}
}
