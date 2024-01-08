package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func eventTextInput(evt *sdl.TextInputEvent, term types.Term) {
	term.Return([]byte(evt.GetText()))
}

func eventKeyPress(evt *sdl.KeyboardEvent, term types.Term) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Keysym.Sym {
	case sdl.K_ESCAPE:
		term.Return([]byte{codes.AsciiEscape})
	case sdl.K_TAB:
		term.Return([]byte{'\t'})
	case sdl.K_RETURN:
		term.Return([]byte{'\n'})
	case sdl.K_BACKSPACE:
		term.Return([]byte{codes.IsoBackspace})

	case sdl.K_UP:
		term.Return(codes.AnsiUp)
	case sdl.K_DOWN:
		term.Return(codes.AnsiDown)
	case sdl.K_LEFT:
		term.Return(codes.AnsiBackwards)
	case sdl.K_RIGHT:
		term.Return(codes.AnsiForwards)

	case sdl.K_PAGEDOWN:
		term.Return(codes.AnsiPageDown)
	case sdl.K_PAGEUP:
		term.Return(codes.AnsiPageUp)

	// F-Keys

	case sdl.K_F1:
		term.Return(codes.AnsiF1VT100)
	case sdl.K_F2:
		term.Return(codes.AnsiF2VT100)
	case sdl.K_F3:
		term.Return(codes.AnsiF3VT100)
	case sdl.K_F4:
		term.Return(codes.AnsiF4VT100)
	case sdl.K_F5:
		term.Return(codes.AnsiF5)
	case sdl.K_F6:
		term.Return(codes.AnsiF6)
	case sdl.K_F7:
		term.Return(codes.AnsiF7)
	case sdl.K_F8:
		term.Return(codes.AnsiF8)
	case sdl.K_F9:
		term.Return(codes.AnsiF9)
	case sdl.K_F10:
		term.Return(codes.AnsiF10)
	case sdl.K_F11:
		term.Return(codes.AnsiF11)
	case sdl.K_F12:
		term.Return(codes.AnsiF12)

	}
}
