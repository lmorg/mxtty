package typeface

import (
	"image"
	"image/color"
	"image/draw"
	"unsafe"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

// TODO: refine this function
func renderGlyphHarfbuzz(char rune, fg, bg *types.Colour, cellRect *sdl.Rect) (*sdl.Surface, error) {
	if false {
		//return fontFile.RenderGlyphBlended(char, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
		return fontFile.RenderGlyphBlended(char, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	}

	//fill := color.Color(bg)
	if bg == nil {
		bg = types.SGR_DEFAULT.Bg
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(cellRect.W), int(cellRect.H)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	//data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	typeface, err := font.ParseTTF(assets.Reader(assets.TYPEFACE))
	if err != nil {
		return nil, err
	}

	textRenderer := &render.Renderer{
		FontSize: float32(config.Config.Terminal.TypeFace.FontSize),
		Color:    fg,
	}

	_ = textRenderer.DrawString(string(char), img, typeface)

	pitch := cellRect.W * 4
	//r, g, b, a := bg.RGBA()
	//return sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&img.Pix[0]), cellRect.W, cellRect.H, 32, int(pitch), r, g, b, a)
	return sdl.CreateRGBSurfaceWithFormatFrom(unsafe.Pointer(&img.Pix[0]), cellRect.W, cellRect.H, 32, pitch, uint32(sdl.PIXELFORMAT_RGBA32))
}
