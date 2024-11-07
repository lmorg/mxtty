//go:build darwin
// +build darwin

package rendersdl

import (
	"github.com/lmorg/mxtty/types"
)

func (sr *sdlRender) Start(term types.Term) {
	// registerHotkey seg faults on macOS if it's not run in a different
	// goroutine
	go sr.registerHotkey()
	sr.eventLoop(term)
}
