package elementTable

import (
	"log"
	"strconv"

	"github.com/lmorg/mxtty/types"
)

func (el *ElementTable) beginTerm() {
	el.colOffset = make([][]int32, 1)
	el._colOffset = -1
	el._stackTerm = make([][]*elTableRecordT, 1)
	el._stackTerm[0] = make([]*elTableRecordT, 1)
	el._stackTerm[0][0] = new(elTableRecordT)
}

func (el *ElementTable) readTerm(cell *types.Cell) {
	el._colOffset++
	y := len(el._stackTerm) - 1
	x := len(el._stackTerm[y]) - 1

	if cell == nil {
		el._stackTerm = append(el._stackTerm, []*elTableRecordT{new(elTableRecordT)})
		el.colOffset = append(el.colOffset, []int32{})
		el._colOffset = -1
		return
	}

	if cell.Char == ' ' /* || cell.Char == '\t' */ {
		switch {
		case y > 0 && x == len(el._stackTerm[0])-1:
			el._stackTerm[y][x].append(cell)

		case len(el._stackTerm[y][x].cells) != 0:
			el._stackTerm[y] = append(el._stackTerm[y], new(elTableRecordT))
		}
		return
	}

	if len(el._stackTerm[y][x].cells) == 0 {
		el.colOffset[len(el.colOffset)-1] = append(el.colOffset[len(el.colOffset)-1], el._colOffset)
	}

	el._stackTerm[y][x].append(cell)
}

func (el *ElementTable) endTerm() {
	if int(el.renderer.TermSize().Y-5) <= len(el._stackTerm) {
		el.size = &types.XY{
			X: el.renderer.TermSize().X,
			Y: el.renderer.TermSize().Y - 5,
		}
	} else {
		el.size = &types.XY{
			X: el.renderer.TermSize().X,
			Y: int32(len(el._stackTerm)),
		}
	}

	// initialise cache
	el._tableCache = append(el._tableCache, make([]string, len(el._stackTerm[0])+1))
	el._tableCache[0][0] = _ROW_ID
	for x := range el._stackTerm[0] {
		el._tableCache[0][x+1] = el._stackTerm[0][x].String
	}

	for y := 1; y < len(el._stackTerm); y++ {
		if len(el._stackTerm[y]) == 1 && el._stackTerm[y][0].String == "" {
			continue
		}
		el._tableCache = append(el._tableCache, make([]string, len(el._stackTerm[y])+1))
		el._tableCache[y][0] = strconv.Itoa(y)
		for x := range el._stackTerm[y] {
			el._tableCache[y][x+1] = el._stackTerm[y][x].String

		}
	}

	return
}

func (el *ElementTable) drawTerm(rect *types.Rect) *types.XY {
	var err error
	pos := new(types.XY)
	pos.Y = rect.Start.Y
	for x := range el._stackTerm[el._sqlResult[0]] {
		for i, cell := range el._stackTerm[0][x].cells {
			pos.X = rect.Start.X + el.colOffset[0][x] + int32(i)

			err = el.renderer.PrintCell(cell, pos)
			if err != nil {
				panic(err) // TODO: don't panic!
			}
		}
	}

	if el.orderBy > -1 {
		var orderGlyph = types.Cell{
			Sgr:  types.SGR_DEFAULT.Copy(),
			Char: arrowGlyph[el.orderDesc],
		}
		el.renderer.PrintCell(&orderGlyph, &types.XY{
			X: rect.Start.X + el.colOffset[0][el.orderBy] + int32(len(el._stackTerm[0][el.orderBy].cells)),
			Y: rect.Start.Y,
		})
	}

	var lineWrapping types.XY

	for i := 0; i < len(el._sqlResult); i++ {
		y := el._sqlResult[i]
		lineWrapping.X = 0
		pos.Y = rect.Start.Y + int32(i) + 1 + lineWrapping.Y
		if pos.Y > rect.End.Y {
			break
		}

		for x := range el._stackTerm[y] {
			var colOffset int32
			if len(el.colOffset[y]) <= x {
				// this works around a bug in tmux
				colOffset = el.colOffset[y][len(el.colOffset[y])-1]
			} else {
				colOffset = el.colOffset[y][x]
			}
			for col, cell := range el._stackTerm[y][x].cells {
				pos.X = rect.Start.X + colOffset + int32(col) - lineWrapping.X

				if pos.X >= el.renderer.TermSize().X {
					lineWrapping.X += el.renderer.TermSize().X
					pos.X = 0
					lineWrapping.Y++
					pos.Y++
				}

				if pos.Y > rect.End.Y {
					return nil
				}

				err = el.renderer.PrintCell(cell, pos)
				if err != nil {
					log.Printf("ERROR: cannot write cell: %s", err.Error())
				}
			}
		}
	}

	return nil
}
