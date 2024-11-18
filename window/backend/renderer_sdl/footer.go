package rendersdl

import (
	"fmt"

	"github.com/lmorg/mxtty/app"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) StatusBarText(s string) {
	sr.footerText = s
	sr.TriggerRedraw()
}

func (sr *sdlRender) renderFooter() {
	if sr.footer == 0 {
		return
	}

	_ = sr.createRendererTexture()

	rect := &sdl.Rect{
		X: 0,
		Y: (sr.term.GetSize().Y * sr.glyphSize.Y) + sr.border,
		W: (sr.term.GetSize().X * sr.glyphSize.X) + (sr.border * 3),
		H: (sr.glyphSize.Y * sr.footer) + (sr.border * 2),
	}

	fill := types.SGR_COLOUR_BLACK_BRIGHT
	_ = sr.renderer.SetDrawColor(fill.Red, fill.Green, fill.Blue, 255)
	_ = sr.renderer.FillRect(rect)

	sr.restoreRendererTexture()

	pos := &types.XY{Y: sr.term.GetSize().Y}

	if !config.Config.Window.StatusBar {
		goto tmuxIntegration
	}

	if sr.footerText == "" {
		sr.footerText = fmt.Sprintf("%s (version %s)  |  [F5] Show / hide window", app.Title, app.Version())
	}

	sr._footerRenderStatusBar(pos)
	pos.Y++

tmuxIntegration:
	if sr.tmux == nil {
		// This shouldn't happen, but saves a crash in case of this getting
		// invoked before tmux has finished getting set up
		return
	}

	_ = sr.createRendererTexture()
	rect.Y += sr.glyphSize.Y
	_ = sr.renderer.SetDrawColor(fill.Red, fill.Green, fill.Blue, 255)
	_ = sr.renderer.FillRect(rect)
	sr.restoreRendererTexture()

	if sr.windowTabs == nil {
		sr._footerCacheTmuxWindowTabs(pos)
	}

	sr._footerRenderTmuxWindowTabs(pos)
}

func (sr *sdlRender) _footerRenderStatusBar(pos *types.XY) {
	footer := make([]types.Cell, sr.term.GetSize().X)

	var i int
	text := []rune(sr.footerText)
	for ; i < len(text) && i < len(footer); i++ {
		footer[i].Char = text[i]
		footer[i].Sgr = types.SGR_DEFAULT.Copy()
	}

	sr.PrintCellBlock(footer[:i], pos)
}

func (sr *sdlRender) _footerCacheTmuxWindowTabs(pos *types.XY) {
	tabList := &tabListT{
		mouseOver: -1,
	}

	heading := []rune("Window tab list →  ")

	cell := types.Cell{
		Char: ' ',
		Sgr:  types.SGR_DEFAULT.Copy(),
	}

	tabList.cells = append(tabList.cells, cell)
	for _, r := range heading {

		cell.Char = r
		tabList.cells = append(tabList.cells, cell)
	}

	tabList.boundaries = []int32{0}
	var x int32

	tabList.windows = sr.tmux.RenderWindows()
	for i, win := range tabList.windows {
		if win.Active {
			tabList.active = i
		}
		for _, r := range win.Name {
			cell.Char = r
			tabList.cells = append(tabList.cells, cell)
			x++
		}

		cell.Char = ' '
		tabList.cells = append(tabList.cells, cell, cell)
		x += 2
		tabList.boundaries = append(tabList.boundaries, x)
	}

	tabList.offset = &types.XY{X: int32(len(heading)), Y: pos.Y}

	sr.windowTabs = tabList
}

func (sr *sdlRender) _footerRenderTmuxWindowTabs(pos *types.XY) {
	sr.PrintCellBlock(sr.windowTabs.cells, pos)
	sr.DrawTable(sr.windowTabs.offset, 0, sr.windowTabs.boundaries[1:])

	var (
		topLeftCellX     = sr.windowTabs.offset.X + sr.windowTabs.boundaries[sr.windowTabs.active]
		topLeftCellY     = sr.windowTabs.offset.Y
		bottomRightCellX = sr.windowTabs.boundaries[sr.windowTabs.active+1] - sr.windowTabs.boundaries[sr.windowTabs.active]
		bottomRightCellY = int32(1)
	)

	activeRect := &sdl.Rect{
		X: (topLeftCellX * sr.glyphSize.X) + sr.border,
		Y: (topLeftCellY * sr.glyphSize.Y) + sr.border,
		W: (bottomRightCellX * sr.glyphSize.X) + 1,
		H: (bottomRightCellY * sr.glyphSize.Y) + 1,
	}
	sr._drawHighlightRect(activeRect, 0, 230)

	if sr.windowTabs.mouseOver == -1 {
		return
	}

	topLeftCellX = sr.windowTabs.offset.X + sr.windowTabs.boundaries[sr.windowTabs.mouseOver]
	bottomRightCellX = sr.windowTabs.boundaries[sr.windowTabs.mouseOver+1] - sr.windowTabs.boundaries[sr.windowTabs.mouseOver]

	highlightRect := &sdl.Rect{
		X: (topLeftCellX * sr.glyphSize.X) + sr.border,
		Y: (topLeftCellY * sr.glyphSize.Y) + sr.border,
		W: (bottomRightCellX * sr.glyphSize.X),
		H: (bottomRightCellY * sr.glyphSize.Y),
	}
	sr._drawHighlightRect(highlightRect, highlightAlphaBorder, highlightAlphaFill)
}
