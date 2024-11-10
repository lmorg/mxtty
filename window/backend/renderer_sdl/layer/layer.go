package layer

import "github.com/veandco/go-sdl2/sdl"

type RenderStackT struct {
	Texture *sdl.Texture
	SrcRect *sdl.Rect
	DstRect *sdl.Rect
	Destroy bool
}
