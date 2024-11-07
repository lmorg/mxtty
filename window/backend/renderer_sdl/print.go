package rendersdl

import (
	"unsafe"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const dropShadowOffset = 2

var (
	textShadow    = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	textHighlight = sdl.Color{R: 50, G: 50, B: 255, A: 255}
)

func (sr *sdlRender) PrintCell(cell *types.Cell, cellPos *types.XY) error {
	fg, bg := sgrOpts(cell.Sgr)

	sr.setFontStyle(cell.Sgr.Bitwise)
	r := cell.Char
	if r == 0 {
		r = ' '
	}

	rect := &sdl.Rect{
		X: (sr.glyphSize.X * cellPos.X) + sr.border,
		Y: (sr.glyphSize.Y * cellPos.Y) + sr.border,
		W: sr.glyphSize.X,
		H: sr.glyphSize.Y,
	}

	isCellHighlighted := isCellHighlighted(sr, rect)

	// render background colour

	if bg != nil {
		var pixel uint32
		if isCellHighlighted {
			pixel = sdl.MapRGBA(sr.surface.Format, textHighlight.R, textHighlight.G, textHighlight.B, 255)
		} else {
			pixel = sdl.MapRGBA(sr.surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		}
		err := sr.surface.FillRect(rect, pixel)
		if err != nil {
			return err
		}
	}

	// render drop shadow

	if config.Config.Terminal.TypeFace.DropShadow && (bg == nil || isCellHighlighted) {
		rect2 := &sdl.Rect{
			X: (sr.glyphSize.X * cellPos.X) + sr.border + dropShadowOffset,
			Y: (sr.glyphSize.Y * cellPos.Y) + sr.border + dropShadowOffset,
			W: sr.glyphSize.X,
			H: sr.glyphSize.Y,
		}

		var c sdl.Color
		if isCellHighlighted && bg == nil {
			c = textHighlight
		} else {
			c = textShadow
		}
		text2, err := sr.font.RenderGlyphBlended(r, c)
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
	if isCellHighlighted {
		text.SetBlendMode(sdl.BLENDMODE_ADD)
	}

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

func sgrOpts(sgr *types.Sgr) (fg *types.Colour, bg *types.Colour) {
	if sgr.Bitwise.Is(types.SGR_INVERT) {
		bg, fg = sgr.Fg, sgr.Bg
	} else {
		fg, bg = sgr.Fg, sgr.Bg
	}

	if unsafe.Pointer(bg) == unsafe.Pointer(types.SGR_DEFAULT.Bg) {
		bg = nil
	}

	return fg, bg
}
