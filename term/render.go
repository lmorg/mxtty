package virtualterm

import (
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	if !term.visible {
		return
	}

	term._mutex.Lock()

	screen := term.visibleScreen()

	if !config.Config.Terminal.TypeFace.Ligatures || term._mouseButtonDown || term._searchHighlight {
		term._renderCells(screen)
	} else {
		term._renderLigatures(screen)
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) _renderCells(screen types.Screen) {
	pos := new(types.XY)
	elementStack := make(map[types.Element]bool) // no duplicates

	for ; pos.Y < term.size.Y; pos.Y++ {
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			switch {
			case screen[pos.Y].Cells[pos.X].Element != nil:
				_, ok := elementStack[screen[pos.Y].Cells[pos.X].Element]
				if !ok {
					elementStack[screen[pos.Y].Cells[pos.X].Element] = true
					offset := getElementXY(screen[pos.Y].Cells[pos.X].Char)
					screen[pos.Y].Cells[pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}

			case screen[pos.Y].Cells[pos.X].Char == 0:
				continue

			case screen[pos.Y].Cells[pos.X].Sgr == nil:
				continue

			default:
				if screen[pos.Y].Cells[pos.X].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCell(screen[pos.Y].Cells[pos.X], pos)
			}
		}
	}
}

func (term *Term) _renderLigatures(screen types.Screen) {
	var (
		pos          = new(types.XY)
		elementStack = make(map[types.Element]bool) // no duplicates
		hash         uint64
		defaultHash  = types.SGR_DEFAULT.HashValue()
	)

	for ; pos.Y < term.size.Y; pos.Y++ {
		if screen[pos.Y].Cells[0].Sgr == nil {
			hash = defaultHash
		} else {
			hash = screen[pos.Y].Cells[0].Sgr.HashValue()
		}

		var start int32
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			newHash := defaultHash // ^uint64(0)
			if screen[pos.Y].Cells[pos.X].Sgr != nil {
				newHash = screen[pos.Y].Cells[pos.X].Sgr.HashValue()
			}

			if screen[pos.Y].Cells[pos.X].Element != nil {
				_, ok := elementStack[screen[pos.Y].Cells[pos.X].Element]
				if !ok {
					elementStack[screen[pos.Y].Cells[pos.X].Element] = true
					offset := getElementXY(screen[pos.Y].Cells[pos.X].Char)
					screen[pos.Y].Cells[pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}
				newHash = defaultHash //^uint64(0)
			}

			if hash != newHash {
				if screen[pos.Y].Cells[start].Sgr != nil && screen[pos.Y].Cells[start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:pos.X], &types.XY{X: start, Y: pos.Y})
				hash = newHash
				start = pos.X
			}

			if screen[pos.Y].Cells[pos.X].Char == 0 || screen[pos.Y].Cells[pos.X].Element != nil {
				term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:pos.X], &types.XY{X: start, Y: pos.Y})
				start = pos.X + 1
			}
		}

		if start < pos.X {
			if screen[pos.Y].Cells[start].Sgr != nil && screen[pos.Y].Cells[start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
				continue // blink
			}
			term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:], &types.XY{X: start, Y: pos.Y})
		}
	}
}

func (term *Term) _blinkCursor() {
	if term._hideCursor {
		return
	}

	if term._slowBlinkState {
		term.renderer.DrawHighlightRect(term.curPos(), &types.XY{1, 1})
		term.renderer.DrawHighlightRect(term.curPos(), &types.XY{1, 1})
	}
}
