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
		sr.term.Reply(codes.AnsiUp)
	case sdl.K_DOWN:
		sr.term.Reply(codes.AnsiDown)
	case sdl.K_LEFT:
		sr.term.Reply(codes.AnsiBackwards)
	case sdl.K_RIGHT:
		sr.term.Reply(codes.AnsiForwards)

	case sdl.K_PAGEDOWN:
		sr.term.Reply(codes.AnsiPageDown)
	case sdl.K_PAGEUP:
		sr.term.Reply(codes.AnsiPageUp)

	// F-Keys

	case sdl.K_F1:
		sr.term.Reply(codes.AnsiF1VT100)
	case sdl.K_F2:
		sr.term.Reply(codes.AnsiF2VT100)
	case sdl.K_F3:
		sr.term.Reply(codes.AnsiF3VT100)
	case sdl.K_F4:
		sr.term.Reply(codes.AnsiF4VT100)
	case sdl.K_F5:
		sr.term.Reply(codes.AnsiF5)
	case sdl.K_F6:
		sr.term.Reply(codes.AnsiF6)
	case sdl.K_F7:
		sr.term.Reply(codes.AnsiF7)
	case sdl.K_F8:
		sr.term.Reply(codes.AnsiF8)
	case sdl.K_F9:
		sr.term.Reply(codes.AnsiF9)
	case sdl.K_F10:
		sr.term.Reply(codes.AnsiF10)
	case sdl.K_F11:
		sr.term.Reply(codes.AnsiF11)
	case sdl.K_F12:
		sr.term.Reply(codes.AnsiF12)

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
	_MOUSE_BUTTON_LEFT = 1 << iota
	_MOUSE_BUTTON_RIGHT
	_MOUSE_BUTTON_MIDDLE
)

func (tw *termWidgetT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Button {
	case _MOUSE_BUTTON_LEFT, _MOUSE_BUTTON_RIGHT:
		sr.term.MouseClick(evt.Button, &types.XY{
			X: (evt.X - sr.border) / sr.glyphSize.X,
			Y: (evt.Y - sr.border) / sr.glyphSize.Y,
		})

	case _MOUSE_BUTTON_MIDDLE, _MOUSE_BUTTON_LEFT | _MOUSE_BUTTON_RIGHT:
		sr.highlighter = &highlighterT{
			button: evt.Button,
			rect:   &sdl.Rect{X: evt.X, Y: evt.Y},
		}
		if sr.keyModifier != 0 {
			sr.highlighter.modifier(sr.keyModifier)
		}
		sr.keyModifier = 0
	}
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
