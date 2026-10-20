package workspace

import (
	"fmt"
	"strings"
)

// EnvsTmux returns the env vars as a tmux CreatePane command flags
func (t *TermT) EnvsTmux() string {
	return EnvsTmux(t.Envs)
}

func EnvsTmux(envMap map[string]string) string {
	var envStr strings.Builder

	for key, value := range envMap {
		envStr.WriteString(fmt.Sprintf("-e %s=%s ", key, QuoteString(value)))
	}

	return envStr.String()
}

func QuoteString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return fmt.Sprintf(`"%s"`, s)
}
