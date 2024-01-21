package virtualterm

/*
	RESET TAB STOPS
*/

func (term *Term) csiResetTabStops() {
	term._tabStops = make([]int32, 0)
	for i := int32(0); i < term.size.X; i += term._tabWidth {
		term._tabStops = append(term._tabStops, i)
	}
}

func (term *Term) nextTabStop() int32 {
	for _, tabStop := range term._tabStops {
		if tabStop > term.curPos.X {
			return tabStop - term.curPos.X
		}
	}

	// end of screen
	return term.size.X - term.curPos.X
}

func (term *Term) printTab() {
	term._printTab(term.nextTabStop())
}

func (term *Term) _printTab(tabWidth int32) {
	for i := 1; i < int(tabWidth); i++ {
		term.writeCell(' ')
	}
}

func (term *Term) csiClearTab() {
	term.csiEraseCharacters(term.nextTabStop())
}
