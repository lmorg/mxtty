package elementCsv

import (
	"strings"

	"github.com/lmorg/mxtty/types"
	"golang.design/x/clipboard"
)

func (el *ElementCsv) MouseClick(pos *types.XY, button uint8, clicks uint8, callback types.EventIgnoredCallback) {
	pos.X -= el.renderOffset

	if pos.Y != 0 {
		if button != 1 {
			callback()
			return
		}
		switch clicks {
		case 1:
			if int(pos.Y) > len(el.table) {
				callback()
				return
			}
			for i := range el.boundaries {
				if pos.X <= el.boundaries[i] {
					var start int32
					if i != 0 {
						start = el.boundaries[i-1]
					}
					cell := string(el.table[pos.Y-1][start:el.boundaries[i]])
					clipboard.Write(clipboard.FmtText, []byte(strings.TrimSpace(cell)))
					el.renderer.DisplayNotification(types.NOTIFY_INFO, "Cell copied to clipboard")
					return
				}
			}
			callback()
			return

		case 2:
			el.renderer.DisplayInputBox("Please input desired SQL filter:", el.filter, func(filter string) {
				el.filter = filter
				err := el.runQuery()
				if err != nil {
					el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot sort table: "+err.Error())
				}
			})

		default:
			callback()
			return
		}

		return
	}

	var column int
	for column = range el.boundaries {
		if int(pos.X) < int(el.boundaries[column]) {
			break
		}
	}

	column++ // columns count from 1 because of rowid

	switch button {
	case 1:
		if el.orderByIndex == column {
			el.orderDesc = !el.orderDesc
		} else {
			el.orderByIndex = column
			el.orderDesc = false
		}

	case 3:
		el.orderByIndex = 0
	}

	err := el.runQuery()
	if err != nil {
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot sort table: "+err.Error())
	}
}

func (el *ElementCsv) MouseWheel(_ *types.XY, movement *types.XY, callback types.EventIgnoredCallback) {
	if movement.X == 0 {
		callback()
		return
	}

	termX := el.renderer.GetTermSize().X
	width := el.boundaries[len(el.boundaries)-1]

	if width < termX {
		callback()
		return
	}

	el.renderOffset += (-movement.X * el.renderer.GetGlyphSize().X)

	if el.renderOffset > 0 {
		el.renderOffset = 0
	}

	if el.renderOffset < -(width - termX) {
		el.renderOffset = -(width - termX)
	}
}

func (el *ElementCsv) MouseMotion(pos *types.XY, move *types.XY, callback types.EventIgnoredCallback) {
	/*pos.X -= el.renderOffset

	var column int
	for column = range el.boundaries {
		if int(pos.X) < int(el.boundaries[column]) {
			break
		}
	}*/

	if pos.Y < 1 || pos.X > el.boundaries[len(el.boundaries)-1] {
		el.highlight = nil
		return
	}

	el.highlight = pos
	el.renderer.TriggerRedraw()

	//callback()
}

func (el *ElementCsv) MouseOut() {
	el.highlight = nil
	el.renderer.TriggerRedraw()
}
