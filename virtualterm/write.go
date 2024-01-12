package virtualterm

import "github.com/lmorg/mxtty/types"

func (term *Term) writeCell(r rune) {
	if term.curPos.X >= term.size.X {
		overflow := term.curPos.X - (term.size.X - 1)
		term.curPos.X = 0

		if overflow > 0 && term.csiMoveCursorDownwards(1) > 0 {
			term.csiScrollUp(1)
			term.csiMoveCursorDownwards(1)
		}
	}

	cell := term.cell()
	cell.Char = r
	cell.Sgr = term.sgr.Copy()

	if term._activeElement != nil {
		term.sgr.Bitwise.Unset(types.APC_BEGIN_ELEMENT)
		cell.Element = term._activeElement
		cell.Sgr.Bitwise.Set(types.APC_ELEMENT)
		term._activeElement.ReadCell(cell) // might need to copy value rather than pointer
	}

	term.curPos.X++
}
