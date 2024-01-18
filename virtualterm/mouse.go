package virtualterm

import (
	"log"

	"github.com/lmorg/murex/utils/json"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) MouseClick(button uint8, pos *types.XY) {
	log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))
	if (*term.cells)[pos.Y][pos.X].Element != nil {
		relPos := types.XY{X: pos.X, Y: pos.Y - term.findElementStart(pos)}
		(*term.cells)[pos.Y][pos.X].Element.MouseClick(button, &relPos)
	}
}

func (term *Term) findElementStart(pos *types.XY) int32 {
	y := pos.Y - 1
	for ; y >= 0; y-- {
		if (*term.cells)[y][0].Element != (*term.cells)[pos.Y][pos.X].Element {
			break
		}
	}

	return y + 1
}
