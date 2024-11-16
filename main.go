package main

import (
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/debug/pprof"
	"github.com/lmorg/mxtty/ptty"
	virtualterm "github.com/lmorg/mxtty/term"
	"github.com/lmorg/mxtty/tmux"
	"github.com/lmorg/mxtty/window/backend"
)

func main() {
	pprof.Start()
	defer pprof.CleanUp()

	getFlags()

	if config.Config.Tmux.Enabled {
		tmuxSession()
	} else {
		regularSession()
	}
}

func regularSession() {
	renderer, size := backend.Initialise()
	defer renderer.Close()

	term := virtualterm.NewTerminal(renderer, size, false)
	pty, err := ptty.NewPTY(size)
	if err != nil {
		panic(err)
	}

	term.Start(pty)
	backend.Start(renderer, term)
}

func tmuxSession() {
	err := tmux.NewTmuxAttachSession()
	if err != nil {
		panic(err)
	}
}
