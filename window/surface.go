package window

import (
	"github.com/lmorg/mxtty/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window    *sdl.Window
	surface   *sdl.Surface
	font      *ttf.Font
	glyphSize *typeface.SizeT
	border    int32 = 5
	width     int32 = 1024
	height    int32 = 768
)

func init() {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}
}

func Create(caption string) error {
	var err error

	// Create a window for us to draw the text on
	window, err = sdl.CreateWindow(caption, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	surface, err = window.GetSurface()
	return err
}

func SetTypeFace(f *ttf.Font) (int32, int32) {
	font = f
	glyphSize = typeface.GetSize()
	w, h := window.GetSize()

	xCells := (w - (border * 2)) / glyphSize.Width
	yCells := (h - (border * 2)) / glyphSize.Height

	return xCells, yCells
}

func Update() error {
	return window.UpdateSurface()
}

func Close() {
	window.Destroy()
	sdl.Quit()
}

/*func PrintRune(r rune, posX, posY int32) error {
	rect := &sdl.Rect{
		X: (glyphSize.Width * posX) + border,
		Y: (glyphSize.Height * posY) + border,
		W: glyphSize.Width,
		H: glyphSize.Height,
	}

	text, err := font.RenderGlyphSolid(r, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	err = surface.FillRect(rect, 0)
	if err != nil {
		return err
	}
	err = text.Blit(nil, surface, rect)
	if err != nil {
		return err
	}

	return nil
}*/

func PrintRuneColour(r rune, posX, posY int32, fg *sdl.Color, bg *sdl.Color) error {
	rect := &sdl.Rect{
		X: (glyphSize.Width * posX) + border,
		Y: (glyphSize.Height * posY) + border,
		W: glyphSize.Width,
		H: glyphSize.Height,
	}

	text, err := font.RenderGlyphSolid(r, *fg)
	if err != nil {
		return err
	}
	defer text.Free()

	pixel := sdl.MapRGBA(surface.Format, bg.R, bg.G, bg.B, bg.A)
	err = surface.FillRect(rect, pixel)
	if err != nil {
		return err
	}
	err = text.Blit(nil, surface, rect)
	if err != nil {
		return err
	}

	return nil
}

var blinkColour = map[bool]sdl.Color{
	true:  {R: 255, G: 255, B: 255, A: 255},
	false: {R: 0, G: 0, B: 0, A: 255},
}

func PrintBlink(state bool, posX, posY int32) error {
	text, err := font.RenderGlyphSolid('_', blinkColour[state])
	if err != nil {
		return err
	}
	defer text.Free()

	//rect := &sdl.Rect{X: 400 - (text.W / 2), Y: 300 - (text.H / 2), W: 0, H: 0}
	rect := &sdl.Rect{
		X: (glyphSize.Width * posX) + border,
		Y: (glyphSize.Height * posY) + border,
		//W: glyphSize.Width,
		//H: glyphSize.Height,
	}

	// Draw the text around the center of the window
	err = text.Blit(nil, surface, rect)
	if err != nil {
		return err
	}

	return nil
}

/*
func PrintTextFont(font *ttf.Font, s string) error {
	// Create a red text with the font
	text, err := font.RenderUTF8Solid(s, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	//rect := &sdl.Rect{X: 400 - (text.W / 2), Y: 300 - (text.H / 2), W: 0, H: 0}
	rect := &sdl.Rect{X: 5, Y: 5}

	// Draw the text around the center of the window
	err = text.Blit(nil, surface, rect)
	if err != nil {
		return err
	}

	return nil
}
*/
