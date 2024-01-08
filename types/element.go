package types

type Element interface {
	Size() *XY
	Draw(*XY)
	Close()
}
