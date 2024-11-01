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

type _highlightMode uint8

const (
	_HIGHLIGHT_MODE_PNG _highlightMode = 0 + iota
	_HIGHLIGHT_MODE_SQUARE
	_HIGHLIGHT_MODE_FULL_LINES
	_HIGHLIGHT_MODE_LINE_RANGE
)

type highlighterT struct {
	button uint8
	rect   *sdl.Rect
	mode   _highlightMode
}

func (hl *highlighterT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	// do nothing
}

func (hl *highlighterT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	if evt.Keysym.Sym == sdl.K_ESCAPE {
		sr.highlighter = nil
		sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
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
		hl.setMode(_HIGHLIGHT_MODE_SQUARE)

	case mod&sdl.KMOD_SHIFT != 0:
		fallthrough
	case mod&sdl.KMOD_LSHIFT != 0:
		fallthrough
	case mod&sdl.KMOD_RSHIFT != 0:
		hl.setMode(_HIGHLIGHT_MODE_LINE_RANGE)

	case mod&sdl.KMOD_ALT != 0:
		fallthrough
	case mod&sdl.KMOD_LALT != 0:
		fallthrough
	case mod&sdl.KMOD_RALT != 0:
		hl.setMode(_HIGHLIGHT_MODE_FULL_LINES)
	}
}

func (hl *highlighterT) setMode(mode _highlightMode) {
	hl.mode = mode
	switch mode {
	case _HIGHLIGHT_MODE_LINE_RANGE:
		sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM))
	default:
		sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
	}
}

func (hl *highlighterT) eventMouseButton(sr *sdlRender, _ *sdl.MouseButtonEvent) {
	hl.button = 0
	sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))

	switch hl.mode {
	case _HIGHLIGHT_MODE_PNG:
		if hl.rect.W < sr.glyphSize.X && hl.rect.H < sr.glyphSize.Y {
			sr.clipboardPasteText()
		}
		// clipboard copy will happen automatically on next redraw
		sr.TriggerRedraw()

	case _HIGHLIGHT_MODE_FULL_LINES:
		normaliseRect(hl.rect)
		rect := sr.rectPxToCells(hl.rect)
		lines := sr.term.CopyLines(rect.Y, rect.H)
		clipboard.Write(clipboard.FmtText, lines)
		sr.highlighter = nil
		count := bytes.Count(lines, []byte{'\n'}) + 1
		sr.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%d lines have been copied to clipboard", count))

	case _HIGHLIGHT_MODE_SQUARE:
		normaliseRect(hl.rect)
		rect := sr.rectPxToCells(hl.rect)
		lines := sr.term.CopySquare(&types.XY{X: rect.X, Y: rect.Y}, &types.XY{X: rect.W, Y: rect.H})
		clipboard.Write(clipboard.FmtText, lines)
		sr.highlighter = nil
		sr.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%dx%d grid has been copied to clipboard", rect.W-rect.X+1, rect.H-rect.Y+1))

	case _HIGHLIGHT_MODE_LINE_RANGE:
		rect := sr.rectPxToCells(hl.rect)
		if rect.X-rect.W < 2 && rect.X-rect.W > -2 && rect.Y-rect.H < 2 && rect.Y-rect.H > -2 {
			sr.highlighter = nil
			return
		}
		lines := sr.term.CopyRange(&types.XY{X: rect.X, Y: rect.Y}, &types.XY{X: rect.W, Y: rect.H})
		clipboard.Write(clipboard.FmtText, lines)
		sr.highlighter = nil
		count := bytes.Count(lines, []byte{'\n'}) + 1
		sr.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%d lines have been copied to clipboard", count))

	default:
		panic(fmt.Sprintf("TODO: unmet conditional '%d'", hl.mode))
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
	var rect *sdl.Rect

	switch sr.highlighter.mode {
	case _HIGHLIGHT_MODE_PNG:
		alphaBorder, alphaFill = 190, 64
		rect = &sdl.Rect{X: sr.highlighter.rect.X, Y: sr.highlighter.rect.Y, W: sr.highlighter.rect.W, H: sr.highlighter.rect.H}

	case _HIGHLIGHT_MODE_SQUARE:
		alphaBorder, alphaFill = 64, 0
		rect = &sdl.Rect{X: sr.highlighter.rect.X, Y: sr.highlighter.rect.Y, W: sr.highlighter.rect.W, H: sr.highlighter.rect.H}

	case _HIGHLIGHT_MODE_LINE_RANGE, _HIGHLIGHT_MODE_FULL_LINES:
		return

	default:

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
}

func isCellHighlighted(sr *sdlRender, rect *sdl.Rect) bool {
	if sr.highlighter == nil || sr.highlighter.button == 0 {
		return false
	}

	hlRect := *sr.highlighter.rect
	if sr.highlighter.mode != _HIGHLIGHT_MODE_LINE_RANGE {
		normaliseRect(&hlRect)
	}
	runeCell := sr.rectPxToCells(rect)
	hlCell := sr.rectPxToCells(&hlRect)

	switch sr.highlighter.mode {
	case _HIGHLIGHT_MODE_FULL_LINES:
		return runeCell.Y >= hlCell.Y && runeCell.Y <= hlCell.H

	case _HIGHLIGHT_MODE_LINE_RANGE:
		switch {
		case hlCell.H < hlCell.Y: // select up
			// start multiline
			return ((runeCell.X <= hlCell.X && runeCell.Y == hlCell.Y) ||
				// middle multiline
				(runeCell.Y < hlCell.Y && runeCell.Y > hlCell.H) ||
				// end multiline
				(runeCell.X >= hlCell.W && runeCell.Y == hlCell.H))

		case hlCell.Y == hlCell.H:
			// midline
			if hlCell.W < hlCell.X { //backwards
				return runeCell.X <= hlCell.X && runeCell.X >= hlCell.W && runeCell.Y == hlCell.Y
			} else { // forwards
				return runeCell.X >= hlCell.X && runeCell.X <= hlCell.W && runeCell.Y == hlCell.Y
			}

		default: // select down
			// start multiline
			return ((runeCell.X >= hlCell.X && runeCell.Y == hlCell.Y) ||
				// middle multiline
				(runeCell.Y > hlCell.Y && runeCell.Y < hlCell.H) ||
				// end multiline
				(runeCell.X <= hlCell.W && runeCell.Y == hlCell.H))
		}

	case _HIGHLIGHT_MODE_SQUARE:
		return runeCell.X >= hlCell.X && runeCell.X <= hlCell.W &&
			runeCell.Y >= hlCell.Y && runeCell.Y <= hlCell.H

	default:
		return false
	}
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
