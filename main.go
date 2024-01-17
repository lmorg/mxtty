package main

import (
	"github.com/lmorg/mxtty/ptty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
)

func main() {
	getFlags()

	fontName := ""
	shell := "/bin/bash"
	//shell :="/opt/homebrew/bin/murex"

	renderer := backend.Initialise(fontName, 15)
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer)
	pty, err := ptty.NewPTY(term.GetSize())
	if err != nil {
		panic(err.Error())
	}

	term.Start(pty, shell)
	backend.Start(renderer, term)
}
