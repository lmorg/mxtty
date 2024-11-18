package main

import (
	"log"

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

	term := virtualterm.NewTerminal(renderer, size, true)
	pty, err := ptty.NewPTY(size)
	if err != nil {
		panic(err)
	}

	term.Start(pty)
	backend.Start(renderer, term, nil)
}

func tmuxSession() {
	renderer, size := backend.Initialise()
	defer renderer.Close()

	tmuxClient, err := tmux.NewStartSession(renderer, size, tmux.START_ATTACH_SESSION)
	if err != nil {
		if err.Error() != "no sessions" {
			panic(err)
		}

		log.Println("No sessions to attach to. Creating new session.")

		tmuxClient, err = tmux.NewStartSession(renderer, size, tmux.START_NEW_SESSION)
		if err != nil {
			panic(err)
		}
	}

	backend.Start(renderer, tmuxClient.ActivePane().Term(), tmuxClient)
}
