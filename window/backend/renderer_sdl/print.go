package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func (sr *sdlRender) PrintRuneColour(r rune, posX, posY int32, fg *types.Colour, bg *types.Colour, style types.SgrFlag) error {
	//log.Printf("debug: r %d pos %d:%d, fg: %v, bg %v", r, posX, posY, *fg, *bg)

	sr.setFontStyle(style)

	// render background colour

	rect := &sdl.Rect{
		X: (sr.glyphSize.X * posX) + sr.border,
		Y: (sr.glyphSize.Y * posY) + sr.border,
		W: sr.glyphSize.X,
		H: sr.glyphSize.Y,
	}

	if bg != nil {
		pixel := sdl.MapRGBA(sr.surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		err := sr.surface.FillRect(rect, pixel)
		if err != nil {
			return err
		}
	}

	// render drop shadow

	var rect2 *sdl.Rect
	if sr.dropShadow && bg == nil {
		rect2 = &sdl.Rect{
			X: (sr.glyphSize.X * posX) + sr.border + 3,
			Y: (sr.glyphSize.Y * posY) + sr.border + 3,
			W: sr.glyphSize.X,
			H: sr.glyphSize.Y,
		}

		text2, err := sr.font.RenderGlyphBlended(r, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		if err != nil {
			return err
		}
		defer text2.Free()

		err = text2.Blit(nil, sr.surface, rect2)
		if err != nil {
			return err
		}
	}

	// render cell char

	text, err := sr.font.RenderGlyphBlended(r, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	err = text.Blit(nil, sr.surface, rect)
	if err != nil {
		return err
	}

	return nil
}

func (sr *sdlRender) setFontStyle(style types.SgrFlag) {
	if style == sr._fontStyle {
		return
	}

	sr.font.SetStyle(fontStyle(style))
	sr._fontStyle = style
}

func fontStyle(style types.SgrFlag) int {
	var i int

	if style.Is(types.SGR_BOLD) {
		i |= ttf.STYLE_BOLD
	}

	if style.Is(types.SGR_ITALIC) {
		i |= ttf.STYLE_ITALIC
	}

	if style.Is(types.SGR_UNDERLINE) {
		i |= ttf.STYLE_UNDERLINE
	}

	if style.Is(types.SGR_STRIKETHROUGH) {
		i |= ttf.STYLE_STRIKETHROUGH
	}

	return i
}
