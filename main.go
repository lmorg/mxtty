package main

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/psuedotty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	renderer := backend.Start()
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer)
	pty, err := psuedotty.NewPTY(term.GetSize())
	if err != nil {
		panic(err.Error())
	}
	term.Start(pty)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TextInputEvent:
				term.Pty.Secondary.WriteString(evt.GetText())

			case *sdl.KeyboardEvent:
				if evt.State == sdl.RELEASED {
					break
				}

				switch evt.Keysym.Sym {
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
				case sdl.K_ESCAPE:
					term.Pty.Secondary.Write([]byte{codes.AsciiEscape})
				}
			}
		}

		sdl.Delay(15)
	}
}
