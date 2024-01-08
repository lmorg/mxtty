package virtualterm

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/psuedotty"
)

var ENV_VARS = []string{
	"MXTTY=true",
	"MXTTY_VERSION=" + app.Version(),
}

func (term *Term) Start(p *psuedotty.PTY, shell string) {
	term.Pty = p

	go term.exec(shell)
	go term.printLoop()
	go term.slowBlink()
}

func (term *Term) exec(command string) {
	cmd := exec.Command(command)
	cmd.Env = append(os.Environ(), ENV_VARS...)
	cmd.Stdin = term.Pty.Primary
	cmd.Stdout = term.Pty.Primary
	cmd.Stderr = term.Pty.Primary
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Noctty:  false,
		Setctty: true,
		//Ctty:    cmd.Process.Pid,
		//Setpgid: true,
		//Pgid:    cmd.Process.Pid,
		Setsid: true,
	}

	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}

	cmd.SysProcAttr.Ctty = cmd.Process.Pid
	cmd.SysProcAttr.Pgid = cmd.Process.Pid

	err = cmd.Wait()
	if err != nil {
		panic(err.Error())
	}
	os.Exit(0)
}
