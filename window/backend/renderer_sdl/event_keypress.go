package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventTextInput(evt *sdl.TextInputEvent) {
	switch {
	case sr.inputBoxActive:
		sr.inputBoxValue += evt.GetText()

	default:
		sr.term.Reply([]byte(evt.GetText()))
	}
}

func (sr *sdlRender) eventKeyPress(evt *sdl.KeyboardEvent) {
	if evt.State != sdl.PRESSED {
		return
	}

	if sr.inputBoxActive {
		switch evt.Keysym.Sym {
		case sdl.K_ESCAPE:
			sr.closeInputBox()
		case sdl.K_RETURN:
			sr.inputBoxCallback(sr.inputBoxValue)
			sr.closeInputBox()
		case sdl.K_BACKSPACE:
			if sr.inputBoxValue != "" {
				sr.inputBoxValue = sr.inputBoxValue[:len(sr.inputBoxValue)-1]
			} else {
				sr.Bell()
			}
		}
		return
	}

	switch evt.Keysym.Sym {
	default:
		if evt.Keysym.Sym < 10000 {
			sr.eventKeyPressMod(evt)
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

func (sr *sdlRender) eventKeyPressMod(evt *sdl.KeyboardEvent) {
	//log.Printf("DEBUG: keycode %s", string(evt.Keysym.Sym))

	switch evt.Keysym.Mod {
	case sdl.KMOD_CTRL, sdl.KMOD_LCTRL, sdl.KMOD_RCTRL:
		if evt.Keysym.Sym > '`' && evt.Keysym.Sym < 'z' {
			sr.term.Reply([]byte{byte(evt.Keysym.Sym) - 0x60})
		}

	case sdl.KMOD_ALT, sdl.KMOD_LALT:
		switch evt.Keysym.Sym {
		case 'f':

		}
	}
}
