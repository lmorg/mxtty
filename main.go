package main

import (
	"github.com/lmorg/mxtty/debug/pprof"
	"github.com/lmorg/mxtty/ptty"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
)

func main() {
	pprof.Start()
	//defer exit.Exit(0)

	getFlags()

	renderer, size := backend.Initialise()
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer, size)
	pty, err := ptty.NewPTY(size)
	if err != nil {
		panic(err.Error())
	}

	term.Start(pty)
	backend.Start(renderer, term)
}
