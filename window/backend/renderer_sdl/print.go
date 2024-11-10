package rendersdl

import (
	"strings"
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

var (
	//fontAtlas                 fontAtlasT
	fontAtlasCharacterPreload = []string{
		" ",                          // whitespace
		"1234567890",                 // numeric
		"abcdefghijklmnopqrstuvwxyz", // alpha, lower
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ", // alpha, upper
		"`",                          // backtick
		`!"£$%^&*()-=_+`,             // special, top row
		`[]{};'#:@~\|,./<>?`,         // special others
		`»`,                          // murex
	}
)

type fontCacheDefaultLookupT map[rune]*sdl.Rect
type fontAtlasT struct {
	sgrHash   uint64
	lookup    fontCacheDefaultLookupT
	normal    *sdl.Texture
	highlight *sdl.Texture
}
type fontTextureLookupTableT map[uint64][]*fontAtlasT
type fontCacheT struct {
	atlas    *fontAtlasT
	extended fontTextureLookupTableT
}

func NewFontCache(sr *sdlRender) *fontCacheT {
	chars := []rune(strings.Join(fontAtlasCharacterPreload, ""))

	fc := &fontCacheT{
		atlas:    newFontAtlas(chars, types.SGR_DEFAULT, sr.glyphSize, sr.renderer, sr.font),
		extended: make(fontTextureLookupTableT),
	}

	return fc
}

func newFontAtlas(chars []rune, sgr *types.Sgr, glyphSize *types.XY, renderer *sdl.Renderer, font *ttf.Font) *fontAtlasT {
	glyphSizePlusShadow := &types.XY{
		X: glyphSize.X + dropShadowOffset,
		Y: glyphSize.Y + dropShadowOffset,
	}

	return &fontAtlasT{
		sgrHash:   sgr.HashValue(),
		lookup:    _newFontCacheDefaultLookup(chars, glyphSizePlusShadow),
		normal:    _newFontTexture(chars, sgr, glyphSizePlusShadow, renderer, font, false),
		highlight: _newFontTexture(chars, sgr, glyphSizePlusShadow, renderer, font, true),
	}
}

func _newFontCacheDefaultLookup(chars []rune, glyphSize *types.XY) fontCacheDefaultLookupT {
	m := make(fontCacheDefaultLookupT)
	for i, r := range chars {
		m[r] = &sdl.Rect{
			X: int32(i) * glyphSize.X,
			Y: 0,
			W: glyphSize.X,
			H: glyphSize.Y,
		}
	}
	return m
}

func _newFontTexture(chars []rune, sgr *types.Sgr, glyphSize *types.XY, renderer *sdl.Renderer, font *ttf.Font, isCellHighlighted bool) *sdl.Texture {
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, glyphSize.X*int32(len(chars)), glyphSize.Y, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	defer surface.Free()

	pixel := sdl.MapRGBA(surface.Format, types.SGR_DEFAULT.Bg.Red, types.SGR_DEFAULT.Bg.Green, types.SGR_DEFAULT.Bg.Blue, 255)
	err = surface.FillRect(&sdl.Rect{W: surface.W, H: surface.H}, pixel)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	err = surface.SetColorKey(true, pixel)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	cell := &types.Cell{
		Sgr: sgr,
	}

	var i int
	cellRect := &sdl.Rect{W: glyphSize.X, H: glyphSize.Y}

	font.SetStyle(fontStyle(cell.Sgr.Bitwise))

	for i, cell.Char = range chars {
		cellRect.X = int32(i) * glyphSize.X
		err = _printCellToSurface(cell, cellRect, font, surface, isCellHighlighted)
		if err != nil {
			panic(err) // TODO: better error handling please!
		}
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	//defer texture.Destroy() // we don't want to destroy!

	return texture
}

func (fa *fontAtlasT) Render(renderer *sdl.Renderer, dstRect *sdl.Rect, r rune, hash uint64, isCellHighlighted bool) (bool, error) {
	if hash != fa.sgrHash {
		return false, nil
	}

	srcRect, ok := fa.lookup[r]
	if ok {
		if isCellHighlighted {
			return true, renderer.Copy(fa.highlight, srcRect, dstRect)
		} else {
			return true, renderer.Copy(fa.normal, srcRect, dstRect)
		}
	}

	return false, nil
}

func _printCellToSurface(cell *types.Cell, cellRect *sdl.Rect, font *ttf.Font, surface *sdl.Surface, isCellHighlighted bool) error {
	fg, bg := sgrOpts(cell.Sgr)

	/*sr.setFontStyle(cell.Sgr.Bitwise)
	r := cell.Char
	if r == 0 {
		r = ' '
	}*/

	// render background colour

	if bg != nil {
		var pixel uint32
		if isCellHighlighted {
			pixel = sdl.MapRGBA(surface.Format, textHighlight.R, textHighlight.G, textHighlight.B, 255)
		} else {
			pixel = sdl.MapRGBA(surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		}
		err := surface.FillRect(cellRect, pixel)
		if err != nil {
			return err
		}
	}

	// render drop shadow

	if config.Config.Terminal.TypeFace.DropShadow && (bg == nil || isCellHighlighted) {
		shadowRect := &sdl.Rect{
			X: cellRect.X + dropShadowOffset,
			Y: cellRect.Y + dropShadowOffset,
			W: cellRect.W,
			H: cellRect.H,
		}

		var c sdl.Color
		if isCellHighlighted && bg == nil {
			c = textHighlight
		} else {
			c = textShadow
		}
		shadowText, err := font.RenderGlyphBlended(cell.Char, c)
		if err != nil {
			return err
		}
		defer shadowText.Free()

		err = shadowText.Blit(nil, surface, shadowRect)
		if err != nil {
			return err
		}
	}

	// render cell char
	text, err := font.RenderGlyphBlended(cell.Char, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()
	if isCellHighlighted {
		text.SetBlendMode(sdl.BLENDMODE_ADD)
	}

	err = text.Blit(nil, surface, cellRect)
	if err != nil {
		return err
	}

	return nil
}

func (sr *sdlRender) setFontStyle(style types.SgrFlag) {
	/*if style == sr._fontStyle {
		return
	}*/

	sr.font.SetStyle(fontStyle(style))
	//sr._fontStyle = style
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

func (sr *sdlRender) PrintCell(cell *types.Cell, cellPos *types.XY) error {
	if cell.Char == 0 {
		return nil
	}

	dstRect := &sdl.Rect{
		X: (sr.glyphSize.X * cellPos.X) + sr.border,
		Y: (sr.glyphSize.Y * cellPos.Y) + sr.border,
		W: sr.glyphSize.X + dropShadowOffset,
		H: sr.glyphSize.Y + dropShadowOffset,
	}

	isCellHighlighted := isCellHighlighted(sr, dstRect)
	hash := cell.Sgr.HashValue()

	ok, err := sr.fontCache.atlas.Render(sr.renderer, dstRect, cell.Char, hash, isCellHighlighted)
	if ok || err != nil {
		return err
	}

	extAtlases, ok := sr.fontCache.extended[hash]
	if ok {
		for i := range extAtlases {
			ok, err = extAtlases[i].Render(sr.renderer, dstRect, cell.Char, hash, isCellHighlighted)
			if ok || err != nil {
				return err
			}
		}
	}

	atlas := newFontAtlas([]rune{cell.Char}, cell.Sgr, sr.glyphSize, sr.renderer, sr.font)
	sr.fontCache.extended[hash] = append(sr.fontCache.extended[hash], atlas)
	_, err = atlas.Render(sr.renderer, dstRect, cell.Char, hash, isCellHighlighted)
	return err
}
