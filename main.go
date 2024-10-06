package main

import (
	"github.com/lmorg/mxtty/debug/pprof"
	"github.com/lmorg/mxtty/ptty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
)

func main() {
	pprof.Start()
	defer pprof.CleanUp()

	getFlags()

	renderer := backend.Initialise()
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer)
	pty, err := ptty.NewPTY(term.GetSize())
	if err != nil {
		panic(err.Error())
	}

	term.Start(pty)
	backend.Start(renderer, term)
}
