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
			sr.eventHotkey()

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

	err := sr.renderer.SetDrawColor(bg.Red, bg.Green, bg.Blue, 255)
	if err != nil {
		log.Printf("ERROR: error drawing background: %s", err.Error())
	}
	err = sr.renderer.FillRect(rect)
	if err != nil {
		log.Printf("ERROR: error drawing background: %s", err.Error())
	}
}

func render(sr *sdlRender, term types.Term) error {
	x, y := sr.window.GetSize()
	rect := &sdl.Rect{W: x, H: y}

	if sr.highlighter != nil && sr.highlighter.button == 0 {
		texture, err := sr.renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, x, y)
		if err != nil {
			sr.highlighter = nil
			return err
		}
		defer texture.Destroy()
		err = sr.renderer.SetRenderTarget(texture)
		if err != nil {
			sr.highlighter = nil
			return err
		}
	}

	sr.drawBg(term, rect)

	term.Render()

	for i := range sr.fnStack {
		sr.fnStack[i]()
	}
	sr.fnStack = make([]func(), 0) // clear image stack

	if sr.highlighter != nil && sr.highlighter.button == 0 {
		sr.copyRendererToClipboard()
		return nil
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
