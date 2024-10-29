package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	cells := term.visibleScreen()

	var err error
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
				err = term.renderer.PrintCell(&cells[pos.Y][pos.X], pos)
				if err != nil {
					log.Printf("ERROR: error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) Render()", pos.X, pos.Y, string(cells[pos.Y][pos.X].Char), err.Error())
				}
			}
		}
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) _blinkCursor() {
	if term._hideCursor {
		return
	}

	// copy cell
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
	err := term.renderer.PrintCell(cell, term.curPos())
	if err != nil {
		log.Printf("ERROR: error in %s [cursorBlink]: %s", "(t *Term) _blinkCursor()", err.Error())
	}
}
