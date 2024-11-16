//go:build !darwin
// +build !darwin

package rendersdl

func (sr *sdlRender) registerHotkey() {
	sr._registerHotkey()
}
