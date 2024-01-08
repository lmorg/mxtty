package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/veandco/go-sdl2/sdl"
)

func eventTextInput(evt *sdl.TextInputEvent, term *virtualterm.Term) {
	term.Pty.Secondary.WriteString(evt.GetText())
}

func eventKeyPress(evt *sdl.KeyboardEvent, term *virtualterm.Term) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Keysym.Sym {
	case sdl.K_ESCAPE:
		term.Pty.Secondary.Write([]byte{codes.AsciiEscape})
	case sdl.K_TAB:
		term.Pty.Secondary.Write([]byte{'\t'})
	case sdl.K_RETURN:
		term.Pty.Secondary.Write([]byte{'\n'})
	case sdl.K_BACKSPACE:
		term.Pty.Secondary.Write([]byte{codes.IsoBackspace})

	case sdl.K_UP:
		term.Pty.Secondary.Write(codes.AnsiUp)
	case sdl.K_DOWN:
		term.Pty.Secondary.Write(codes.AnsiDown)
	case sdl.K_LEFT:
		term.Pty.Secondary.Write(codes.AnsiBackwards)
	case sdl.K_RIGHT:
		term.Pty.Secondary.Write(codes.AnsiForwards)

	case sdl.K_PAGEDOWN:
		term.Pty.Secondary.Write(codes.AnsiPageDown)
	case sdl.K_PAGEUP:
		term.Pty.Secondary.Write(codes.AnsiPageUp)

	// F-Keys

	case sdl.K_F1:
		term.Pty.Secondary.Write(codes.AnsiF1VT100)
	case sdl.K_F2:
		term.Pty.Secondary.Write(codes.AnsiF2VT100)
	case sdl.K_F3:
		term.Pty.Secondary.Write(codes.AnsiF3VT100)
	case sdl.K_F4:
		term.Pty.Secondary.Write(codes.AnsiF4VT100)
	case sdl.K_F5:
		term.Pty.Secondary.Write(codes.AnsiF5)
	case sdl.K_F6:
		term.Pty.Secondary.Write(codes.AnsiF6)
	case sdl.K_F7:
		term.Pty.Secondary.Write(codes.AnsiF7)
	case sdl.K_F8:
		term.Pty.Secondary.Write(codes.AnsiF8)
	case sdl.K_F9:
		term.Pty.Secondary.Write(codes.AnsiF9)
	case sdl.K_F10:
		term.Pty.Secondary.Write(codes.AnsiF10)
	case sdl.K_F11:
		term.Pty.Secondary.Write(codes.AnsiF11)
	case sdl.K_F12:
		term.Pty.Secondary.Write(codes.AnsiF12)

	}
}
