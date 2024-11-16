package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/tmux"
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

func Initialise() (types.Renderer, *types.XY) {
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
	if config.Config.Window.StatusBar {
		sr.footer++
	}
	if config.Config.Tmux.Enabled {
		sr.footer++
	}

	sr._quit = make(chan bool)
	sr._redraw = make(chan bool)
	sr.ligCache = newCachedLigatures(sr)

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

	sr.preloadNotificationGlyphs()
	sr.fontCache = NewFontCache(sr)

	return sr, sr._getSizeCells()
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

	setLghtOrDarkMode()

	sr.renderer, err = sdl.CreateRenderer(sr.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	err = sr.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return err
	}

	sr.ShowAndFocusWindow()
	return nil
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
}

func setLghtOrDarkMode() {
	if config.Config.Terminal.LightMode {
		highlightBlendMode = sdl.BLENDMODE_BLEND
		textShadow.A = 32
		types.SGR_DEFAULT.Fg, types.SGR_DEFAULT.Bg = types.SGR_DEFAULT.Bg, types.SGR_DEFAULT.Fg
	} else {
		highlightBlendMode = sdl.BLENDMODE_ADD
		textShadow.A = 255
	}
}

func (sr *sdlRender) Start(term types.Term, tmuxClient any) {
	sr.term = term
	sr.tmux = tmuxClient.(*tmux.Tmux)

	sr.registerHotkey()
	go sr.refreshInterval()

	sr.eventLoop()
}
