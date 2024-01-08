package rendersdl

import (
	"log"
	"sync/atomic"

	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	border int32 = 5
	width  int32 = 1024
	height int32 = 768
)

var focused *sdlRender

func Initialise(fontName string, fontSize int) types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	focused = new(sdlRender)
	err = focused.createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open(fontName, fontSize)
	if err != nil {
		panic(err.Error())
	}

	focused.setTypeFace(font)

	return focused
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

	rwops, err := sdl.RWFromMem(assets.Get(assets.ICON_BMP))
	if err != nil {
		return err
	}

	icon, err := sdl.LoadBMPRW(rwops, true)
	if err != nil {
		return err
	}

	sr.window.SetIcon(icon)

	sr.surface, err = sr.window.GetSurface()
	return err
}

func (sr *sdlRender) setTypeFace(f *ttf.Font) {
	sr.font = f
	sr.glyphSize = typeface.GetSize()
	sr.termSize = sr.getTermSize()
}

func (sr *sdlRender) getTermSize() *types.Rect {
	x, y := sr.window.GetSize()

	return &types.Rect{
		X: (x - (border * 2)) / sr.glyphSize.X,
		Y: (y - (border * 2)) / sr.glyphSize.Y,
	}
}

func Start(r types.Renderer, term *virtualterm.Term) {
	c := virtualterm.SGR_COLOUR_BLACK
	pixel := sdl.MapRGBA(focused.surface.Format, c.Red, c.Green, c.Blue, 255)
	err := focused.surface.FillRect(&sdl.Rect{W: focused.surface.W, H: focused.surface.H}, pixel)
	if err != nil {
		log.Printf("error drawing background: %s", err.Error())
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {

			case *sdl.QuitEvent:
				running = false

			case *sdl.WindowEvent:
				eventWindow(r, evt, term)

			case *sdl.TextInputEvent:
				eventTextInput(evt, term)

			case *sdl.KeyboardEvent:
				eventKeyPress(evt, term)

			}
		}

		sdl.Delay(5)
		term.Render()

		if atomic.CompareAndSwapInt32(&focused.updateTitle, 1, 0) {
			focused.window.SetTitle(focused.title)
		}

		err = focused.window.UpdateSurface()
		if err != nil {
			log.Printf("error in renderer: %s", err.Error())
		}
	}
}
