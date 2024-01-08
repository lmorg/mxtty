package typeface

import (
	"github.com/flopp/go-findfont"
	"github.com/lmorg/mxtty/virtualterm/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var fontSize *types.Rect

func init() {
	err := ttf.Init()
	if err != nil {
		panic(err.Error())
	}
}

func Close() {
	ttf.Quit()
}

func Open(name string, size int) (*ttf.Font, error) {
	path, err := findfont.Find(name)
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}

	fontSize, err = getSize(font)
	return font, err
}

func GetSize() *types.Rect {
	return fontSize
}

func getSize(font *ttf.Font) (*types.Rect, error) {
	surface, err := font.RenderGlyphSolid('W', sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		return nil, err
	}
	return &types.Rect{
		X: surface.W,
		Y: surface.H,
	}, nil
}
