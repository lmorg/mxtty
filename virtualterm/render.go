package virtualterm

import (
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	cells := term.visibleScreen()

	term._blinkCursor()

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
				newHash := cells[pos.Y][pos.X].Sgr.HashValue()
				if hash != newHash {
					term.renderer.PrintCellBlock(cells[pos.Y][start:pos.X], &types.XY{X: start, Y: pos.Y})
					hash = newHash
					start = pos.X
				}
			}
		}

		if start < pos.X {
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
	}

	/*// copy cell
	cell := term.copyCurrentCell(term.currentCell())

	// format cell
	if cell.Char == 0 {
		cell.Char = ' '
		cell.Sgr.Fg, cell.Sgr.Bg = types.BlinkColour[true], types.BlinkColour[false]
		cell.Sgr.Bitwise = 0
	} else {
		cell.Sgr.Bg = term.sgr.Bg
	}

	if term._slowBlinkState {
		cell.Sgr.Fg, cell.Sgr.Bg = cell.Sgr.Bg, cell.Sgr.Fg
	}

	// print cell
	term.renderer.PrintCell(cell, term.curPos())*/
}
