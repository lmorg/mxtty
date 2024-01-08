package virtualterm

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/lmorg/mxtty/psuedotty"
)

func (term *Term) Start(p *psuedotty.PTY) {
	term.Pty = p

	//command :="/opt/homebrew/bin/murex"
	command := "/bin/zsh"

	go term.exec(command)
	go term.printLoop()
	go term.renderLoop()
	go term.slowBlink()
}

func (term *Term) exec(command string) {
	cmd := exec.Command(command)
	//cmd.Env = append(os.Environ(), "TERM=mxtty")
	cmd.Stdin = term.Pty.Primary
	cmd.Stdout = term.Pty.Primary
	cmd.Stderr = term.Pty.Primary

	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		panic(err.Error())
	}
	os.Exit(0)
}

func (term *Term) renderLoop() {
	for {
		time.Sleep(5 * time.Millisecond)
		term.Render()
		err := term.renderer.Update()
		if err != nil {
			log.Printf("error in renderer: %s", err.Error())
		}
	}
}
