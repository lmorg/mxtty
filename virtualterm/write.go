package virtualterm

import (
	"fmt"
	"unsafe"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

type _debugWriteCell struct {
	types.XY
	Rune string
}

func (term *Term) writeCell(r rune) {
	if term._insertOrReplace == _STATE_IRM_INSERT {
		term.csiInsertCharacters(1)
	}

	cell := term.cell()
	cell.Char = r
	cell.Sgr = term.sgr.Copy()

	if debug.Enabled {
		debug.Log(_debugWriteCell{term.curPos, string(r)})
	}

	if term._activeElement != nil {
		cell.Element = term._activeElement
		term._activeElement.ReadCell(cell)
	}

	term.curPos.X++

	if term.curPos.X >= term.size.X && term._noAutoLineWrap {
		term.curPos.X--
	}
}

func (term *Term) appendScrollBuf() {
	if unsafe.Pointer(term.cells) == unsafe.Pointer(&term._normBuf) {
		if len(term._scrollBuf) < config.SCROLLBACK_HISTORY {
			term._scrollBuf = append(term._scrollBuf, term._normBuf[0])
		} else {
			term._scrollBuf = append(term._scrollBuf[1:], term._normBuf[0])
		}
		if term._scrollOffset > 0 {
			term._scrollOffset++
			if term._scrollMsg != nil {
				term._scrollMsg.SetMessage(fmt.Sprintf("Viewing scrollback history. %d lines from end", term._scrollOffset))
				//term.renderer.TriggerRedraw()
			}
		}
	}
}
