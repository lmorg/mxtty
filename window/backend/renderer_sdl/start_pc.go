//go:build !darwin
// +build !darwin

package rendersdl

func (sr *sdlRender) registerHotkey() {
	return // currently not supported on Linux
	sr._registerHotkey()
}
