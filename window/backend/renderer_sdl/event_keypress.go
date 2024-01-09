package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func eventTextInput(evt *sdl.TextInputEvent, term types.Term) {
	term.Reply([]byte(evt.GetText()))
}

func eventKeyPress(evt *sdl.KeyboardEvent, term types.Term) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Keysym.Sym {
	case sdl.K_ESCAPE:
		term.Reply([]byte{codes.AsciiEscape})
	case sdl.K_TAB:
		term.Reply([]byte{'\t'})
	case sdl.K_RETURN:
		term.Reply([]byte{'\n'})
	case sdl.K_BACKSPACE:
		term.Reply([]byte{codes.IsoBackspace})

	case sdl.K_UP:
		term.Reply(codes.AnsiUp)
	case sdl.K_DOWN:
		term.Reply(codes.AnsiDown)
	case sdl.K_LEFT:
		term.Reply(codes.AnsiBackwards)
	case sdl.K_RIGHT:
		term.Reply(codes.AnsiForwards)

	case sdl.K_PAGEDOWN:
		term.Reply(codes.AnsiPageDown)
	case sdl.K_PAGEUP:
		term.Reply(codes.AnsiPageUp)

	// F-Keys

	case sdl.K_F1:
		term.Reply(codes.AnsiF1VT100)
	case sdl.K_F2:
		term.Reply(codes.AnsiF2VT100)
	case sdl.K_F3:
		term.Reply(codes.AnsiF3VT100)
	case sdl.K_F4:
		term.Reply(codes.AnsiF4VT100)
	case sdl.K_F5:
		term.Reply(codes.AnsiF5)
	case sdl.K_F6:
		term.Reply(codes.AnsiF6)
	case sdl.K_F7:
		term.Reply(codes.AnsiF7)
	case sdl.K_F8:
		term.Reply(codes.AnsiF8)
	case sdl.K_F9:
		term.Reply(codes.AnsiF9)
	case sdl.K_F10:
		term.Reply(codes.AnsiF10)
	case sdl.K_F11:
		term.Reply(codes.AnsiF11)
	case sdl.K_F12:
		term.Reply(codes.AnsiF12)

	}
}
