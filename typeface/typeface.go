package typeface

import (
	"github.com/flopp/go-findfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SizeT struct {
	Width  int32
	Height int32
}

var fontSize *SizeT

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

func GetSize() *SizeT {
	return fontSize
}

func getSize(font *ttf.Font) (*SizeT, error) {
	surface, err := font.RenderGlyphSolid('W', sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		return nil, err
	}
	return &SizeT{
		Height: surface.H,
		Width:  surface.W,
	}, nil
}
