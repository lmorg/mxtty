package ptty

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/utils/exit"
)

func MultiplexBypass(command string) {
	fifo := os.Getenv("MXTTY_FIFO")

	size, err := pty.GetsizeFull(os.Stdin)
	if err != nil {
		panic(err) // TODO: don't panic
	}

	secondary, primary, err := pty.Open()
	if err != nil {
		panic(err) // TODO: don't panic
	}

	err = pty.Setsize(primary, size)
	if err != nil {
		panic(err) // TODO: don't panic
	}

	p := &PTY{
		primary:   primary,
		secondary: secondary,
		stream:    make(chan rune),
	}

	f, err := os.OpenFile(fifo, os.O_WRONLY, 0)
	if err != nil {
		panic(err) // TODO: don't panic
	}

	go multiplexBypass(p, f)
	execute(p.primary, command)
}

func execute(f *os.File, command string) {
	cmd := exec.Command(command)
	cmd.Env = config.SetEnv()
	cmd.Stdin = f
	cmd.Stdout = f
	cmd.Stderr = f
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
	exit.Exit(0)
}

func multiplexBypass(p *PTY, fifo *os.File) {
	var (
		err    error
		bypass bool
		b      = make([]byte, 10*1024)
		n      int
		apc    = map[bool]string{
			false: "\x1b_begin;",
			true:  "\x1b_end;",
		}
		c    int
		file = map[bool]*os.File{
			false: os.Stdout,
			true:  fifo,
		}
	)

	for {
		n, err = p.secondary.Read(b)
		if err != nil {
			log.Printf("DEBUG: error returned from multiplexBypass()->Read(): %s", err)
		}

		for i := 0; i < n; i++ {
			if b[i] == apc[bypass][c] {
				i++
				c++

			} else {
				c = 0
			}

			if c == len(apc[bypass]) {
				_, err = file[bypass].Write(b[:n-c])
				if err != nil {
					log.Printf("ERROR: cannot write partial buffer to %s: %s", file[bypass].Name(), err)
				}

				bypass = !bypass

				_, err = file[bypass].Write(b[c:n])
				if err != nil {
					log.Printf("ERROR: cannot write partial buffer to %s: %s", file[bypass].Name(), err)
				}
			}

		}

		switch c {
		case 0:
			_, err = file[bypass].Write(b[:n])
			if err != nil {
				log.Printf("ERROR: cannot write buffer to %s: %s", file[bypass].Name(), err)
			}

		default:
			_, err = file[!bypass].Write(b[:n-c])
			if err != nil {
				log.Printf("ERROR: cannot write partial buffer to %s: %s", file[!bypass].Name(), err)
			}
		}
	}
}
