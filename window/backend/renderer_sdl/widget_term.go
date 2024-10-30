package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

type termWidgetT struct{}

func (tw *termWidgetT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	sr.term.Reply([]byte(evt.GetText()))
}

func (tw *termWidgetT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	sr.keyModifier = evt.Keysym.Mod

	switch evt.Keysym.Sym {
	default:
		if evt.Keysym.Sym < 10000 {
			tw.eventKeyPressMod(sr, evt)
		}

	case sdl.K_ESCAPE:
		sr.term.Reply([]byte{codes.AsciiEscape})
	case sdl.K_TAB:
		sr.term.Reply([]byte{'\t'})
	case sdl.K_RETURN:
		sr.term.Reply([]byte{'\n'})
	case sdl.K_BACKSPACE:
		sr.term.Reply([]byte{codes.AsciiBackspace})

	case sdl.K_UP:
		sr.ansiReply(codes.AnsiUp)
	case sdl.K_DOWN:
		sr.ansiReply(codes.AnsiDown)
	case sdl.K_LEFT:
		sr.ansiReply(codes.AnsiLeft)
	case sdl.K_RIGHT:
		sr.ansiReply(codes.AnsiRight)

	case sdl.K_PAGEDOWN:
		sr.ansiReply(codes.AnsiPageDown)
	case sdl.K_PAGEUP:
		sr.ansiReply(codes.AnsiPageUp)

	// F-Keys

	case sdl.K_F1:
		sr.ansiReply(codes.AnsiF1)
	case sdl.K_F2:
		sr.ansiReply(codes.AnsiF2)
	case sdl.K_F3:
		sr.ansiReply(codes.AnsiF3)
	case sdl.K_F4:
		sr.ansiReply(codes.AnsiF4)
	case sdl.K_F5:
		sr.ansiReply(codes.AnsiF5)
	case sdl.K_F6:
		sr.ansiReply(codes.AnsiF6)
	case sdl.K_F7:
		sr.ansiReply(codes.AnsiF7)
	case sdl.K_F8:
		sr.ansiReply(codes.AnsiF8)
	case sdl.K_F9:
		sr.ansiReply(codes.AnsiF9)
	case sdl.K_F10:
		sr.ansiReply(codes.AnsiF10)
	case sdl.K_F11:
		sr.ansiReply(codes.AnsiF11)
	case sdl.K_F12:
		sr.ansiReply(codes.AnsiF12)

	}
}

func (tw *termWidgetT) eventKeyPressMod(sr *sdlRender, evt *sdl.KeyboardEvent) {
	//log.Printf("DEBUG: keycode %s", string(evt.Keysym.Sym))

	switch {
	case evt.Keysym.Mod&sdl.KMOD_CTRL != 0:
		fallthrough
	case evt.Keysym.Mod&sdl.KMOD_LCTRL != 0:
		fallthrough
	case evt.Keysym.Mod&sdl.KMOD_RCTRL != 0:
		if evt.Keysym.Sym > '`' && evt.Keysym.Sym < 'z' {
			sr.term.Reply([]byte{byte(evt.Keysym.Sym) - 0x60})
		}

		//case sdl.KMOD_ALT, sdl.KMOD_LALT:
		/*switch evt.Keysym.Sym {
		case 'f':

		}*/
	}
}

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

	if evt.Which == sdl.TOUCH_MOUSEID {
		// touchpad events

		switch evt.Button {
		case _MOUSE_BUTTON_LEFT:
			sr.term.MouseClick(evt.Button, &types.XY{
				X: (evt.X - sr.border) / sr.glyphSize.X,
				Y: (evt.Y - sr.border) / sr.glyphSize.Y,
			})

		case _MOUSE_BUTTON_MIDDLE:
			highlighterStart(sr, evt)
		}

	} else {

		// mouse events

		switch evt.Button {
		case _MOUSE_BUTTON_LEFT:
			sr.term.MouseClick(evt.Button, &types.XY{
				X: (evt.X - sr.border) / sr.glyphSize.X,
				Y: (evt.Y - sr.border) / sr.glyphSize.Y,
			})

		case _MOUSE_BUTTON_MIDDLE:
			sr.clipboardPasteText()

		case _MOUSE_BUTTON_RIGHT:
			highlighterStart(sr, evt)

		case _MOUSE_BUTTON_X1:
			highlighterStart(sr, evt)
			sr.highlighter.mode = _HIGHLIGHT_MODE_LINES
		}
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
	sr.highlighter.mode = _HIGHLIGHT_MODE_PNG
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
