package tmux

import (
	"fmt"
	"os"
	"reflect"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

/*
	pane_active                1 if active pane
	pane_at_bottom             1 if pane is at the bottom of window
	pane_at_left               1 if pane is at the left of window
	pane_at_right              1 if pane is at the right of window
	pane_at_top                1 if pane is at the top of window
	pane_bg                    Pane background colour
	pane_bottom                Bottom of pane
	pane_current_command       Current command if available
	pane_current_path          Current path if available
	pane_dead                  1 if pane is dead
	pane_dead_signal           Exit signal of process in dead pane
	pane_dead_status           Exit status of process in dead pane
	pane_dead_time             Exit time of process in dead pane
	pane_fg                    Pane foreground colour
	pane_format                1 if format is for a pane
	pane_height                Height of pane
	pane_id                #D  Unique pane ID
	pane_in_mode               1 if pane is in a mode
	pane_index             #P  Index of pane
	pane_input_off             1 if input to pane is disabled
	pane_key_mode              Extended key reporting mode in this pane
	pane_last                  1 if last pane
	pane_left                  Left of pane
	pane_marked                1 if this is the marked pane
	pane_marked_set            1 if a marked pane is set
	pane_mode                  Name of pane mode, if any
	pane_path                  Path of pane (can be set by application)
	pane_pid                   PID of first process in pane
	pane_pipe                  1 if pane is being piped
	pane_right                 Right of pane
	pane_search_string         Last search string in copy mode
	pane_start_command         Command pane started with
	pane_start_path            Path pane started with
	pane_synchronized          1 if pane is synchronized
	pane_tabs                  Pane tab positions
	pane_title             #T  Title of pane (can be set by application)
	pane_top                   Top of pane
	pane_tty                   Pseudo terminal of pane
	pane_unseen_changes        1 if there were changes in pane while in mode
	pane_width                 Width of pane
*/

var CMD_LIST_PANES = "list-panes"

type PANE_T struct {
	Title  string `tmux:"pane_title"`
	Id     string `tmux:"pane_id"`
	Width  int    `tmux:"pane_width"`
	Height int    `tmux:"pane_height"`
	Active bool   `tmux:"?pane_active,true,false"`
	tmux   *Tmux
	buf    chan rune
	closed bool
}

func (p *PANE_T) File() *os.File { return nil }

func (p *PANE_T) respFromTmux(b []byte) {
	//debug.Log(p.Id)
	for _, r := range []rune(string(b)) {
		p.buf <- r
	}
}

func (p *PANE_T) Read() rune {
	return <-p.buf
}

func (tmux *Tmux) initSessionPanes() error {
	for _, win := range tmux.wins {
		win.panes = make(map[string]*PANE_T)

		panes, err := tmux.sendCommand(CMD_LIST_PANES, reflect.TypeOf(PANE_T{}), "-t", win.Id)
		if err != nil {
			return err
		}

		for i := range panes.([]any) {
			pane := panes.([]any)[i].(*PANE_T)
			pane.tmux = tmux
			pane.buf = make(chan rune)
			debug.Log(pane)
			win.panes[pane.Id] = pane
			if pane.Active {
				win.activePane = pane
			}
			tmux.panes[pane.Id] = pane
		}
	}

	return nil
}

func (p *PANE_T) Write(b []byte) error {
	command := []byte(fmt.Sprintf(`send-keys -t %s `, p.Id))
	command = append(command, b...)
	_, err := p.tmux.SendCommand(command)
	return err
}

func (p *PANE_T) Resize(size *types.XY) error {
	return p.tmux.RefreshClient(size)
}

func (tmux *Tmux) ActivePane() *PANE_T {
	return tmux.activeWindow.activePane
}
