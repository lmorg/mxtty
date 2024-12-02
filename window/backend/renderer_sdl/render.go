package rendersdl

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/renderer_sdl/layer"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) refreshInterval() {
	if config.Config.Window.RefreshInterval == 0 {
		return
	}

	d := time.Duration(config.Config.Window.RefreshInterval) * time.Millisecond
	for {
		time.Sleep(d)
		sr.TriggerRedraw()
	}
}

func (sr *sdlRender) eventLoop() {
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
		case <-sr.pollEventHotkey():
			sr.eventHotkey()

		case <-sr._quit:
			return

		case <-sr._redraw:
			err := render(sr)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			sr.limiter.Unlock()

		case size := <-sr._resize:
			sr._resizeWindow(size)

		case <-time.After(15 * time.Millisecond):
			continue
		}
	}
}

func (sr *sdlRender) drawBg(term types.Term, rect *sdl.Rect) {
	bg := term.Bg()

	texture := sr.createRendererTexture()
	if texture == nil {
		return
	}
	defer sr.restoreRendererTexture()

	var err error

	err = sr.renderer.SetDrawColor(bg.Red, bg.Green, bg.Blue, 255)
	if err != nil {
		log.Printf("ERROR: error drawing background: %v", err)
	}

	err = sr.renderer.FillRect(rect)
	if err != nil {
		log.Printf("ERROR: error drawing background: %v", err)
	}
}

func (sr *sdlRender) AddToElementStack(item *layer.RenderStackT) {
	sr._elementStack = append(sr._elementStack, item)
}

func (sr *sdlRender) AddToOverlayStack(item *layer.RenderStackT) {
	sr._overlayStack = append(sr._overlayStack, item)
}

func (sr *sdlRender) createRendererTexture() *sdl.Texture {
	w, h, err := sr.renderer.GetOutputSize()
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}
	texture, err := sr.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_TARGET, w, h)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}
	err = sr.renderer.SetRenderTarget(texture)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}
	err = texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}
	return texture
}

func (sr *sdlRender) restoreRendererTexture() {
	texture := sr.renderer.GetRenderTarget()
	sr.AddToElementStack(&layer.RenderStackT{texture, nil, nil, true})
	err := sr.renderer.SetRenderTarget(nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func (sr *sdlRender) renderStack(stack *[]*layer.RenderStackT) {
	var err error
	for _, item := range *stack {
		err = sr.renderer.Copy(item.Texture, item.SrcRect, item.DstRect)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		if item.Destroy {
			_ = item.Texture.Destroy()
		}
	}
	*stack = make([]*layer.RenderStackT, 0) // clear image stack
}

func render(sr *sdlRender) error {
	x, y := sr.window.GetSize()
	rect := &sdl.Rect{W: x, H: y}

	sr.drawBg(sr.term, rect)
	sr.term.Render()
	sr.renderFooter()

	if sr.highlighter != nil && sr.highlighter.button == 0 {
		texture := sr.createRendererTexture()
		if texture == nil {
			sr.highlighter = nil
			return nil
		}
		defer texture.Destroy()
	}

	sr.renderStack(&sr._elementStack)

	if sr.highlighter != nil && sr.highlighter.button == 0 {
		sr.copyRendererToClipboard()
		return nil
	}

	switch {
	case sr.inputBox != nil:
		sr.renderInputBox(rect)

	case sr.menu != nil:
		sr.renderMenu(rect)

	default:
		sr.selectionHighlighter()
	}

	sr.renderStack(&sr._overlayStack)

	sr.renderNotification(rect)

	if atomic.CompareAndSwapInt32(&sr.updateTitle, 1, 0) {
		sr.window.SetTitle(sr.title)
	}

	sr.renderer.Present()
	return nil
}
