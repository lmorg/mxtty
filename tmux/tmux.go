package tmux

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/utils/octal"
)

/*
	Reference documentation used:
	- tmux man page: https://man.openbsd.org/tmux#CONTROL_MODE

	*** Control Mode: ***

	tmux offers a textual interface called control mode. This allows applications to communicate with tmux using a simple text-only protocol.

	In control mode, a client sends tmux commands or command sequences terminated by newlines on standard input. Each command will produce one block of output on standard output. An output block consists of a %begin line followed by the output (which may be empty). The output block ends with a %end or %error. %begin and matching %end or %error have three arguments: an integer time (as seconds from epoch), command number and flags (currently not used). For example:

	%begin 1363006971 2 1
	0: ksh* (1 panes) [80x24] [layout b25f,80x24,0,0,2] @2 (active)
	%end 1363006971 2 1

	The refresh-client -C command may be used to set the size of a client in control mode.

	In control mode, tmux outputs notifications. A notification will never occur inside an output block.

	The following notifications are defined:

	%client-detached client
		The client has detached.
	%client-session-changed client session-id name
		The client is now attached to the session with ID session-id, which is named name.
	%config-error error
		An error has happened in a configuration file.
	%continue pane-id
		The pane has been continued after being paused (if the pause-after flag is set, see refresh-client -A).
	%exit [reason]
		The tmux client is exiting immediately, either because it is not attached to any session or an error occurred. If present, reason describes why the client exited.
	%extended-output pane-id age ... : value
		New form of %output sent when the pause-after flag is set. age is the time in milliseconds for which tmux had buffered the output before it was sent. Any subsequent arguments up until a single ‘:’ are for future use and should be ignored.
	%layout-change window-id window-layout window-visible-layout window-flags
		The layout of a window with ID window-id changed. The new layout is window-layout. The window's visible layout is window-visible-layout and the window flags are window-flags.
	%message message
		A message sent with the display-message command.
	%output pane-id value
		A window pane produced output. value escapes non-printable characters and backslash as octal \xxx.
	%pane-mode-changed pane-id
		The pane with ID pane-id has changed mode.
	%paste-buffer-changed name
		Paste buffer name has been changed.
	%paste-buffer-deleted name
		Paste buffer name has been deleted.
	%pause pane-id
		The pane has been paused (if the pause-after flag is set).
	%session-changed session-id name
		The client is now attached to the session with ID session-id, which is named name.
	%session-renamed name
		The current session was renamed to name.
	%session-window-changed session-id window-id
		The session with ID session-id changed its active window to the window with ID window-id.
	%sessions-changed
		A session was created or destroyed.
	%subscription-changed name session-id window-id window-index pane-id ... : value
		The value of the format associated with subscription name has changed to value. See refresh-client -B. Any arguments after pane-id up until a single ‘:’ are for future use and should be ignored.
	%unlinked-window-add window-id
		The window with ID window-id was created but is not linked to the current session.
	%unlinked-window-close window-id
		The window with ID window-id, which is not linked to the current session, was closed.
	%unlinked-window-renamed window-id
		The window with ID window-id, which is not linked to the current session, was renamed.
	%window-add window-id
		The window with ID window-id was linked to the current session.
	%window-close window-id
		The window with ID window-id closed.
	%window-pane-changed window-id pane-id
		The active pane in the window with ID window-id changed to the pane with ID pane-id.
	%window-renamed window-id name
		The window with ID window-id was renamed to name.
		All the notifications listed in the CONTROL MODE section are hooks (without any arguments), except %exit. The following additional hooks are available:
*/

var (
	_RESP_OUTPUT = []byte("%output")
	_RESP_BEGIN  = []byte("%begin")
	_RESP_END    = []byte("%end")
	_RESP_ERROR  = []byte("%error")

	// currently unused
	_RESP_CLIENT_DETACHED         = []byte("%client-detached")
	_RESP_CLIENT_SESSION_CHANGED  = []byte("%client-session-changed")
	_RESP_CONFIG_ERROR            = []byte("%config-error")
	_RESP_CONTINUE                = []byte("%continue")
	_RESP_EXIT                    = []byte("%exit")
	_RESP_EXTENDED_OUTPUT         = []byte("%extended-output")
	_RESP_LAYOUT_CHANGE           = []byte("%layout-change")
	_RESP_MESSAGE                 = []byte("%message")
	_RESP_PANE_MODE_CHANGED       = []byte("%pane-mode-changed")
	_RESP_PASTE_BUFFER_CHANGED    = []byte("%paste-buffer-changed")
	_RESP_PASTE_BUFFER_DELETED    = []byte("%paste-buffer-deleted")
	_RESP_PAUSE                   = []byte("%pause")
	_RESP_SESSION_CHANGED         = []byte("%session-changed")
	_RESP_SESSION_RENAMED         = []byte("%session-renamed")
	_RESP_SESSION_WINDOW_CHANGED  = []byte("%session-window-changed")
	_RESP_SESSIONS_CHANGED        = []byte("%sessions-changed")
	_RESP_SUBSCRIPTION_CHANGED    = []byte("%subscription-changed")
	_RESP_UNLINKED_WINDOW_ADD     = []byte("%unlinked-window-add")
	_RESP_UNLINKED_WINDOW_CLOSE   = []byte("%unlinked-window-close")
	_RESP_UNLINKED_WINDOW_RENAMED = []byte("%unlinked-window-renamed")
	_RESP_WINDOW_ADD              = []byte("%window-add")
	_RESP_WINDOW_CLOSE            = []byte("%window-close")
	_RESP_WINDOW_PANE_CHANGED     = []byte("%window-pane-changed")
	_RESP_WINDOW_RENAMED          = []byte("%window-renamed")
)

var respIgnored = [][]byte{
	_RESP_CLIENT_DETACHED,
	_RESP_CLIENT_SESSION_CHANGED,
	_RESP_CONFIG_ERROR,
	_RESP_CONTINUE,
	_RESP_EXIT,
	_RESP_EXTENDED_OUTPUT,
	_RESP_LAYOUT_CHANGE,
	_RESP_PANE_MODE_CHANGED,
	_RESP_PASTE_BUFFER_CHANGED,
	_RESP_PASTE_BUFFER_DELETED,
	_RESP_PAUSE,
	_RESP_SESSION_CHANGED,
	_RESP_SESSION_RENAMED,
	_RESP_SESSION_WINDOW_CHANGED,
	_RESP_SESSIONS_CHANGED,
	_RESP_SUBSCRIPTION_CHANGED,
	_RESP_UNLINKED_WINDOW_ADD,
	_RESP_UNLINKED_WINDOW_CLOSE,
	_RESP_UNLINKED_WINDOW_RENAMED,
	_RESP_WINDOW_ADD,
	_RESP_WINDOW_CLOSE,
	_RESP_WINDOW_PANE_CHANGED,
}

type Tmux struct {
	cmd  *exec.Cmd
	tty  *os.File
	resp chan *tmuxResponseT
	win  map[string]*WINDOW_T
	pane map[string]*PANE_T

	activeWindow *WINDOW_T
}

type tmuxResponseT struct {
	Message [][]byte
	IsErr   bool
}

func NewTmuxAttachSession(renderer types.Renderer, size *types.XY) (*Tmux, error) {
	tmux := new(Tmux)
	tmux.resp = make(chan *tmuxResponseT)
	tmux.pane = make(map[string]*PANE_T)

	var err error
	resp := new(tmuxResponseT)
	go func() {
		<-tmux.resp // ignore the first block
	}()

	tmux.cmd = exec.Command("tmux", "-CC", "attach-session")
	tmux.cmd.Env = config.SetEnv()
	tmux.tty, err = pty.Start(tmux.cmd)
	if err != nil {
		return nil, err
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
				pane, ok := tmux.pane[string(params[1])]
				if !ok {
					panic(fmt.Sprintf("unknown pane ID: %s", string(params[1])))
				}
				pane.respFromTmux(octal.Unescape(params[2]))

			case bytes.HasPrefix(b, _RESP_BEGIN):
				resp = new(tmuxResponseT)

			case bytes.HasPrefix(b, _RESP_ERROR):
				resp.IsErr = true
				fallthrough

			case bytes.HasPrefix(b, _RESP_END):
				tmux.resp <- resp

			case bytes.HasPrefix(b, _RESP_MESSAGE):
				msg := b[len(_RESP_MESSAGE):]
				renderer.DisplayNotification(types.NOTIFY_INFO, string(msg))

			case bytes.HasPrefix(b, _RESP_WINDOW_RENAMED):
				params := bytes.SplitN(b, []byte{' '}, 3)
				tmux.win[string(params[1])].Name = string(params[2])
				renderer.RefreshWindowList()

			default:
				// ignore anything that looks like a notification
				if ignoreResponse(b) {
					continue
				}

				resp.Message = append(resp.Message, b)
			}
		}
	}()

	err = tmux.RefreshClient(size)
	if err != nil {
		return nil, err
	}

	err = tmux.initSessionWindows()
	if err != nil {
		return nil, err
	}

	err = tmux.initSessionPanes(renderer, size)
	if err != nil {
		return nil, err
	}

	tmux.ActivePane().term.MakeVisible(true)

	return tmux, nil // could shouldn't reach this point
}

func ignoreResponse(b []byte) bool {
	if len(b) > 0 && b[0] == '%' {
		for _, notification := range respIgnored {
			if bytes.HasPrefix(b, notification) {
				return true
			}
		}
	}

	return false
}

func (tmux *Tmux) SendCommand(b []byte) (*tmuxResponseT, error) {
	debug.Log(string(b))
	_, err := tmux.tty.Write(append(b, '\n'))
	if err != nil {
		return nil, err
	}

	resp := <-tmux.resp

	if resp.IsErr {
		return nil, fmt.Errorf("tmux command failed: %s", string(bytes.Join(resp.Message, []byte(": "))))
	}

	return resp, nil
}
