package rendersdl

import (
	"strings"
	"unsafe"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/renderer_sdl/layer"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const dropShadowOffset int32 = 1

const (
	_HLTEXTURE_NONE      = iota
	_HLTEXTURE_SELECTION // should always be first non-zero value
	_HLTEXTURE_SEARCH_RESULT
	_HLTEXTURE_MATCH_RANGE
	_HLTEXTURE_LAST // placeholder for rect calculations. Must always come last
)

var textShadow = []sdl.Color{
	_HLTEXTURE_NONE:          {R: 0, G: 0, B: 0, A: 0}, // A controlled by LightMode
	_HLTEXTURE_SELECTION:     {R: 64, G: 64, B: 255, A: 255},
	_HLTEXTURE_SEARCH_RESULT: {R: 64, G: 64, B: 255, A: 192},
	_HLTEXTURE_MATCH_RANGE:   {R: 64, G: 255, B: 64, A: 128},
}

var (
	fontAtlasCharacterPreload = []string{
		" ",                          // whitespace
		"1234567890",                 // numeric
		"abcdefghijklmnopqrstuvwxyz", // alpha, lower
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ", // alpha, upper
		"`",                          // backtick
		`!"£$%^&*()-=_+`,             // special, top row
		`[]{};'#:@~\|,./<>?`,         // special, others ascii
		`↑↓←→`,                       // special, mxtty
		`»…`,                         // special, murex
		`┏┓┗┛━─╶┃┠┨╔╗╚╝═║╟╢█`, //     // box drawing
	}
)

type fontCacheDefaultLookupT map[rune]*sdl.Rect
type fontAtlasT struct {
	sgrHash uint64
	lookup  fontCacheDefaultLookupT
	texture []*sdl.Texture
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
		sgrHash: sgr.HashValue(),
		lookup:  _newFontCacheDefaultLookup(chars, glyphSizePlusShadow),
		texture: []*sdl.Texture{
			_HLTEXTURE_NONE:          _newFontTexture(chars, sgr, glyphSizePlusShadow, renderer, font, _HLTEXTURE_NONE),
			_HLTEXTURE_SELECTION:     _newFontTexture(chars, sgr, glyphSizePlusShadow, renderer, font, _HLTEXTURE_SELECTION),
			_HLTEXTURE_SEARCH_RESULT: _newFontTexture(chars, sgr, glyphSizePlusShadow, renderer, font, _HLTEXTURE_SEARCH_RESULT),
		},
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

func _newFontTexture(chars []rune, sgr *types.Sgr, glyphSize *types.XY, renderer *sdl.Renderer, font *ttf.Font, hlTexture int) *sdl.Texture {
	surface := _newFontSurface(glyphSize, int32(len(chars)))
	defer surface.Free()

	cell := &types.Cell{
		Sgr: sgr,
	}

	cellRect := &sdl.Rect{W: glyphSize.X, H: glyphSize.Y}

	font.SetStyle(fontStyle(cell.Sgr.Bitwise))

	var i int
	var err error
	for i, cell.Char = range chars {
		cellRect.X = int32(i) * glyphSize.X
		err = _printCellToSurface(cell, cellRect, font, surface, hlTexture)
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

func _newFontSurface(glyphSize *types.XY, nCharacters int32) *sdl.Surface {
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, glyphSize.X*nCharacters, glyphSize.Y*_HLTEXTURE_LAST, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	pixel := sdl.MapRGBA(surface.Format, types.SGR_DEFAULT.Bg.Red, types.SGR_DEFAULT.Bg.Green, types.SGR_DEFAULT.Bg.Blue, 255)
	err = surface.FillRect(&sdl.Rect{W: surface.W, H: surface.H}, pixel)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	err = surface.SetColorKey(true, pixel)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}

	return surface
}

func _printCellToSurface(cell *types.Cell, cellRect *sdl.Rect, font *ttf.Font, surface *sdl.Surface, hlTexture int) error {
	fg, bg := sgrOpts(cell.Sgr, hlTexture == _HLTEXTURE_SELECTION)

	// render background colour

	if bg != nil {
		var pixel uint32
		if hlTexture != 0 {
			pixel = sdl.MapRGBA(surface.Format, textShadow[hlTexture].R, textShadow[hlTexture].G, textShadow[hlTexture].B, 255)
		} else {
			pixel = sdl.MapRGBA(surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		}
		fillRect := &sdl.Rect{
			X: cellRect.X,
			Y: cellRect.Y,
			W: cellRect.W - dropShadowOffset + 10, //,
			H: cellRect.H - dropShadowOffset + 10, //-1,
		}
		err := surface.FillRect(fillRect, pixel)
		if err != nil {
			return err
		}
	}

	// render drop shadow

	if (config.Config.TypeFace.DropShadow && bg == nil) ||
		hlTexture > _HLTEXTURE_SELECTION {

		shadowText, err := font.RenderGlyphBlended(cell.Char, textShadow[hlTexture])
		if err != nil {
			return err
		}
		defer shadowText.Free()

		shadowRect := &sdl.Rect{
			X: cellRect.X + dropShadowOffset,
			Y: cellRect.Y + dropShadowOffset,
			W: cellRect.W,
			H: cellRect.H,
		}
		_ = shadowText.Blit(nil, surface, shadowRect)

		if hlTexture > _HLTEXTURE_SELECTION {
			shadowRect = &sdl.Rect{
				X: cellRect.X - dropShadowOffset,
				Y: cellRect.Y + dropShadowOffset,
				W: cellRect.W,
				H: cellRect.H,
			}
			_ = shadowText.Blit(nil, surface, shadowRect)
			shadowRect = &sdl.Rect{
				X: cellRect.X - dropShadowOffset,
				Y: cellRect.Y - dropShadowOffset,
				W: cellRect.W,
				H: cellRect.H,
			}
			_ = shadowText.Blit(nil, surface, shadowRect)
			shadowRect = &sdl.Rect{
				X: cellRect.X + dropShadowOffset,
				Y: cellRect.Y - dropShadowOffset,
				W: cellRect.W,
				H: cellRect.H,
			}
			_ = shadowText.Blit(nil, surface, shadowRect)
		}
	}

	// render cell char
	text, err := typeface.RenderGlyph(cell.Char, fg, bg, cellRect)
	if err != nil {
		return err
	}
	defer text.Free()

	if hlTexture != 0 {
		_ = text.SetBlendMode(sdl.BLENDMODE_ADD)
	}

	err = text.Blit(nil, surface, cellRect)
	if err != nil {
		return err
	}
	if config.Config.TypeFace.Ligatures && cell.Sgr.Bitwise.Is(types.SGR_BOLD) {
		_ = text.SetBlendMode(sdl.BLENDMODE_ADD)
		_ = text.Blit(nil, surface, cellRect)
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

func fontStyle(style types.SgrFlag) ttf.Style {
	var i ttf.Style

	if style.Is(types.SGR_BOLD) && !config.Config.TypeFace.Ligatures {
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

func sgrOpts(sgr *types.Sgr, forceBg bool) (fg *types.Colour, bg *types.Colour) {
	if sgr.Bitwise.Is(types.SGR_INVERT) {
		bg, fg = sgr.Fg, sgr.Bg
	} else {
		fg, bg = sgr.Fg, sgr.Bg
	}

	if unsafe.Pointer(bg) == unsafe.Pointer(types.SGR_DEFAULT.Bg) && !forceBg {
		bg = nil
	}

	return fg, bg
}

func (sr *sdlRender) PrintCell(cell *types.Cell, cellPos *types.XY) {
	if cell.Char == 0 {
		return
	}

	dstRect := &sdl.Rect{
		X: (sr.glyphSize.X * cellPos.X) + _PANE_LEFT_MARGIN,
		Y: (sr.glyphSize.Y * cellPos.Y) + _PANE_TOP_MARGIN,
		W: sr.glyphSize.X + dropShadowOffset,
		H: sr.glyphSize.Y + dropShadowOffset,
	}

	hlTexture := _HLTEXTURE_NONE
	if cell.Sgr.Bitwise.Is(types.SGR_HIGHLIGHT_SEARCH_RESULT) {
		hlTexture = _HLTEXTURE_SEARCH_RESULT
	}
	if isCellHighlighted(sr, dstRect) {
		hlTexture = _HLTEXTURE_SELECTION
	}
	hash := cell.Sgr.HashValue()

	ok := sr.fontCache.atlas.Render(sr, dstRect, cell.Char, hash, hlTexture)
	if ok {
		return
	}

	extAtlases, ok := sr.fontCache.extended[hash]
	if ok {
		for i := range extAtlases {
			ok = extAtlases[i].Render(sr, dstRect, cell.Char, hash, hlTexture)
			if ok {
				return
			}
		}
	}

	atlas := newFontAtlas([]rune{cell.Char}, cell.Sgr, sr.glyphSize, sr.renderer, sr.font)
	sr.fontCache.extended[hash] = append(sr.fontCache.extended[hash], atlas)
	atlas.Render(sr, dstRect, cell.Char, hash, hlTexture)
}

func (fa *fontAtlasT) Render(sr *sdlRender, dstRect *sdl.Rect, r rune, hash uint64, hlMode int) bool {
	if hash != fa.sgrHash {
		return false
	}

	srcRect, ok := fa.lookup[r]
	if !ok {
		return false
	}

	texture := fa.texture[hlMode]

	sr.AddToElementStack(&layer.RenderStackT{texture, srcRect, dstRect, false})

	return ok
}
