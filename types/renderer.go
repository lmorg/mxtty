package types

type Renderer interface {
	Start(Term)
	FocusWindow()
	TermSize() *XY
	PrintCell(*Cell, *XY) error
	GetWindowTitle() string
	SetWindowTitle(string)
	Bell()
	TriggerRedraw()
	TriggerQuit()
	NewElement(elementType ElementID) Element
	DisplayNotification(NotificationType, string)
	DisplaySticky(NotificationType, string) Notification
	DisplayInputBox(string, string, func(string))
	AddRenderFnToStack(func())
	GetWindowMeta() any
	ResizeWindow(*XY)
	Close()
}

type Image interface {
	Size() *XY
	Asset() any
	Draw(size *XY, pos *XY)
	Close()
}
