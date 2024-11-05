package types

type Renderer interface {
	Start(Term)
	ShowAndFocusWindow()
	GetTermSize() *XY
	PrintCell(*Cell, *XY) error
	DrawTable(*XY, int32, []int32)
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
	SetKeyboardFnMode(KeyboardMode)
	Close()
}

type Image interface {
	Size() *XY
	Asset() any
	Draw(size *XY, pos *XY)
	Close()
}
