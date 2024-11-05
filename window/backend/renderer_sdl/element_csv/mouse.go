package elementCsv

import (
	"strings"

	"github.com/lmorg/mxtty/types"
	"golang.design/x/clipboard"
)

func (el *ElementCsv) MouseClick(button uint8, pos *types.XY, callback types.MouseClickCallback) {
	if pos.Y != 0 {
		switch button {
		case 1:
			for i := range el.boundaries {
				if pos.X <= el.boundaries[i] {
					var start int32
					if i != 0 {
						start = el.boundaries[i-1]
					}
					cell := el.table[pos.Y-1][start:el.boundaries[i]]
					clipboard.Write(clipboard.FmtText, []byte(strings.TrimSpace(cell)))
					el.renderer.DisplayNotification(types.NOTIFY_INFO, "Cell copied to clipboard")
					return
				}
			}
			callback()

		case 3:
			el.renderer.DisplayInputBox("Please input desired SQL filter:", el.filter, func(filter string) {
				el.filter = filter
				err := el.runQuery()
				if err != nil {
					el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot sort table: "+err.Error())
				}
			})

		default:
			callback()
		}

		return
	}

	var column, width int
	for column = range el.width {
		width += el.width[column]
		if int(pos.X) < width-1 {
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
