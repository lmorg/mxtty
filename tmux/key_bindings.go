package tmux

import (
	"fmt"
	"strings"

	"github.com/lmorg/mxtty/debug"
)

type keyBindsT struct {
	tmux    map[string]map[string]string
	prefix  string
	fnTable map[string]fnKeyT
}

type fnKeyT func(*Tmux, string) error

var defaultFnKeys = map[string]fnKeyT{
	"Create a new window":  fnKeyNewWindow,
	"Kill the active pane": fnKeyPane,
}

func (tmux *Tmux) _getDefaultTmuxKeyBindings() error {
	const (
		PREFIX = iota
		KEY
		NOTE
	)

	resp, err := tmux.SendCommand([]byte(`list-keys -N -a`))
	if err != nil {
		return err
	}

	tmux.keys.tmux = make(map[string]map[string]string)
	tmux.keys.fnTable = make(map[string]fnKeyT)

	for i := range resp.Message {
		split := strings.SplitN(string(resp.Message[i]), " ", 3)
		if tmux.keys.tmux[split[PREFIX]] == nil {
			tmux.keys.tmux[split[PREFIX]] = make(map[string]string)
		}

		note := strings.TrimSpace(split[NOTE])

		tmux.keys.tmux[split[PREFIX]][split[KEY]] = note

		if fn, ok := defaultFnKeys[note]; ok {
			tmux.keys.fnTable[split[KEY]] = fn
		}
	}

	debug.Log(tmux.keys.tmux)
	debug.Log(fmt.Sprintf("len(tmux.keys.tmux) == %d", len(tmux.keys.tmux)))
	//debug.Log(tmux.keys.fnTable)

	if len(tmux.keys.tmux) == 1 {
		for tmux.keys.prefix = range tmux.keys.tmux {
			// assign key prefix to mxtty
			debug.Log(tmux.keys.prefix)
		}
	}

	return nil
}
