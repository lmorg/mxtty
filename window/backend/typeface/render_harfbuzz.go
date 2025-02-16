package typeface

import (
	"fmt"
	"image"
	"os"
	"unsafe"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type fontHarfbuzz struct {
	//size *types.XY
	fonts [3]*font.Face
	sdl   *fontSdl // just needed while I figure out how to do everything in harfbuzz
}

func (f *fontHarfbuzz) Init() error {
	f.sdl = new(fontSdl)
	return f.sdl.Init()
}

func (f *fontHarfbuzz) Open(name string, size int) (err error) {
	var file font.Resource

	if name != "" {
		file, err = os.Open(name)
	}
	if name == "" || err != nil {
		file = assets.Reader(assets.TYPEFACE_DEFAULT)
	}

	f.fonts[_FONT_DEFAULT], err = font.ParseTTF(file)
	if err != nil {
		return err
	}

	f.fonts[_FONT_FALLBACK], err = font.ParseTTF(assets.Reader(assets.TYPEFACE_FALLBACK))
	if err != nil {
		panic(err)
	}
	f.fonts[_FONT_EMOJI], err = font.ParseTTF(assets.Reader(assets.TYPEFACE_EMOJI))
	if err != nil {
		panic(err)
	}

	return f.sdl.Open(name, size)
}

func (f *fontHarfbuzz) GetSize() *types.XY {
	return f.sdl.GetSize()
}

// RenderGlyph should be called from a font atlas
func (f *fontHarfbuzz) RenderGlyph(char rune, fg *types.Colour, cellRect *sdl.Rect) (*sdl.Surface, error) {
	img := image.NewNRGBA(image.Rect(0, 0, int(cellRect.W), int(cellRect.H)))

	textRenderer := &render.Renderer{
		FontSize: float32(config.Config.TypeFace.FontSize),
		Color:    fg,
	}

	for fontId := range f.fonts {
		if f.glyphIsProvided(fontId, char) {
			_ = textRenderer.DrawString(string(char), img, f.fonts[fontId])
			goto found
		}
		debug.Log(fmt.Sprintf("[harfbuzz] emoji not in font %d: %s", fontId, string(char)))
	}

	// not found
	_ = textRenderer.DrawString(string(char), img, f.fonts[_FONT_DEFAULT])

found:

	return sdl.CreateRGBSurfaceWithFormatFrom(
		unsafe.Pointer(&img.Pix[0]),
		cellRect.W, cellRect.H,
		32, cellRect.W*4, uint32(sdl.PIXELFORMAT_RGBA32),
	)
}

func (f *fontHarfbuzz) glyphIsProvided(fontId int, r rune) bool {
	_, found := f.fonts[fontId].NominalGlyph(r)
	return found
}

func (f *fontHarfbuzz) Close() {
	f.sdl.Close()
}

func (f *fontHarfbuzz) Deprecated_GetFont() *ttf.Font {
	return f.sdl.Deprecated_GetFont()
}
