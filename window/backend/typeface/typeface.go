package typeface

import (
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type typefaceRenderer interface {
	Init() error
	Open(string, int) error
	GetSize() *types.XY
	RenderGlyph(rune, *types.Colour, *types.Colour, *sdl.Rect) (*sdl.Surface, error)
	Deprecated_GetFont() *ttf.Font
	Close()
}

var renderer typefaceRenderer

func Init() error {
	if config.Config.TypeFace.HarfbuzzRenderer {
		renderer = new(fontHarfbuzz)
	} else {
		renderer = new(fontSdl)
	}

	return renderer.Init()
}

func Open(name string, size int) (err error) {
	return renderer.Open(name, size)
}

func GetSize() *types.XY {
	return renderer.GetSize()
}

func RenderGlyph(char rune, fg, bg *types.Colour, cellRect *sdl.Rect) (*sdl.Surface, error) {
	return renderer.RenderGlyph(char, fg, bg, cellRect)
}

func Close() {
	ttf.Quit()
}

func Deprecated_GetFont() *ttf.Font {
	return renderer.Deprecated_GetFont()
}
