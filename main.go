package main

import (
	"os"
	"os/exec"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/psuedotty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	renderer := backend.Start()
	defer renderer.Close()

	virtTerm := virtualterm.NewTerminal(renderer)
	pty, err := psuedotty.NewPTY(virtTerm)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		//cmd := exec.Command("/opt/homebrew/bin/murex")
		cmd := exec.Command("/bin/zsh")
		//cmd.Env = append(os.Environ(), "TERM=mxtty")
		cmd.Stdin = pty.Primary
		cmd.Stdout = pty.Primary
		cmd.Stderr = pty.Primary

		err := cmd.Start()
		if err != nil {
			panic(err.Error())
		}

		err = cmd.Wait()
		if err != nil {
			panic(err.Error())
		}
		os.Exit(0)
	}()

	// Run infinite loop until user closes the window
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TextInputEvent:
				pty.Secondary.WriteString(evt.GetText())

			case *sdl.KeyboardEvent:
				if evt.State == sdl.RELEASED {
					break
				}

				switch evt.Keysym.Sym {
				case sdl.K_TAB:
					pty.Secondary.Write([]byte{'\t'})
				case sdl.K_RETURN:
					pty.Secondary.Write([]byte{'\n'})
				case sdl.K_BACKSPACE:
					pty.Secondary.Write([]byte{codes.IsoBackspace})
				case sdl.K_UP:
					pty.Secondary.Write(codes.AnsiUp)
				case sdl.K_DOWN:
					pty.Secondary.Write(codes.AnsiDown)
				case sdl.K_LEFT:
					pty.Secondary.Write(codes.AnsiBackwards)
				case sdl.K_RIGHT:
					pty.Secondary.Write(codes.AnsiForwards)
				case sdl.K_ESCAPE:
					pty.Secondary.Write([]byte{codes.AsciiEscape})
				}
			}
		}

		sdl.Delay(15)
	}
}
