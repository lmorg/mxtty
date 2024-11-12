package rendersdl

import (
	"fmt"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

const footer = 1

func (sr *sdlRender) StatusbarText(s string) {
	sr.footerText = s
	sr.TriggerRedraw()
}

func (sr *sdlRender) renderFooter() {
	_ = sr.createRendererTexture()

	rect := &sdl.Rect{
		X: 0,
		Y: (sr.term.GetSize().Y * sr.glyphSize.Y) + sr.border,
		W: (sr.term.GetSize().X * sr.glyphSize.X) + (sr.border * 3),
		H: sr.glyphSize.Y + (sr.border * 2),
	}
	fill := types.SGR_COLOUR_BLACK_BRIGHT
	_ = sr.renderer.SetDrawColor(fill.Red, fill.Green, fill.Blue, 255)
	_ = sr.renderer.FillRect(rect)

	if sr.footerText == "" {
		sr.footerText = fmt.Sprintf("%s (version %s)  |  [F5] Show / hide window", app.Title, app.Version())
	}

	pos := &types.XY{Y: sr.term.GetSize().Y}
	cells := sr._makeFooter()

	sr.restoreRendererTexture()
	sr.PrintCellBlock(cells, pos)
}

func (sr *sdlRender) _makeFooter() []types.Cell {
	footer := make([]types.Cell, sr.term.GetSize().X)

	var i int
	text := []rune(sr.footerText)
	for ; i < len(text) && i < len(footer); i++ {
		footer[i].Char = text[i]
		footer[i].Sgr = types.SGR_DEFAULT.Copy()
	}

	return footer[:i]
}
