package virtualterm

import (
	"fmt"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Match(pos *types.XY) {
	if term.IsAltBuf() {
		term.renderer.DisplayNotification(types.NOTIFY_WARN, "Match is not supported in alt buffer")
		return
	}

	r := (*term.cells)[pos.Y][pos.X].Rune()
	switch r {
	default:
		term.renderer.DisplayNotification(types.NOTIFY_WARN, fmt.Sprintf("Match cannot run match against the character '%s'", string(r)))
		return
	}
}
