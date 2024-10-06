package virtualterm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/config"
)

var ENV_VARS = []string{
	"MXTTY=true",
	"MXTTY_VERSION=" + app.Version(),
	"TERM=xterm-256color",
	"TERM_PROGRAM=mxtty",
}

func init() {
	exe, _ := os.Executable()
	ENV_VARS = append(ENV_VARS, "MXTTY_EXE="+exe)

	_ = os.Unsetenv("TMUX")
	_ = os.Unsetenv("TERM")
	_ = os.Unsetenv("TERM_PROGRAM")

	ENV_VARS = append(os.Environ(), ENV_VARS...)
}

func (term *Term) exec() {
	defaultErr := _exec(term.Pty.File(), config.Config.Shell.Default)
	if defaultErr == nil {
		// success, no need to run fallback shell
		term.renderer.TriggerQuit()
	}

	fallbackErr := _exec(term.Pty.File(), config.Config.Shell.Fallback)
	if fallbackErr == nil {
		// success, no need to run fallback shell
		term.renderer.TriggerQuit()
	}

	panic(fmt.Sprintf(
		"Cannot launch either shell: (Default) %s: (Fallback) %s",
		defaultErr, fallbackErr))
}

func _exec(tty *os.File, command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = ENV_VARS
	cmd.Stdin = tty
	cmd.Stdout = tty
	cmd.Stderr = tty
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Noctty:  false,
		Setctty: true,
		//Ctty:    int(term.Pty.File().Fd()),
		//Setpgid: true,
		Setsid: true,
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	//cmd.SysProcAttr.Ctty = cmd.Process.Pid
	cmd.SysProcAttr.Pgid = cmd.Process.Pid

	err = cmd.Wait()
	if err != nil && strings.HasPrefix(err.Error(), "Signal") {
		return err
	}

	return nil
}
