package rendersdl

import (
	"bytes"
	"fmt"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"golang.design/x/clipboard"
)

var (
	highlightBorder = &types.Colour{0x31, 0x6d, 0xb0}
	highlightFill   = &types.Colour{0x1c, 0x3e, 0x64}
)

const (
	_HIGHLIGHT_MODE_PNG = 0 + iota
	_HIGHLIGHT_MODE_SQUARE
	_HIGHLIGHT_MODE_LINES
)

type highlighterT struct {
	button uint8
	rect   *sdl.Rect
	mode   uint8
}

func (hl *highlighterT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	// do nothing
}

func (hl *highlighterT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	if evt.Keysym.Sym == sdl.K_ESCAPE {
		sr.highlighter = nil
		return
	}

	hl.modifier(evt.Keysym.Mod)
}

func (hl *highlighterT) modifier(mod uint16) {
	switch {
	case mod&sdl.KMOD_CTRL != 0:
		fallthrough
	case mod&sdl.KMOD_LCTRL != 0:
		fallthrough
	case mod&sdl.KMOD_RCTRL != 0:
		hl.mode = _HIGHLIGHT_MODE_SQUARE

	case mod&sdl.KMOD_SHIFT != 0:
		fallthrough
	case mod&sdl.KMOD_LSHIFT != 0:
		fallthrough
	case mod&sdl.KMOD_RSHIFT != 0:
		hl.mode = _HIGHLIGHT_MODE_LINES
	}
}

func (hl *highlighterT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	hl.button = 0

	normaliseRect(hl.rect)

	switch hl.mode {
	case _HIGHLIGHT_MODE_PNG:
		if hl.rect.W < sr.glyphSize.X && hl.rect.H < sr.glyphSize.Y {
			sr.highlighter = nil
			b := clipboard.Read(clipboard.FmtText)
			if len(b) != 0 {
				sr.term.Reply(b)
			} else {
				sr.DisplayNotification(types.NOTIFY_INFO, "Clipboard does not contain text to paste")
			}
		}
		// clipboard copy will happen automatically on next redraw
		sr.TriggerRedraw()

	case _HIGHLIGHT_MODE_LINES:
		rect := sr.rectPxToCells(hl.rect)
		lines := sr.term.CopyLines(rect.Y, rect.H)
		clipboard.Write(clipboard.FmtText, lines)
		sr.highlighter = nil
		count := bytes.Count(lines, []byte{'\n'}) + 1
		sr.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%d lines have been copied to clipboard", count))

	case _HIGHLIGHT_MODE_SQUARE:
		rect := sr.rectPxToCells(hl.rect)
		lines := sr.term.CopySquare(&types.XY{X: rect.X, Y: rect.Y}, &types.XY{X: rect.W, Y: rect.H})
		clipboard.Write(clipboard.FmtText, lines)
		sr.highlighter = nil
		sr.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%dx%d grid has been copied to clipboard", rect.W-rect.X+1, rect.H-rect.Y+1))
	}

}

func (hl *highlighterT) eventMouseWheel(sr *sdlRender, evt *sdl.MouseWheelEvent) {
	// do nothing
}

func (hl *highlighterT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	hl.rect.W += evt.XRel
	hl.rect.H += evt.YRel
	sr.TriggerRedraw()
}

func (sr *sdlRender) selectionHighlighter() {
	if sr.highlighter == nil {
		return
	}

	var alphaBorder, alphaFill uint8
	if sr.highlighter.mode == _HIGHLIGHT_MODE_PNG {
		alphaBorder, alphaFill = 190, 64
	} else {
		alphaBorder, alphaFill = 64, 32
	}

	sr.renderer.SetDrawColor(highlightBorder.Red, highlightBorder.Green, highlightBorder.Blue, alphaBorder)
	rect := sdl.Rect{
		X: sr.highlighter.rect.X - 1,
		Y: sr.highlighter.rect.Y - 1,
		W: sr.highlighter.rect.W + 2,
		H: sr.highlighter.rect.H + 2,
	}
	sr.renderer.DrawRect(&rect)
	rect = sdl.Rect{
		X: sr.highlighter.rect.X,
		Y: sr.highlighter.rect.Y,
		W: sr.highlighter.rect.W,
		H: sr.highlighter.rect.H,
	}
	sr.renderer.DrawRect(&rect)

	// fill background
	sr.renderer.SetDrawColor(highlightFill.Red, highlightFill.Green, highlightFill.Blue, alphaFill)
	rect = sdl.Rect{
		X: sr.highlighter.rect.X + 1,
		Y: sr.highlighter.rect.Y + 1,
		W: sr.highlighter.rect.W - 2,
		H: sr.highlighter.rect.H - 2,
	}
	sr.renderer.FillRect(&rect)
}

func isCellHighlighted(sr *sdlRender, rect *sdl.Rect) bool {
	if sr.highlighter == nil || sr.highlighter.button == 0 {
		return false
	}

	hlRect := *sr.highlighter.rect
	normaliseRect(&hlRect)

	runeCell := sr.rectPxToCells(rect)
	hlCell := sr.rectPxToCells(&hlRect)

	switch sr.highlighter.mode {
	default:
		return false
	case _HIGHLIGHT_MODE_LINES:
		return runeCell.Y >= hlCell.Y && runeCell.Y <= hlCell.H
	case _HIGHLIGHT_MODE_SQUARE:
		return runeCell.X >= hlCell.X && runeCell.X <= hlCell.W &&
			runeCell.Y >= hlCell.Y && runeCell.Y <= hlCell.H
	}
}

func normaliseRect(rect *sdl.Rect) {
	if rect.W < 0 {
		rect.X += rect.W
		rect.W *= -1
	}

	if rect.H < 0 {
		rect.Y += rect.H
		rect.H *= -1
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
