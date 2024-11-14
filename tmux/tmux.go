package tmux

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend"
)

type Tmux struct {
	cmd   *exec.Cmd
	tty   *os.File
	resp  chan *tmuxResponseT
	wins  map[string]*windowT
	panes map[string]*paneT

	activeWindow *windowT
}

var (
	_RESP_OUTPUT = []byte("%output")
	_RESP_BEGIN  = []byte("%begin")
	_RESP_END    = []byte("%end")
	_RESP_ERROR  = []byte("%error")
)

type tmuxResponseT struct {
	Message [][]byte
	IsErr   bool
}

func NewTmuxAttachSession() error {
	tmux := new(Tmux)
	tmux.resp = make(chan *tmuxResponseT)
	tmux.panes = make(map[string]*paneT)

	var err error
	resp := new(tmuxResponseT)
	go func() {
		<-tmux.resp // ignore the first block
	}()

	tmux.cmd = exec.Command("tmux", "-CC", "attach-session")
	tmux.tty, err = pty.Start(tmux.cmd)
	if err != nil {
		return err
	}

	_, _ = tmux.tty.Read(make([]byte, 7))
	// Discard the following because it's just setting mode:
	//    \u001bP1000p

	go func() {
		scanner := bufio.NewScanner(tmux.tty)

		for scanner.Scan() {
			b := scanner.Bytes()
			debug.Log(string(b))
			switch {
			case bytes.HasPrefix(b, _RESP_OUTPUT):
				params := bytes.SplitN(b, []byte{' '}, 3)
				pane, ok := tmux.panes[string(params[1])]
				if !ok {
					panic(fmt.Sprintf("unknown pane ID: %s", string(params[1])))
				}
				pane.respFromTmux(unescapeOctal(params[2]))

			case bytes.HasPrefix(b, _RESP_BEGIN):
				resp = new(tmuxResponseT)

			case bytes.HasPrefix(b, _RESP_ERROR):
				resp.IsErr = true
				fallthrough

			case bytes.HasPrefix(b, _RESP_END):
				tmux.resp <- resp

			default:
				//if len(b) > 0 && b[0] == '%' {

				//				}
				resp.Message = append(resp.Message, b)
			}
		}
	}()

	renderer, size := backend.Initialise()
	defer renderer.Close()

	err = tmux.initSessionWindows()
	if err != nil {
		return err
	}

	err = tmux.initSessionPanes()
	if err != nil {
		return err
	}

	term := virtualterm.NewTerminal(renderer, size, false)
	term.Start(tmux.ActivePane())

	s := fmt.Sprintf("refresh-client -l %s\n", tmux.ActivePane().Id)
	_, err = tmux.SendCommand([]byte(s))
	if err != nil {
		return err
	}

	backend.Start(renderer, term)
	return nil // could shouldn't reach this point
}

func (tmux *Tmux) SendCommand(b []byte) (*tmuxResponseT, error) {
	debug.Log(string(b))
	_, err := tmux.tty.Write(b)
	if err != nil {
		return nil, err
	}

	resp := <-tmux.resp
	if resp.IsErr {
		return nil, fmt.Errorf("tmux command failed: %s", string(bytes.Join(resp.Message, []byte(": "))))
	}

	return resp, nil
}
