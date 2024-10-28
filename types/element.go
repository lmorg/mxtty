package types

type Element interface {
	Generate(*ApcSlice) error
	Size() *XY
	Draw(*XY, *XY)
	MouseClick(uint8, *XY)
	Close()
}

type ElementID int

const (
	ELEMENT_ID_IMAGE ElementID = iota
	ELEMENT_ID_TABLE
	ELEMENT_ID_FOLDABLE
)
