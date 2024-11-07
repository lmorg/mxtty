package rendersdl

import (
	"sync/atomic"

	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) SetWindowTitle(title string) {
	sr.title = title
	atomic.CompareAndSwapInt32(&sr.updateTitle, 0, 1)
}

func (sr *sdlRender) GetWindowTitle() string {
	return sr.window.GetTitle()
}

func (sr *sdlRender) GetWindowMeta() any {
	return sr.window
}

func (sr *sdlRender) ShowAndFocusWindow() {
	defer sr.window.Raise()
	defer sr.window.Show()

	displayNum := screenUnderCursor()
	if displayNum == -1 {
		return
	}
	displayBounds, err := sdl.GetDisplayUsableBounds(displayNum)
	if err != nil {
		return
	}

	winW, _ := sr.window.GetSize()

	posX := displayBounds.W - winW
	if width < 0 {
		winW, posX = displayBounds.W, 0
	}
	sr.window.SetPosition(posX, displayBounds.Y)
	sr.window.SetSize(winW, displayBounds.H)
}

func (sr *sdlRender) hideWindow() {
	sr.window.Hide()
}

func screenUnderCursor() int {
	displayCount, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return -1
	}

	x, y, _ := sdl.GetGlobalMouseState()
	for i := 0; i < displayCount; i++ {
		displayBounds, err := sdl.GetDisplayBounds(i)
		if err != nil {
			return -1
		}

		if x >= displayBounds.X && x <= displayBounds.W &&
			y >= displayBounds.Y && y <= displayBounds.H {
			return i
		}
	}

	return -1
}
