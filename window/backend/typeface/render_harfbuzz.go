package typeface

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"unsafe"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type fontHarfbuzz struct {
	//size *types.XY
	font *font.Face
	sdl  *fontSdl // just needed while I figure out how to do everything in harfbuzz
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
		file = assets.Reader(assets.TYPEFACE)
	}

	f.font, err = font.ParseTTF(file)
	if err != nil {
		return err
	}

	return f.sdl.Open(name, size)
}

func (f *fontHarfbuzz) GetSize() *types.XY {
	return f.sdl.GetSize()
}

// RenderGlyph should be called from a font atlas
func (f *fontHarfbuzz) RenderGlyph(char rune, fg, bg *types.Colour, cellRect *sdl.Rect) (*sdl.Surface, error) {
	if bg == nil {
		bg = types.SGR_DEFAULT.Bg
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(cellRect.W), int(cellRect.H)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	textRenderer := &render.Renderer{
		FontSize: float32(config.Config.TypeFace.FontSize),
		Color:    fg,
	}

	_ = textRenderer.DrawString(string(char), img, f.font)

	return sdl.CreateRGBSurfaceWithFormatFrom(
		unsafe.Pointer(&img.Pix[0]),
		cellRect.W, cellRect.H,
		32, cellRect.W*4, uint32(sdl.PIXELFORMAT_RGBA32),
	)
}

func (f *fontHarfbuzz) Close() {
	f.sdl.Close()
}

func (f *fontHarfbuzz) Deprecated_GetFont() *ttf.Font {
	return f.sdl.Deprecated_GetFont()
}
