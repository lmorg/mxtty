package virtualterm

import "sort"

/*
	TAB STOPS
*/

func (term *Term) c1AddTabStop() {
	term._tabStops = append(term._tabStops, term.curPos.X)
	sort.Slice(term._tabStops, func(i, j int) bool {
		return term._tabStops[i] < term._tabStops[j]
	})
}

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
	term._printTab(term.nextTabStop() - 1)
}

func (term *Term) _printTab(tabWidth int32) {
	for i := 1; i < int(tabWidth); i++ {
		term.writeCell(' ')
	}
}

func (term *Term) csiClearTab() {
	term.csiEraseCharacters(term.nextTabStop())
}
