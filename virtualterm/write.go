package virtualterm

func (term *Term) writeCell(r rune) {
	//debug.Log(term.curPos)

	if term.curPos.X >= term.size.X {
		if term._noAutoLineWrap {
			term.curPos.X--

		} else {
			overflow := term.curPos.X - (term.size.X - 1)
			term.curPos.X = 0

			if overflow > 0 && term.csiMoveCursorDownwards(1) > 0 {
				term.csiScrollUp(1)
				term.csiMoveCursorDownwards(1)
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

	/*if term.curPos.X >= term.size.X {
		if term._noAutoLineWrap {
			term.curPos.X--

		} else {
			overflow := term.curPos.X - (term.size.X - 1)
			term.curPos.X = 0

			if overflow > 0 && term.csiMoveCursorDownwards(1) > 0 {
				term.csiScrollUp(1)
				term.csiMoveCursorDownwards(1)
			}
		}
	}*/
}
