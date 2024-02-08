package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	var cells = *term.cells
	if term._scrollOffset != 0 {
		// render scrollback buffer
		start := len(term._scrollBuf) - term._scrollOffset
		cells = term._scrollBuf[start:]
		if len(cells) < int(term.size.Y) {
			cells = append(cells, term._normBuf...)
		}
	}

	var err error
	elementStack := make(map[types.Element]*types.Rect)
	pos := new(types.XY)

	for ; pos.Y < term.size.Y; pos.Y++ {
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			switch {
			case cells[pos.Y][pos.X].Element != nil:
				rect, ok := elementStack[cells[pos.Y][pos.X].Element]
				if !ok { // create rect
					elementStack[cells[pos.Y][pos.X].Element] = &types.Rect{
						Start: &types.XY{X: pos.X, Y: pos.Y},
						End:   &types.XY{X: pos.X, Y: pos.Y},
					}
				} else { // update rect
					rect.End.X, rect.End.Y = pos.X, pos.Y
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

	for el, rect := range elementStack {
		size := el.Draw(rect)
		if size != nil {
			term._elementResizeGrow(el, rect.Start, size)
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
