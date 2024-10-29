package elementCsv

import (
	"github.com/lmorg/mxtty/types"
)

func (el *ElementCsv) MouseClick(button uint8, pos *types.XY) {
	if pos.Y != 0 {
		el.renderer.DisplayInputBox("Please input desired SQL filter:", el.filter, func(filter string) {
			el.filter = filter
			err := el.runQuery()
			if err != nil {
				el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot sort table: "+err.Error())
			}
		})
		return
	}

	i := len(el.colOffset[0]) - 1
	for ; i > 0; i-- {
		if pos.X >= el.colOffset[0][i] {
			break
		}
	}

	switch button {
	case 1:
		if el.orderBy == i {
			el.orderDesc = !el.orderDesc
		} else {
			el.orderBy = i
			el.orderDesc = false
		}

	case 3:
		el.orderBy = -1
	}

	err := el.runQuery()
	if err != nil {
		//log.Printf("ERROR: cannot sort table: %s", err.Error())
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot sort table: "+err.Error())
	}
}
