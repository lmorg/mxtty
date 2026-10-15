package tmux

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed bash.sh.tmpl
var bashScriptTemplate string

var bashScript = template.Must(template.New("bash.sh.tmpl").Parse(bashScriptTemplate))

func writeBashScript(commandLine, callbackID string) (string, error) {
	file, err := os.CreateTemp("", "ttyphoon-*.sh")
	if err != nil {
		return "", err
	}
	err = bashScript.Execute(file, struct {
		CommandLine string
		CallbackID  string
	}{commandLine, callbackID})
	_ = file.Close()
	return file.Name(), err
}
