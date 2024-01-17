package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	var err error
	elementLookup := make(map[types.Element]*types.Rect)
	pos := new(types.XY)

	for ; int(pos.Y) < len(*term.cells); pos.Y++ {
		for pos.X = 0; int(pos.X) < len((*term.cells)[pos.Y]); pos.X++ {
			switch {
			case (*term.cells)[pos.Y][pos.X].Sgr == nil:
				continue

			case (*term.cells)[pos.Y][pos.X].Element != nil:
				rect, ok := elementLookup[(*term.cells)[pos.Y][pos.X].Element]
				if ok { // update rect
					rect.End.X, rect.End.Y = pos.X, pos.Y
				} else { // create rect
					elementLookup[(*term.cells)[pos.Y][pos.X].Element] = &types.Rect{
						Start: &types.XY{X: pos.X, Y: pos.Y},
						End:   &types.XY{X: pos.X, Y: pos.Y},
					}
				}

			case (*term.cells)[pos.Y][pos.X].Char == 0:
				continue

			default:
				err = term.renderer.PrintCell(&(*term.cells)[pos.Y][pos.X], pos)
				if err != nil {
					log.Printf("ERROR: error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) Render()", pos.X, pos.Y, string((*term.cells)[pos.Y][pos.X].Char), err.Error())
				}
			}
		}
	}

	for el, rect := range elementLookup {
		el.Draw(rect)
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) _blinkCursor() {
	if term._hideCursor {
		return
	}

	// copy cell
	cell := term.copyCell(term.cell())

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
	err := term.renderer.PrintCell(cell, &term.curPos)
	if err != nil {
		log.Printf("ERROR: error in %s [cursorBlink]: %s", "(t *Term) _blinkCursor()", err.Error())
	}
}
