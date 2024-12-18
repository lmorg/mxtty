package virtualterm

func (term *Term) phraseAppend(r rune) {
	if term.IsAltBuf() {
		return
	}

	*term._rowPhrase = append(*term._rowPhrase, r)
}

/*func (term *Term) phraseErasePhrase() {
	if term.IsAltBuf() {
		return
	}

	term._rowPhrase = &[]rune{}
}*/

func (term *Term) phraseSetToRowPos() {
	if term.IsAltBuf() {
		return
	}

	term._rowPhrase = (*term.screen)[term.curPos().Y].Phrase
}
