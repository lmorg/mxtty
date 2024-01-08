package rendersdl

import (
	"github.com/lmorg/mxtty/virtualterm/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window    *sdl.Window
	surface   *sdl.Surface
	font      *ttf.Font
	glyphSize *types.Rect
	border    int32 = 5
	width     int32 = 1024
	height    int32 = 768
)

func Start() *types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	err = createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open("hasklig.ttf", 14)
	//font, err := typeface.Open("monaco.ttf", 16)
	if err != nil {
		panic(err.Error())
	}

	return &types.Renderer{
		Close:          close,
		Size:           setTypeFace(font),
		PrintRuneColor: printRuneColour,
		PrintBlink:     printBlink,
		Update:         update,
	}
}

func createWindow(caption string) error {
	var err error

	// Create a window for us to draw the text on
	window, err = sdl.CreateWindow(caption, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	surface, err = window.GetSurface()
	return err
}

func setTypeFace(f *ttf.Font) *types.Rect {
	font = f
	glyphSize = typeface.GetSize()
	x, y := window.GetSize()

	return &types.Rect{
		X: (x - (border * 2)) / glyphSize.X,
		Y: (y - (border * 2)) / glyphSize.Y,
	}
}

func close() {
	typeface.Close()
	window.Destroy()
	sdl.Quit()
}
