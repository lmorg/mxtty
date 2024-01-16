package virtualterm

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/types"
)

var ENV_VARS = []string{
	"MXTTY=true",
	"MXTTY_VERSION=" + app.Version(),
	"TERM=xterm-256color",
}

func init() {
	exe, _ := os.Executable()
	ENV_VARS = append(ENV_VARS, "MXTTY_EXE="+exe)
}

func (term *Term) Start(pty types.Pty, shell string) {
	term.Pty = pty

	go term.exec(shell)
	go term.readLoop()
	go term.slowBlink()
}

func (term *Term) exec(command string) {
	cmd := exec.Command(command)
	cmd.Env = append(os.Environ(), ENV_VARS...)
	cmd.Stdin = term.Pty.File()
	cmd.Stdout = term.Pty.File()
	cmd.Stderr = term.Pty.File()
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
