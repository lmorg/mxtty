package rendersdl

import (
	"strings"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

// PrintCellBlock is much slower because it doesn't cache textures
func (sr *sdlRender) PrintCellBlock(cells []types.Cell, cellPos *types.XY) {
	if len(cells) == 0 {
		return
	}

	r := make([]rune, len(cells))
	for i := range cells {
		r[i] = cells[i].Rune()
	}
	s := strings.TrimRight(string(r), " ")
	if s == "" {
		return
	}

	surface := _newFontSurface(sr.glyphSize, int32(len(cells)))
	defer surface.Free()

	sgr := cells[0].Sgr
	if sgr == nil {
		sgr = types.SGR_DEFAULT
	}

	sr.font.SetStyle(fontStyle(sgr.Bitwise))

	fg, bg := sgrOpts(sgr)

	cellBlockRect := &sdl.Rect{
		W: sr.glyphSize.X * int32(len(cells)),
		H: sr.glyphSize.Y,
	}

	// render background colour

	if bg != nil {
		var pixel uint32
		//if isCellHighlighted {
		//	pixel = sdl.MapRGBA(surface.Format, textHighlight.R, textHighlight.G, textHighlight.B, 255)
		//} else {
		pixel = sdl.MapRGBA(surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		//}
		err := surface.FillRect(cellBlockRect, pixel)
		if err != nil {
			panic(err) // TODO: better error handling please!
		}
	}

	if config.Config.Terminal.TypeFace.DropShadow && bg == nil { // (bg == nil) || isCellHighlighted) {
		shadowRect := &sdl.Rect{
			X: cellBlockRect.X + dropShadowOffset,
			Y: cellBlockRect.Y + dropShadowOffset,
			W: cellBlockRect.W,
			H: cellBlockRect.H,
		}

		var c sdl.Color
		//if isCellHighlighted && bg == nil {
		//	c = textHighlight
		//} else {
		c = textShadow
		//}
		//shadowText, err := font.RenderGlyphBlended(cell.Char, c)
		shadowText, err := sr.font.RenderUTF8Blended(s, c)
		if err != nil {
			panic(err) // TODO: better error handling please!
		}
		defer shadowText.Free()

		err = shadowText.Blit(nil, surface, shadowRect)
		if err != nil {
			panic(err) // TODO: better error handling please!
		}
	}

	// render cell char
	//text, err := font.RenderGlyphBlended(cell.Char, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	text, err := sr.font.RenderUTF8Blended(s, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	defer text.Free()
	//if isCellHighlighted {
	//	text.SetBlendMode(sdl.BLENDMODE_ADD)
	//}

	err = text.Blit(nil, surface, cellBlockRect)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	if sgr.Bitwise.Is(types.SGR_BOLD) {
		text.SetBlendMode(sdl.BLENDMODE_ADD)
		err = text.Blit(nil, surface, cellBlockRect)
		if err != nil {
			panic(err) // TODO: better error handling please!
		}
	}

	texture, err := sr.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	defer texture.Destroy()

	dstRect := &sdl.Rect{
		X: (sr.glyphSize.X * cellPos.X) + sr.border,
		Y: (sr.glyphSize.Y * cellPos.Y) + sr.border,
		W: cellBlockRect.W + dropShadowOffset,
		H: cellBlockRect.H + dropShadowOffset,
	}

	sr.renderer.Copy(texture, cellBlockRect, dstRect)
}
