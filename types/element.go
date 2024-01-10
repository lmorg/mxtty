package types

type Element interface {
	Draw()
	ReadCell(*Cell)
	Close()
}

type ElementID int

const (
	ELEMENT_ID_IMAGE ElementID = iota
	ELEMENT_ID_TABLE
	ELEMENT_ID_FOLDABLE
)

type ELEMENT_TABLE interface {
	Element
	//EventCallback()
}
