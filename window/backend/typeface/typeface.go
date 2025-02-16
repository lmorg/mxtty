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
	fontFile *ttf.Font
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

func Open(name string, size int) (err error) {
	if name != "" {
		fontFile, err = openSystemTtf(name, size)
	}
	if name == "" || err != nil {
		fontFile, err = openCompiledTtf(size)
	}

	if err != nil {
		return err
	}

	fontFile.SetHinting(ttf.HINTING_MONO)

	fontSize, err = getSize(fontFile)
	return err
}

func GetSize() *types.XY {
	return fontSize
}

func getSize(font *ttf.Font) (*types.XY, error) {
	x, y, err := font.SizeUTF8("W")
	return &types.XY{int32(x), int32(y)}, err
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
