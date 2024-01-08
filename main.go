package main

import (
	"github.com/lmorg/mxtty/psuedotty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
)

func main() {
	fontName := "hasklig.ttf"
	//fontName := "monaco.ttf"
	shell := "/bin/bash"
	//shell :="/opt/homebrew/bin/murex"

	renderer := backend.Initialise(fontName, 15)
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer)
	pty, err := psuedotty.NewPTY(term.GetSize())
	if err != nil {
		panic(err.Error())
	}

	term.Start(pty, shell)
	backend.Start(renderer, term)
}
