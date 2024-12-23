package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/debug"
	virtualterm "github.com/lmorg/mxtty/term"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/utils/octal"
	runebuf "github.com/lmorg/mxtty/utils/rune_buf"
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
	Title     string `tmux:"pane_title"`
	Id        string `tmux:"pane_id"`
	Width     int    `tmux:"pane_width"`
	Height    int    `tmux:"pane_height"`
	Active    bool   `tmux:"?pane_active,true,false"`
	WindowId  string `tmux:"window_id"`
	tmux      *Tmux
	buf       *runebuf.Buf
	prefixTtl time.Time
	term      types.Term
}

func (tmux *Tmux) initSessionPanes(renderer types.Renderer, size *types.XY) error {
	panes, err := tmux.sendCommand(CMD_LIST_PANES, reflect.TypeOf(PANE_T{}), "-s")
	if err != nil {
		return err
	}

	for i := range panes.([]any) {
		pane := panes.([]any)[i].(*PANE_T)
		pane.tmux = tmux

		pane.buf = runebuf.New()
		debug.Log(pane)
		tmux.win[pane.WindowId].panes[pane.Id] = pane
		if pane.Active {
			tmux.win[pane.WindowId].activePane = pane
		}
		tmux.pane[pane.Id] = pane

		term := virtualterm.NewTerminal(renderer, size, false)
		pane.term = term
		term.Start(pane)

		command := fmt.Sprintf("capture-pane -J -e -p -t %s", pane.Id)
		resp, err := tmux.SendCommand([]byte(command))
		if err != nil {
			renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
		} else {
			b := bytes.Join(resp.Message, []byte{'\r', '\n'}) // CRLF
			pane.buf.Write(b)
		}

		command = fmt.Sprintf(`display-message -p -t %s "#{e|+:#{cursor_y},1};#{e|+:#{cursor_x},1}H"`, pane.Id)
		resp, err = tmux.SendCommand([]byte(command))
		if err != nil {
			renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
		} else {
			b := append([]byte{codes.AsciiEscape, '['}, resp.Message[0]...)
			pane.buf.Write(b)
		}
	}

	return nil
}

func (tmux *Tmux) newPane(paneId string) *PANE_T {
	pane := &PANE_T{
		Id:   paneId,
		tmux: tmux,
		buf:  runebuf.New(),
	}

	term := virtualterm.NewTerminal(tmux.renderer, tmux.renderer.GetWindowSizeCells(), false)
	term.Start(pane)
	pane.term = term

	tmux.pane[pane.Id] = pane

	go pane._updateInfo(tmux.renderer)

	return pane
}

func (pane *PANE_T) _updateInfo(renderer types.Renderer) {
	err := pane.tmux.updatePaneInfo(pane.Id)
	if err != nil {
		renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
	}
}

type paneInfo struct {
	Id        string `tmux:"pane_id"`
	Title     string `tmux:"pane_title"`
	Width     int    `tmux:"pane_width"`
	Height    int    `tmux:"pane_height"`
	Active    bool   `tmux:"?pane_active,true,false"`
	WindowId  string `tmux:"window_id"`
	WinActive bool   `tmux:"?window_active,true,false"`
}

// updatePaneInfo, paneId is optional. Leave blank to update all panes
func (tmux *Tmux) updatePaneInfo(paneId string) error {
	var filter string
	if paneId != "" {
		filter = fmt.Sprintf("-f '#{m:#{pane_id},%s}'", paneId)
	}

	v, err := tmux.sendCommand(CMD_LIST_PANES, reflect.TypeOf(paneInfo{}), "-s", filter)
	if err != nil {
		return err
	}

	panes, ok := v.([]any)
	if !ok {
		return fmt.Errorf("expecting an array of panes, instead got %T", v)
	}

	for i := range panes {

		info, ok := panes[i].(*paneInfo)
		if !ok {
			return fmt.Errorf("expecting info on a pane, instead got %T", info)
		}

		pane, ok := tmux.pane[info.Id]
		if !ok {
			pane = tmux.newPane(info.Id)
		}
		pane.Title = info.Title
		pane.Width = info.Width
		pane.Height = info.Height
		pane.Active = info.Active
		pane.WindowId = info.WindowId
		pane.term.MakeVisible(info.WinActive)
		pane.term.HasFocus(info.Active)
		pane.term.Resize(&types.XY{X: int32(info.Width), Y: int32(info.Height)})

		tmux.win[pane.WindowId].panes[pane.Id] = pane
		if pane.Active {
			tmux.win[pane.WindowId].activePane = pane
		}
	}

	return nil
}

func (p *PANE_T) File() *os.File { return nil }

func (p *PANE_T) Read() (rune, error) {
	return p.buf.Read()
}

func (p *PANE_T) Write(b []byte) error {
	if len(b) == 0 {
		return errors.New("nothing to write")
	}

	ok, err := p._hotkey(b)
	if ok {
		if err != nil {
			p.tmux.renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
		}
		return nil
	}

	var flags string
	if b[0] == 0 {
		b = b[1:]
	} else {
		flags = "-l"
		b = octal.Escape(b)
	}

	command := []byte(fmt.Sprintf(`send-keys %s -t %s `, flags, p.Id))
	command = append(command, b...)
	_, err = p.tmux.SendCommand(command)
	return err
}

func (p *PANE_T) _hotkey(b []byte) (bool, error) {
	var key string
	if b[0] == 0 {
		key = string(b[1 : len(b)-1])
	} else {
		key = string(b)
	}
	debug.Log(key)

	if p.prefixTtl.Before(time.Now()) {
		if key != p.tmux.keys.prefix {
			// standard key, do nothing
			return false, nil
		}

		// prefix key pressed
		p.prefixTtl = time.Now().Add(2 * time.Second)
		return true, nil
	}

	// run tmux function
	fn, ok := p.tmux.keys.fnTable[key]
	debug.Log(ok)
	if !ok {
		// no function to run, lets treat as standard key
		p.prefixTtl = time.Now()
		return false, nil
	}

	return true, fn.fn(p.tmux)
}

func (p *PANE_T) Resize(size *types.XY) error {
	command := fmt.Sprintf("resize-pane -t %s -x %d -y %d", p.Id, size.X, size.Y)
	_, err := p.tmux.SendCommand([]byte(command))
	if err != nil {
		p.Width = int(size.X)
		p.Height = int(size.Y)
		return err
	}

	return p.tmux.RefreshClient(size)
}

func (tmux *Tmux) ActivePane() *PANE_T {
	return tmux.activeWindow.activePane
}

func (p *PANE_T) Term() types.Term {
	return p.term
}
