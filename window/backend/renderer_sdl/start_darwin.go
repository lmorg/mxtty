//go:build darwin
// +build darwin

package rendersdl

func (sr *sdlRender) registerHotkey() {
	go sr._registerHotkey()
}
