package rendersdl

import (
	"github.com/lmorg/mxtty/virtualterm/types"
	"github.com/veandco/go-sdl2/sdl"
)

func printRuneColour(r rune, posX, posY int32, fg *types.Colour, bg *types.Colour) error {
	//log.Printf("debug: r %d pos %d:%d, fg: %v, bg %v", r, posX, posY, *fg, *bg)
	rect := &sdl.Rect{
		X: (glyphSize.X * posX) + border,
		Y: (glyphSize.Y * posY) + border,
		W: glyphSize.X,
		H: glyphSize.Y,
	}

	//text, err := font.RenderGlyphSolid(r, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	text, err := font.RenderGlyphSolid(r, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	pixel := sdl.MapRGBA(surface.Format, bg.Red, bg.Green, bg.Blue, 255)
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

func printBlink(state bool, posX, posY int32) error {
	text, err := font.RenderGlyphSolid('_', blinkColour[state])
	if err != nil {
		return err
	}
	defer text.Free()

	rect := &sdl.Rect{
		X: (glyphSize.X * posX) + border,
		Y: (glyphSize.Y * posY) + border,
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

func update() error {
	return window.UpdateSurface()
}
