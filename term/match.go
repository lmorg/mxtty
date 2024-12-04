package virtualterm

import "github.com/lmorg/mxtty/types"

func (term *Term) Match(pos *types.XY) {
	if term.IsAltBuf() {
		term.renderer.DisplayNotification(types.NOTIFY_WARN, "Match is not supported in alt buffer")
		return
	}
}
