package config

import "os/exec"

/*
	Eventually these will be user configurable rather than compiled time
	options.
*/

var (
	SCROLLBACK_HISTORY = 10000 // lines
	DEFAULT_SHELL      = "/bin/sh"
	FONT_NAME          = ""
	FONT_SIZE          = 14
)

func init() {
	shell, err := exec.LookPath("murex")
	if err == nil {
		DEFAULT_SHELL = shell
	}
}
