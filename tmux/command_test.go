package tmux

import "testing"

func TestCommandLine(t *testing.T) {
	expected := `list-windows -F 'Name:#{window_name}|||Id:#{window_id}|||Width:#{window_width}|||Height:#{window_height}|||Active:#{?window_active,true,false}'` + "\n"

	def := cmdDefinitionT{
		cmd: "list-windows",
		fields: []cmdFieldT{
			{
				name:   "Name",
				format: "window_name",
			},
			{
				name:   "Id",
				format: "window_id",
			},
			{
				name:   "Width",
				format: "window_width",
			},
			{
				name:   "Height",
				format: "window_height",
			},
			{
				name:   "Active",
				format: "?window_active,true,false",
			},
		},
	}

	actual := string(def.CmdLine())

	if expected != actual {
		t.Errorf("Incorrect string returned:")
		t.Logf("  Expected: %s", expected)
		t.Logf("  Actual:   %s", actual)
		t.Logf("  exp byte: %v", []byte(expected))
		t.Logf("  act byte: %v", []byte(actual))
	}
}
