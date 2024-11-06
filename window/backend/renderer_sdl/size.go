package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) convertPxToCellXY(x, y int32) *types.XY {
	xy := &types.XY{
		X: (x - sr.border) / sr.glyphSize.X,
		Y: (y - sr.border) / sr.glyphSize.Y,
	}

	if xy.X < 0 {
		xy.X = 0
	} else if xy.X >= sr.term.GetSize().X {
		xy.X = sr.term.GetSize().X - 1
	}
	if xy.Y < 0 {
		xy.Y = 0
	} else if xy.Y >= sr.term.GetSize().Y {
		xy.Y = sr.term.GetSize().Y - 1
	}

	return xy
}

func normaliseRect(rect *sdl.Rect) {
	if rect.W < 0 {
		rect.X += rect.W
		rect.W = -rect.W
	}

	if rect.H < 0 {
		rect.Y += rect.H
		rect.H = -rect.H
	}
}

func (sr *sdlRender) rectPxToCells(rect *sdl.Rect) *sdl.Rect {
	return &sdl.Rect{
		X: (rect.X - sr.border) / sr.glyphSize.X,
		Y: (rect.Y - sr.border) / sr.glyphSize.Y,
		W: ((rect.X + rect.W - sr.border) / sr.glyphSize.X),
		H: ((rect.Y + rect.H - sr.border) / sr.glyphSize.Y),
	}
}
