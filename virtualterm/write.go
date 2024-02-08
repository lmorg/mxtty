package virtualterm

import (
	"fmt"
	"unsafe"
)

func (term *Term) writeCell(r rune) {
	//debug.Log(term.curPos)

	if term.curPos.X >= term.size.X {
		if term._noAutoLineWrap {
			term.curPos.X--

		} else {
			overflow := term.curPos.X - (term.size.X - 1)
			term.curPos.X = 0

			if overflow > 0 {
				term.lineFeed()
			}
		}
	}

	cell := term.cell()
	cell.Char = r
	cell.Sgr = term.sgr.Copy()

	if term._activeElement != nil {
		cell.Element = term._activeElement
		term._activeElement.ReadCell(cell)
	}

	term.curPos.X++
}

func (term *Term) appendScrollBuf() {
	if unsafe.Pointer(term.cells) == unsafe.Pointer(&term._normBuf) {
		term._scrollBuf = append(term._scrollBuf, term._normBuf[0])
		if term._scrollOffset > 0 {
			term._scrollOffset++
			if term._scrollMsg != nil {
				term._scrollMsg.SetMessage(fmt.Sprintf("Viewing scrollback history. %d lines from end", term._scrollOffset))
				term.renderer.TriggerRedraw()
			}
		}
	}
}
