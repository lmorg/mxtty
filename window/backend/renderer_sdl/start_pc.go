//go:build !darwin
// +build !darwin

package rendersdl

import "github.com/lmorg/mxtty/types"

func (sr *sdlRender) Start(term types.Term) {
	sr.registerHotkey()
	sr.eventLoop(term)
}
