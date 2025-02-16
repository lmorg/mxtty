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

type fontSdl struct {
	size *types.XY
	font *ttf.Font
}

func (f *fontSdl) Init() error { return ttf.Init() }

func (f *fontSdl) Open(name string, size int) (err error) {
	if name != "" {
		err = f.openSystemTtf(name, size)
	}
	if name == "" || err != nil {
		err = f.openCompiledTtf(size)
	}

	if err != nil {
		return err
	}

	f.font.SetHinting(ttf.HINTING_MONO)

	err = f._getSize()
	return err
}

func (f *fontSdl) _getSize() error {
	x, y, err := f.font.SizeUTF8("W")
	f.size = &types.XY{int32(x), int32(y)}
	return err
}

func (f *fontSdl) GetSize() *types.XY {
	return f.size
}

func (f *fontSdl) openSystemTtf(name string, size int) error {
	path, err := findfont.Find(name)
	if err != nil {
		log.Printf("error in findfont.Find(): %s", err.Error())
		log.Println("defaulting to compiled log...")
	}

	f.font, err = ttf.OpenFont(path, size)
	if err != nil {
		return fmt.Errorf("error in ttf.OpenFont(): %s", err.Error())
	}

	return nil
}

func (f *fontSdl) openCompiledTtf(size int) error {
	rwops, err := sdl.RWFromMem(assets.Get(assets.TYPEFACE))
	if err != nil {
		return fmt.Errorf("error in sdl.RWFromMem(): %s", err.Error())
	}

	f.font, err = ttf.OpenFontRW(rwops, 0, size)
	if err != nil {
		return fmt.Errorf("error in ttf.OpenFontRW(): %s", err.Error())
	}
	return nil
}

// RenderGlyph should be called from a font atlas
func (f *fontSdl) RenderGlyph(char rune, fg, bg *types.Colour, cellRect *sdl.Rect) (*sdl.Surface, error) {
	return f.font.RenderGlyphBlended(char, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
}

func (f *fontSdl) Close() {
	ttf.Quit()
}

func (f *fontSdl) Deprecated_GetFont() *ttf.Font {
	return f.font
}
