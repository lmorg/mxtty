package tmux

import (
	"reflect"
	"testing"
)

func TestCommandLine(t *testing.T) {
	expected := `list-windows  -F 'Name:#{window_name}|||Id:#{window_id}|||Width:#{window_width}|||Height:#{window_height}|||Active:#{?window_active,true,false}'`

	var _test_CMD_LIST_WINDOWS = "list-windows"

	type _test_CMD_LIST_WINDOWS_T struct {
		Name       string `tmux:"window_name"`
		Id         string `tmux:"window_id"`
		Width      int    `tmux:"window_width"`
		Height     int    `tmux:"window_height"`
		Active     bool   `tmux:"?window_active,true,false"`
		panes      map[string]*PANE_T
		activePane *PANE_T
	}

	actual := string(mkCmdLine(
		_test_CMD_LIST_WINDOWS,
		reflect.TypeOf(_test_CMD_LIST_WINDOWS_T{}),
	))

	if expected != actual {
		t.Errorf("Incorrect string returned:")
		t.Logf("  Expected: %s", expected)
		t.Logf("  Actual:   %s", actual)
		t.Logf("  exp byte: %v", []byte(expected))
		t.Logf("  act byte: %v", []byte(actual))
	}
}
