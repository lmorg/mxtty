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

func TestBashReturnsOutputAndExitError(t *testing.T) {
	output, err := (&Bash{}).CallStructured(context.Background(), &bashInputT{
		Script: "printf 'before failure\\n'; exit 7",
	})
	if err == nil || !strings.Contains(err.Error(), "exited unsuccessfully") {
		t.Fatalf("CallStructured() error = %v, want Bash exit error", err)
	}
	if output != "before failure\n" {
		t.Fatalf("output = %q, want script output", output)
	}
}
