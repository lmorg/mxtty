package rendersdl

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
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

	sr.setTypeFace(font)

	sr.border = 5

	sr._quit = make(chan bool)
	sr._event = make(chan bool)
	sr._redraw = make(chan bool)

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

	sr.surface, err = sr.window.GetSurface()
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
	c := virtualterm.SGR_COLOUR_BLACK
	pixel := sdl.MapRGBA(sr.surface.Format, c.Red, c.Green, c.Blue, 255)
	err := sr.surface.FillRect(&sdl.Rect{W: sr.surface.W, H: sr.surface.H}, pixel)
	if err != nil {
		log.Printf("error drawing background: %s", err.Error())
	}

	for {
		slowPoll := time.After(250 * time.Millisecond)

	checkEvent:
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {

			case *sdl.QuitEvent:
				go sr.triggerQuit()

			case *sdl.WindowEvent:
				eventWindow(sr, evt, term)
				go sr.triggerEvent()

			case *sdl.TextInputEvent:
				eventTextInput(evt, term)
				go sr.triggerEvent()

			case *sdl.KeyboardEvent:
				eventKeyPress(evt, term)
				go sr.triggerEvent()
			}
		}

		select {
		case <-sr._quit:
			return

		case <-sr._event:
			update(sr, term)

		case <-slowPoll:
			update(sr, term)

		case <-time.After(15 * time.Millisecond):
			goto checkEvent
		}

		//sdl.Delay(250)

	}
}

func update(sr *sdlRender, term types.Term) {
	term.Render()

	if atomic.CompareAndSwapInt32(&sr.updateTitle, 1, 0) {
		sr.window.SetTitle(sr.title)
	}

	err := sr.window.UpdateSurface()
	if err != nil {
		log.Printf("error in renderer: %s", err.Error())
	}
}
