package elementTable

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (el *ElementTable) MouseClick(button uint8, pos *types.XY) {
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

	var err error
	el._sqlResult, err = el.runQuery()
	if err != nil {
		log.Printf("ERROR: cannot sort table: %s", err.Error())
	}
}
