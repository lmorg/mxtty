package tmux

import "fmt"

func fnKeyNewWindow(tmux *Tmux, _ string) error {
	_, err := tmux.SendCommand([]byte("new-window"))
	return err
}

func fnKeyPane(tmux *Tmux, paneId string) error {
	command := fmt.Sprintf("kill-pane -t %s", paneId)
	_, err := tmux.SendCommand([]byte(command))
	return err
}
