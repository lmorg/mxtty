package virtualterm

func (term *Term) mxapcTableBegin(parameters apcSlice) {
	//el := term.renderer.NewElement(types.ELEMENT_ID_TABLE, nil, nil)
}

func (term *Term) mxapcTableEnd(parameters apcSlice) {
	if term._activeElement == nil {
		return
	}
}
