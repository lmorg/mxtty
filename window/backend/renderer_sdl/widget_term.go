package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

type termWidgetT struct{}

func (tw *termWidgetT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	sr.term.Reply([]byte(evt.GetText()))
}

func (tw *termWidgetT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	sr.keyModifier = evt.Keysym.Mod

	debug.Log(evt.Keysym.Sym)

	switch evt.Keysym.Sym {
	case sdl.K_LSHIFT, sdl.K_RSHIFT, sdl.K_LALT, sdl.K_RALT,
		sdl.K_LCTRL, sdl.K_RCTRL, sdl.K_LGUI, sdl.K_RGUI,
		sdl.K_CAPSLOCK, sdl.K_NUMLOCKCLEAR, sdl.K_SCROLLLOCK, sdl.K_SPACE:
		// modifier keys pressed on their own shouldn't trigger anything
		return
	}

	if evt.Keysym.Sym < 256 && evt.Keysym.Sym > sdl.K_SPACE &&
		(evt.Keysym.Mod == sdl.KMOD_NONE ||
			evt.Keysym.Mod&sdl.KMOD_CAPS != 0 || evt.Keysym.Mod&sdl.KMOD_NUM != 0) {
		// lets let eventTextInput() handle this so we don't need to think
		// about keyboard layouts and shift chars like if shift+2 == '@' or '"'
		return
	}

	mod := keyEventModToCodesModifier(evt.Keysym.Mod)
	keyCode := sr.keyCodeLookup(evt.Keysym.Sym)
	b := codes.GetAnsiEscSeq(sr.keyboardMode.Get(), keyCode, mod)
	if len(b) > 0 {
		sr.term.Reply(b)
	}
}

// SDL doesn't name these, so lets name it ourselves for convenience
const (
	_MOUSE_BUTTON_LEFT = 1 + iota
	_MOUSE_BUTTON_MIDDLE
	_MOUSE_BUTTON_RIGHT
	_MOUSE_BUTTON_X1
	_MOUSE_BUTTON_X2
)

func (tw *termWidgetT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Button {
	case _MOUSE_BUTTON_LEFT:
		sr.term.MouseClick(evt.Button, &types.XY{
			X: (evt.X - sr.border) / sr.glyphSize.X,
			Y: (evt.Y - sr.border) / sr.glyphSize.Y,
		}, func() {
			highlighterStart(sr, evt)
			sr.highlighter.setMode(_HIGHLIGHT_MODE_LINE_RANGE)
		})

	case _MOUSE_BUTTON_MIDDLE:
		sr.clipboardPasteText()

	case _MOUSE_BUTTON_RIGHT:
		highlighterStart(sr, evt)

	case _MOUSE_BUTTON_X1:
		highlighterStart(sr, evt)
		sr.highlighter.setMode(_HIGHLIGHT_MODE_SQUARE)
	}
}

func highlighterStart(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	sr.highlighter = &highlighterT{
		button: evt.Button,
		rect:   &sdl.Rect{X: evt.X, Y: evt.Y},
	}
	if sr.keyModifier != 0 {
		sr.highlighter.modifier(sr.keyModifier)
	}
	sr.keyModifier = 0
}

func (tw *termWidgetT) eventMouseWheel(sr *sdlRender, evt *sdl.MouseWheelEvent) {
	if evt.Direction == sdl.MOUSEWHEEL_FLIPPED {
		sr.term.MouseWheel(int(-evt.Y))
	} else {
		sr.term.MouseWheel(int(evt.Y))
	}
}

func (tw *termWidgetT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	// do nothing
}
