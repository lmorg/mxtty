package types

type Element interface {
	Generate(*ApcSlice) error
	Write(rune) error
	Rune(*XY) rune
	Size() *XY
	Draw(*XY, *XY)
	MouseClick(uint8, *XY, MouseClickCallback)
	Close()
}

type ElementID int

const (
	ELEMENT_ID_IMAGE ElementID = iota
	ELEMENT_ID_CSV
	ELEMENT_ID_FOLDABLE
)
