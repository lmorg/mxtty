package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func (sr *sdlRender) PrintRuneColour(r rune, posX, posY int32, fg *types.Colour, bg *types.Colour, style types.SgrFlag) error {
	//log.Printf("debug: r %d pos %d:%d, fg: %v, bg %v", r, posX, posY, *fg, *bg)
	rect := &sdl.Rect{
		X: (sr.glyphSize.X * posX) + sr.border,
		Y: (sr.glyphSize.Y * posY) + sr.border,
		W: sr.glyphSize.X,
		H: sr.glyphSize.Y,
	}

	/*rect2 := &sdl.Rect{
		X: (sr.glyphSize.X * posX) + sr.border + 1,
		Y: (sr.glyphSize.Y * posY) + sr.border + 1,
		W: sr.glyphSize.X,
		H: sr.glyphSize.Y,
	}*/

	sr.font.SetStyle(fontStyle(style))

	text, err := sr.font.RenderGlyphBlended(r, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	/*sr.font.SetStyle(ttf.STYLE_BOLD | fontStyle(style))

	text2, err := sr.font.RenderGlyphBlended(r, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		return err
	}
	defer text2.Free()*/

	pixel := sdl.MapRGBA(sr.surface.Format, bg.Red, bg.Green, bg.Blue, 1)
	err = sr.surface.FillRect(rect, pixel)
	if err != nil {
		return err
	}

	/*err = text2.Blit(nil, sr.surface, rect2)
	if err != nil {
		return err
	}*/

	err = text.Blit(nil, sr.surface, rect)
	if err != nil {
		return err
	}

	return nil
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
