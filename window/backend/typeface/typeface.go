package typeface

import (
	"fmt"
	"log"

	"github.com/flopp/go-findfont"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	fontSize *types.XY
	adjustX  int32 = 0
	adjustY  int32 = 1
)

func init() {
	err := ttf.Init()
	if err != nil {
		panic(err.Error())
	}
}

func Close() {
	ttf.Quit()
}

func Open(name string, size int) (font *ttf.Font, err error) {
	if name != "" {
		font, err = openSystemTtf(name, size)
	}
	if name == "" || err != nil {
		font, err = openCompiledTtf(size)
	}

	if err != nil {
		return nil, err
	}

	font.SetHinting(ttf.HINTING_MONO)

	fontSize, err = getSize(font)
	return font, err
}

func GetSize() *types.XY {
	return fontSize
}

func getSize(font *ttf.Font) (*types.XY, error) {
	surface, err := font.RenderGlyphSolid('W', sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		return nil, fmt.Errorf("error in font.RenderGlyphBlended(): %s", err.Error())
	}
	return &types.XY{
		X: surface.W + adjustX,
		Y: surface.H + adjustY,
	}, nil
}

func openSystemTtf(name string, size int) (*ttf.Font, error) {
	path, err := findfont.Find(name)
	if err != nil {
		//return nil, fmt.Errorf("error in findfont.Find(): %s", err.Error())
		log.Printf("error in findfont.Find(): %s", err.Error())
		log.Println("defaulting to compiled log...")
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, fmt.Errorf("error in ttf.OpenFont(): %s", err.Error())
	}
	return font, nil
}

func openCompiledTtf(size int) (*ttf.Font, error) {
	rwops, err := sdl.RWFromMem(assets.Get(assets.TYPEFACE))
	if err != nil {
		return nil, fmt.Errorf("error in sdl.RWFromMem(): %s", err.Error())
	}

	font, err := ttf.OpenFontRW(rwops, 0, size)
	if err != nil {
		return nil, fmt.Errorf("error in ttf.OpenFontRW(): %s", err.Error())
	}
	return font, nil
}
