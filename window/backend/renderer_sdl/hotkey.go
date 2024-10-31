package rendersdl

import (
	"fmt"

	"github.com/lmorg/mxtty/types"
	"golang.design/x/hotkey"
)

func (sr *sdlRender) registerHotkey() {
	sr.hk = hotkey.New([]hotkey.Modifier{}, hotkey.KeyF5)
	err := sr.hk.Register()
	if err != nil {
		sr.DisplayNotification(types.NOTIFY_ERROR, fmt.Sprintf("Unable to set hotkey: %s", err.Error()))
	}
}
