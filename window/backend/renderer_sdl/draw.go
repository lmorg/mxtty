package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/renderer_sdl/layer"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) DrawTable(pos *types.XY, height int32, boundaries []int32) {
	var err error

	fg := types.SGR_DEFAULT.Fg

	texture := sr.createRendererTexture()
	if texture == nil {
		return
	}
	defer sr.restoreRendererTexture()

	sr.renderer.SetDrawColor(fg.Red, fg.Green, fg.Blue, 128)

	X := (pos.X * sr.glyphSize.X) + _PANE_LEFT_MARGIN
	Y := (pos.Y * sr.glyphSize.Y) + _PANE_TOP_MARGIN
	H := Y + ((height + 1) * sr.glyphSize.Y)

	err = sr.renderer.DrawLine(X, Y, X, H)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	for i := range boundaries {
		x := X + (boundaries[i] * sr.glyphSize.X)
		err = sr.renderer.DrawLine(x, Y, x, H)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}
	}

	x := X + (boundaries[len(boundaries)-1] * sr.glyphSize.X)
	y := Y + ((height + 1) * sr.glyphSize.Y)
	err = sr.renderer.DrawLine(X, y, x, y)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	sr.renderer.SetDrawColor(fg.Red, fg.Green, fg.Blue, 100)

	for i := int32(0); i <= height; i++ {
		y = Y + (i * sr.glyphSize.Y)
		err = sr.renderer.DrawLine(X, y, x, y)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}
	}
}

const (
	highlightAlphaBorder = 190
	highlightAlphaFill   = 128
)

var highlightBlendMode sdl.BlendMode // controlled by LightMode

func (sr *sdlRender) DrawHighlightRect(topLeftCell, bottomRightCell *types.XY) {
	sr._drawHighlightRect(
		&sdl.Rect{
			X: (topLeftCell.X * sr.glyphSize.X) + _PANE_LEFT_MARGIN,
			Y: (topLeftCell.Y * sr.glyphSize.Y) + _PANE_TOP_MARGIN,
			W: (bottomRightCell.X * sr.glyphSize.X),
			H: (bottomRightCell.Y * sr.glyphSize.Y),
		},
		highlightAlphaBorder, highlightAlphaFill)
}

func (sr *sdlRender) _drawHighlightRect(rect *sdl.Rect, alphaBorder, alphaFill byte) {
	texture := sr.createRendererTexture()
	if texture == nil {
		return
	}
	defer sr.renderer.SetRenderTarget(nil)

	err := texture.SetBlendMode(highlightBlendMode)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	sr.renderer.SetDrawColor(highlightBorder.Red, highlightBorder.Green, highlightBorder.Blue, alphaBorder)
	rect.X -= 1
	rect.Y -= 1
	rect.W += 2
	rect.H += 2

	sr.renderer.DrawRect(rect)
	rect.X += 1
	rect.Y += 1
	rect.W -= 2
	rect.H -= 2
	sr.renderer.DrawRect(rect)

	// fill background

	sr.renderer.SetDrawColor(highlightFill.Red, highlightFill.Green, highlightFill.Blue, alphaFill)
	rect.X += 1
	rect.Y += 1
	rect.W -= 2
	rect.H -= 2
	sr.renderer.FillRect(rect)

	sr.AddToOverlayStack(&layer.RenderStackT{texture, nil, nil, true})
}
