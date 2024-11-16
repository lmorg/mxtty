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

func (term *Term) writeCell(r rune, el types.Element) {
	if term.writeToElement(r) {
		return
	}

	if term._insertOrReplace == _STATE_IRM_INSERT {
		term.csiInsertCharacters(1)
	}

	if term._curPos.X >= term.size.X && !term._noAutoLineWrap {
		term._curPos.X = 0
		term.lineFeed()
	}

	cell := term.currentCell()
	cell.Char = r
	cell.Sgr = term.sgr.Copy()
	cell.Element = el

	if debug.Enabled {
		debug.Log(_debugWriteCell{term._curPos, string(r)})
	}

	if term._insertOrReplace == _STATE_IRM_REPLACE {
		term._curPos.X++
	}

	if term._ssFrequency == 0 {
		term.renderer.TriggerRedraw()
	}
}

func (term *Term) writeToElement(r rune) (ok bool) {
	if term._activeElement == nil {
		return false
	}

	err := term._activeElement.Write(r)
	if err == nil {
		return true
	}

	term.renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
	term._activeElement = nil
	return false
}

func (term *Term) appendScrollBuf() {
	if unsafe.Pointer(term.cells) == unsafe.Pointer(&term._normBuf) {
		if len(term._scrollBuf) < config.Config.Terminal.ScrollbackHistory {
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
