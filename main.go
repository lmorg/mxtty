package main

import (
	"os"
	"os/exec"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/psuedotty"
	"github.com/lmorg/mxtty/typeface"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	x int32 = 80
	y int32 = 25
)

func main() {
	defer typeface.Close()
	defer window.Close()

	err := window.Create("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open("hasklig.ttf", 14)
	//font, err := typeface.Open("monaco.ttf", 16)
	if err != nil {
		panic(err.Error())
	}

	x, y = window.SetTypeFace(font)

	virtTerm := virtualterm.NewTerminal(x, y)
	//virtTerm.Write([]rune(stuff))
	//virtTerm.ExportMxTTY()
	pty, err := psuedotty.NewPTY(x, y)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		p := make([]byte, 2048)
		for {
			i, err := pty.Secondary.Read(p)
			if err != nil {
				panic(err.Error())
				continue
			}
			virtTerm.Write([]rune(string(p[:i])))
			virtTerm.ExportMxTTY()
		}
	}()

	/*go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			_ = window.Update()
			//if err != nil {
			//	panic(err)
			//}
		}
	}()*/

	go func() {
		//cmd := exec.Command("/opt/homebrew/bin/murex")
		cmd := exec.Command("/bin/zsh")
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
					pty.Secondary.Write([]byte{codes.AsciiBackspace})
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
