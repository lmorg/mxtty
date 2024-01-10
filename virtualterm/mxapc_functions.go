package virtualterm

type mxapc interface {
	//Append(*cell)
	//CrLf()
}

type cellCollection struct {
	value string
	cells []*cell
	width int32
}

type mxapcTable struct {
	headingOffset int
	table         map[string]cellCollection
}

func (table *mxapcTable) Append(cell *cell) {

}

func (term *Term) mxapcTableBegin(parameters apcSlice) {

}

func (term *Term) mxapcTableEnd(parameters apcSlice) {

}
