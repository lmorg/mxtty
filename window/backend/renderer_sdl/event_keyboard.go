package rendersdl

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventTextInput(evt *sdl.TextInputEvent) {
	switch {
	case sr.inputBox != nil:
		sr.inputBox.eventTextInput(sr, evt)

	case sr.highlighter != nil:
		sr.highlighter.eventTextInput(sr, evt)

	default:
		sr.termWidget.eventTextInput(sr, evt)
	}
}

func (sr *sdlRender) eventKeyPress(evt *sdl.KeyboardEvent) {
	if evt.State != sdl.PRESSED {
		return
	}

	switch {
	case sr.inputBox != nil:
		sr.inputBox.eventKeyPress(sr, evt)

	case sr.highlighter != nil:
		sr.highlighter.eventKeyPress(sr, evt)

	default:
		sr.termWidget.eventKeyPress(sr, evt)
	}
}
