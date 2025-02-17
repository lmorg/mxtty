package typeface

import (
	"image"
	"log"
	"os"
	"unsafe"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type fontHarfbuzz struct {
	//size *types.XY
	face *font.Face
	fmap *fontscan.FontMap
	sdl  *fontSdl // just needed while I figure out how to do everything in harfbuzz
}

func (f *fontHarfbuzz) Init() error {
	f.fmap = fontscan.NewFontMap(log.Default())
	err := f.fmap.UseSystemFonts("")
	if err != nil {
		return err
	}
	f.fmap.SetQuery(fontscan.Query{Families: []string{"monospace", "emoji", "math", "fantasy"}})

	f.sdl = new(fontSdl)
	return f.sdl.Init()
}

func (f *fontHarfbuzz) Open(name string, size int) (err error) {
	var file font.Resource

	if name != "" {
		file, err = os.Open(name)
	}
	if name == "" || err != nil {
		file = assets.Reader(assets.TYPEFACE)
	}

	f.face, err = font.ParseTTF(file)

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

	if f.glyphIsProvided(0, char) {
		_ = textRenderer.DrawString(string(char), img, f.face)
		goto found
	}

	// not found
	_ = textRenderer.DrawString(string(char), img, f.fmap.ResolveFace(char))

found:

	return sdl.CreateRGBSurfaceWithFormatFrom(
		unsafe.Pointer(&img.Pix[0]),
		cellRect.W, cellRect.H,
		32, cellRect.W*4, uint32(sdl.PIXELFORMAT_RGBA32),
	)
}

func (f *fontHarfbuzz) glyphIsProvided(_ int, r rune) bool {
	_, found := f.face.NominalGlyph(r)
	return found
}

func (f *fontHarfbuzz) Close() {
	f.sdl.Close()
}

func (f *fontHarfbuzz) Deprecated_GetFont() *ttf.Font {
	return f.sdl.Deprecated_GetFont()
}
