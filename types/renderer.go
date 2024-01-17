package types

type NotificationType int

const (
	NOTIFY_INFO = iota
	NOTIFY_WARNING
	NOTIFY_ERROR
)

type Renderer interface {
	Start(Term)
	TermSize() *XY
	Resize() *XY
	PrintCell(*Cell, *XY) error
	GetWindowTitle() string
	SetWindowTitle(string)
	Bell()
	TriggerRedraw()
	NewElement(elementType ElementID, size *XY, data []byte) Element
	Close()
}

type Image interface {
	Size() *XY
	Draw(size *XY, rect *Rect)
	Close()
}
