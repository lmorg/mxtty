package rendersdl

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventLoop(term types.Term) {
	sr.term = term

	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {

			case *sdl.QuitEvent:
				sr.TriggerQuit()

			case *sdl.WindowEvent:
				sr.eventWindow(evt)
				sr.TriggerRedraw()

			case *sdl.TextInputEvent:
				sr.eventTextInput(evt)
				sr.TriggerRedraw()

			case *sdl.KeyboardEvent:
				sr.eventKeyPress(evt)
				sr.TriggerRedraw()

			case *sdl.MouseButtonEvent:
				sr.eventMouseButton(evt)
				sr.TriggerRedraw()

			case *sdl.MouseMotionEvent:
				sr.eventMouseMotion(evt)
				// don't trigger redraw

			case *sdl.MouseWheelEvent:
				sr.eventMouseWheel(evt)
				sr.TriggerRedraw()
			}
		}

		select {
		case <-sr.hk.Keydown():
			if sr.hkToggle {
				sr.window.Hide()
			} else {
				sr.FocusWindow()
			}
			sr.hkToggle = !sr.hkToggle

		case <-sr._quit:
			return

		case <-sr._redraw:
			err := render(sr, term)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			sr.limiter.Unlock()

		case <-time.After(15 * time.Millisecond):
			continue
		}
	}
}

func (sr *sdlRender) drawBg(term types.Term, rect *sdl.Rect) {
	bg := term.Bg()

	pixel := sdl.MapRGBA(sr.surface.Format, bg.Red, bg.Green, bg.Blue, 255)
	err := sr.surface.FillRect(rect, pixel)
	if err != nil {
		log.Printf("ERROR: error drawing background: %s", err.Error())
	}
}

func render(sr *sdlRender, term types.Term) error {
	var err error
	x, y := sr.window.GetSize()
	rect := &sdl.Rect{W: x, H: y}

	sr.surface, err = sdl.CreateRGBSurfaceWithFormat(0, x, y, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		return err
	}
	defer sr.surface.Free()

	sr.drawBg(term, rect)

	term.Render()

	texture, err := sr.renderer.CreateTextureFromSurface(sr.surface)
	if err != nil {
		return err
	}

	err = sr.renderer.Copy(texture, rect, rect)
	if err != nil {
		return err
	}

	for i := range sr.fnStack {
		sr.fnStack[i]()
	}
	sr.fnStack = make([]func(), 0) // clear image stack

	if sr.highlighter != nil && sr.highlighter.button == 0 {
		sr.copySurfaceToClipboard()
	}

	sr.renderNotification(rect)

	if sr.inputBox != nil {
		sr.renderInputBox(rect)
	} else {
		sr.selectionHighlighter()
	}

	if atomic.CompareAndSwapInt32(&sr.updateTitle, 1, 0) {
		sr.window.SetTitle(sr.title)
	}

	sr.renderer.Present()
	return nil
}
