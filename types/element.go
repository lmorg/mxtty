package types

type Element interface {
	Begin(*ApcSlice)
	ReadCell(*Cell)
	End() *XY             // return is optional
	Insert(*ApcSlice) *XY // return is required
	Draw(*Rect) *XY       // return is optional
	MouseClick(uint8, *XY)
	Close()
}

type ElementID int

const (
	ELEMENT_ID_IMAGE ElementID = iota
	ELEMENT_ID_TABLE
	ELEMENT_ID_FOLDABLE
)
