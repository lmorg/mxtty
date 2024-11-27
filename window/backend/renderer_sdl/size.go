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

// GetTermSize only exists so that elements can get the terminal size without
// having access to the term interface.
func (sr *sdlRender) GetTermSize() *types.XY {
	return sr.term.GetSize()
}

// _getSizeCells should only be called upon terminal resizing.
// All other checks for terminal size should come from term.GetSize()
func (sr *sdlRender) _getSizeCells() *types.XY {
	x, y, err := sr.renderer.GetOutputSize()
	if err != nil {
		panic("arg!")
	}
	//x, y := sr.window.GetSize()

	return &types.XY{
		X: ((x - (sr.border * 2)) / sr.glyphSize.X),
		Y: ((y - (sr.border * 2)) / sr.glyphSize.Y) - sr.footer, // inc footer
	}
}

///// resize

func (sr *sdlRender) windowResized() {
	sr.windowTabs = nil
	if sr.term != nil {
		sr.term.Resize(sr._getSizeCells())
	}
}

func (sr *sdlRender) ResizeWindow(size *types.XY) {
	go func() { sr._resize <- size }()
}

func (sr *sdlRender) _resizeWindow(size *types.XY) {
	w := (size.X * sr.glyphSize.X) + (sr.border * 2)
	h := ((size.Y + sr.footer) * sr.glyphSize.Y) + (sr.border * 2)
	sr.window.SetSize(w, h)
	sr.RefreshWindowList()
}
