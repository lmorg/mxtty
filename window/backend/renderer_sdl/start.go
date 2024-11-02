package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"golang.design/x/clipboard"
)

/*
	Reference documentation used:
	- https://github.com/veandco/go-sdl2-examples/tree/master/examples
*/

var (
	width  int32 = 1024
	height int32 = 768
	X      int32 = sdl.WINDOWPOS_UNDEFINED
	Y      int32 = sdl.WINDOWPOS_UNDEFINED
)

func Initialise() types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	rect, err := sdl.GetDisplayUsableBounds(0)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	} else {
		X = rect.W / 2
		Y = 0
		width = rect.W / 2
		height = rect.H
	}

	sr := new(sdlRender)
	err = sr.createWindow(app.Title)
	if err != nil {
		panic(err.Error())
	}

	sr.border = 5

	sr._quit = make(chan bool)
	sr._redraw = make(chan bool)

	font, err := typeface.Open(
		config.Config.Terminal.TypeFace.FontName,
		config.Config.Terminal.TypeFace.FontSize,
	)
	if err != nil {
		panic(err.Error())
	}
	sr.setTypeFace(font)

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	return sr
}

func (sr *sdlRender) createWindow(caption string) error {
	var err error

	// Create a window for us to draw the text on
	sr.window, err = sdl.CreateWindow(
		caption,             // window title
		X, Y, width, height, // window position & dimensions
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_HIDDEN, // window properties
	)
	if err != nil {
		return err
	}

	sr.window.SetWindowOpacity(float32(config.Config.Window.Opacity) / 100)

	err = sr.setIcon()
	if err != nil {
		return err
	}

	sr.initBell()

	sr.renderer, err = sdl.CreateRenderer(sr.window, -1, sdl.RENDERER_ACCELERATED) //|sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		return err
	}

	err = sr.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	sr.ShowAndFocusWindow()
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
	sr.preloadNotificationGlyphs()
}

func (sr *sdlRender) getTermSizeCells() *types.XY {
	x, y := sr.window.GetSize()

	return &types.XY{
		X: ((x - (sr.border * 2)) / sr.glyphSize.X),
		Y: ((y - (sr.border * 2)) / sr.glyphSize.Y),
	}
}
