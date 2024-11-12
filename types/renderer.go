package types

type Renderer interface {
	Start(Term)
	ShowAndFocusWindow()
	GetTermSize() *XY
	GetGlyphSize() *XY
	PrintCell(*Cell, *XY)
	PrintCellBlock([]Cell, *XY)
	DrawTable(*XY, int32, []int32)
	DrawHighlightRect(*XY, *XY)
	GetWindowTitle() string
	SetWindowTitle(string)
	StatusbarText(string)
	Bell()
	TriggerRedraw()
	TriggerQuit()
	NewElement(elementType ElementID) Element
	DisplayNotification(NotificationType, string)
	DisplaySticky(NotificationType, string) Notification
	DisplayInputBox(string, string, func(string))
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
