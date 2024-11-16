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

	cells := term.visibleScreen()

	if !config.Config.Terminal.TypeFace.Ligatures || term._mouseButtonDown {
		term._renderCells(cells)
	} else {
		term._renderLigatures(cells)
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) _renderCells(cells [][]types.Cell) {
	pos := new(types.XY)
	elementStack := make(map[types.Element]bool) // no duplicates

	for ; pos.Y < term.size.Y; pos.Y++ {
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			switch {
			case cells[pos.Y][pos.X].Element != nil:
				_, ok := elementStack[cells[pos.Y][pos.X].Element]
				if !ok {
					elementStack[cells[pos.Y][pos.X].Element] = true
					offset := getElementXY(cells[pos.Y][pos.X].Char)
					cells[pos.Y][pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}

			case cells[pos.Y][pos.X].Char == 0:
				continue

			case cells[pos.Y][pos.X].Sgr == nil:
				continue

			default:
				if cells[pos.Y][pos.X].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCell(&cells[pos.Y][pos.X], pos)
			}
		}
	}
}

func (term *Term) _renderLigatures(cells [][]types.Cell) {
	var (
		pos          = new(types.XY)
		elementStack = make(map[types.Element]bool) // no duplicates
		hash         uint64
		defaultHash  = types.SGR_DEFAULT.HashValue()
	)

	for ; pos.Y < term.size.Y; pos.Y++ {
		if cells[pos.Y][0].Sgr == nil {
			hash = defaultHash
		} else {
			hash = cells[pos.Y][0].Sgr.HashValue()
		}

		var start int32
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			newHash := ^uint64(0)
			if cells[pos.Y][pos.X].Sgr != nil {
				newHash = cells[pos.Y][pos.X].Sgr.HashValue()
			}

			if cells[pos.Y][pos.X].Element != nil {
				_, ok := elementStack[cells[pos.Y][pos.X].Element]
				if !ok {
					elementStack[cells[pos.Y][pos.X].Element] = true
					offset := getElementXY(cells[pos.Y][pos.X].Char)
					cells[pos.Y][pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}
				newHash = ^uint64(0)
			}

			if hash != newHash {
				if cells[pos.Y][start].Sgr != nil && cells[pos.Y][start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCellBlock(cells[pos.Y][start:pos.X], &types.XY{X: start, Y: pos.Y})
				hash = newHash
				start = pos.X
			}

			if cells[pos.Y][pos.X].Char == 0 || cells[pos.Y][pos.X].Element != nil {
				start = pos.X + 1
			}
		}

		if start < pos.X {
			if cells[pos.Y][start].Sgr != nil && cells[pos.Y][start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
				continue // blink
			}
			term.renderer.PrintCellBlock(cells[pos.Y][start:], &types.XY{X: start, Y: pos.Y})
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
