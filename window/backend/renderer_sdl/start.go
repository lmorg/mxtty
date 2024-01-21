package rendersdl

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	width  int32 = 1024
	height int32 = 768
)

func Initialise(fontName string, fontSize int) types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	sr := new(sdlRender)
	err = sr.createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open(fontName, fontSize)
	if err != nil {
		panic(err.Error())
	}

	sr._quit = make(chan bool)
	sr._redraw = make(chan bool)

	sr.setTypeFace(font)
	sr.border = 5
	sr.dropShadow = true

	return sr
}

func (sr *sdlRender) createWindow(caption string) error {
	var err error

	// Create a window for us to draw the text on
	sr.window, err = sdl.CreateWindow(
		caption,                                          // window title
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, // window pos
		width, height, // window dimensions
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE, // window properties
	)
	if err != nil {
		return err
	}

	err = sr.setIcon()
	if err != nil {
		return err
	}

	sr.initBell()

	sr.renderer, err = sdl.CreateRenderer(sr.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	err = sr.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	return err
}

func (sr *sdlRender) setIcon() error {
	rwops, err := sdl.RWFromMem(assets.Get(assets.ICON_APP))
	if err != nil {
		return err
	}

	icon, err := sdl.LoadBMPRW(rwops, true)
	if err != nil {
		return err
	}

	sr.window.SetIcon(icon)

	return nil
}

func (sr *sdlRender) setTypeFace(f *ttf.Font) {
	sr.font = f
	sr.glyphSize = typeface.GetSize()
	sr.termSize = sr.getTermSize()
}

func (sr *sdlRender) getTermSize() *types.XY {
	x, y := sr.window.GetSize()

	return &types.XY{
		X: ((x - (sr.border * 2)) / sr.glyphSize.X),
		Y: ((y - (sr.border * 2)) / sr.glyphSize.Y),
	}
}

func (sr *sdlRender) Start(term types.Term) {
	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {

			case *sdl.QuitEvent:
				go sr.triggerQuit()

			case *sdl.WindowEvent:
				eventWindow(sr, evt, term)
				go sr.TriggerRedraw()

			case *sdl.TextInputEvent:
				eventTextInput(evt, term)
				go sr.TriggerRedraw()

			case *sdl.KeyboardEvent:
				eventKeyPress(evt, term)
				go sr.TriggerRedraw()

			case *sdl.MouseButtonEvent:
				eventMouseButton(evt, term, sr)
				go sr.TriggerRedraw()
			}
		}

		select {
		case <-sr._quit:
			return

		case <-sr._redraw:
			update(sr, term)

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

func update(sr *sdlRender, term types.Term) {
	var err error
	x, y := sr.window.GetSize()
	rect := &sdl.Rect{W: x, H: y}

	sr.surface, err = sdl.CreateRGBSurfaceWithFormat(0, x, y, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) //TODO: don't panic!
	}
	defer sr.surface.Free()

	sr.drawBg(term, rect)

	term.Render()

	texture, err := sr.renderer.CreateTextureFromSurface(sr.surface)
	if err != nil {
		panic(err) //TODO: don't panic!
	}

	err = sr.renderer.Copy(texture, rect, rect)
	if err != nil {
		panic(err) //TODO: don't panic!
	}

	for i := range sr.fnStack {
		sr.fnStack[i]()
	}
	sr.fnStack = make([]func(), 0) // clear image stack

	sr.renderNotification(rect)

	if atomic.CompareAndSwapInt32(&sr.updateTitle, 1, 0) {
		sr.window.SetTitle(sr.title)
	}

	sr.renderer.Present()
}
