package elementTable

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/lmorg/mxtty/types"
)

type elTableRecordT struct {
	cells  []*types.Cell
	String string
}

func (rec *elTableRecordT) append(cell *types.Cell) {
	rec.cells = append(rec.cells, cell)
	rec.String += string(cell.Char)
	//fmt.Println(rec.String)
}

type ElementTable struct {
	renderer    types.Renderer
	name        string
	size        *types.XY
	apc         *types.ApcSlice
	db          *sql.DB
	colOffset   [][]int32           // [row][column]
	_colOffset  int32               // counter
	headOffset  int32               // row
	_headOffset int32               // counter
	_stack      [][]*elTableRecordT // [row][column][cell]
	_tableCache [][]string          // [row][column]
	_sqlResult  []int               // [row]
	orderBy     int                 // row
	orderDesc   bool                // ASC or DESC
}

/*var arrowGlyph = map[bool]rune{
	false: 'ˇ',
	true:  '^',
}*/

var arrowGlyph = map[bool]rune{
	false: '↓',
	true:  '↑',
}

func New(renderer types.Renderer) *ElementTable {
	return &ElementTable{renderer: renderer}
}

func (el *ElementTable) Begin(apc *types.ApcSlice) {
	el.apc = apc

	// initialise table
	el.setName()

	el.colOffset = make([][]int32, 1)
	el._colOffset = -1
	el._stack = make([][]*elTableRecordT, 1)
	el._stack[0] = make([]*elTableRecordT, 1)
	el._stack[0][0] = new(elTableRecordT)
	el.orderBy = -1

	err := el.createDb()
	if err != nil {
		panic(err)
	}
}

func (el *ElementTable) ReadCell(cell *types.Cell) {
	// read cell
	el._colOffset++
	y := len(el._stack) - 1
	x := len(el._stack[y]) - 1

	if cell == nil {
		el._stack = append(el._stack, []*elTableRecordT{new(elTableRecordT)})
		el.colOffset = append(el.colOffset, []int32{})
		el._colOffset = -1
		return
	}

	if cell.Char == ' ' /* || cell.Char == '\t' */ {
		switch {
		case y > 0 && x == len(el._stack[0])-1:
			el._stack[y][x].append(cell)

		case len(el._stack[y][x].cells) != 0:
			el._stack[y] = append(el._stack[y], new(elTableRecordT))
		}
		return
	}

	if len(el._stack[y][x].cells) == 0 {
		el.colOffset[len(el.colOffset)-1] = append(el.colOffset[len(el.colOffset)-1], el._colOffset)
	}

	el._stack[y][x].append(cell)
}

func (el *ElementTable) End() {
	if int(el.renderer.TermSize().Y-5) <= len(el._stack) {
		el.size = &types.XY{
			X: el.renderer.TermSize().X,
			Y: el.renderer.TermSize().Y - 5,
		}
	} else {
		el.size = &types.XY{
			X: el.renderer.TermSize().X,
			Y: int32(len(el._stack)),
		}
	}

	// initialise cache
	el._tableCache = append(el._tableCache, make([]string, len(el._stack[0])+1))
	el._tableCache[0][0] = "___mxapc_row_id"
	for x := range el._stack[0] {
		el._tableCache[0][x+1] = el._stack[0][x].String
	}

	for y := 1; y < len(el._stack); y++ {
		if len(el._stack[y]) == 1 && el._stack[y][0].String == "" {
			continue
		}
		el._tableCache = append(el._tableCache, make([]string, len(el._stack[y])+1))
		el._tableCache[y][0] = strconv.Itoa(y)
		for x := range el._stack[y] {
			el._tableCache[y][x+1] = el._stack[y][x].String

		}
	}

	// initialise db
	var (
		confFailColMismatch      = false
		confMergeTrailingColumns = false
		confTableIncHeadings     = true
	)
	err := el.createTable(confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings)
	if err != nil {
		panic(err)
	}

	el._sqlResult, err = el.runQuery()
	if err != nil {
		panic(err)
	}
}

func (el *ElementTable) Draw(rect *types.Rect) {
	if el.size == nil {
		return
	}

	//fmt.Printf("%s\n", json.LazyLogging(el._sqlResult))

	var err error
	pos := new(types.XY)
	pos.Y = rect.Start.Y
	for x := range el._stack[el._sqlResult[0]] {
		for i, cell := range el._stack[0][x].cells {
			pos.X = rect.Start.X + el.colOffset[0][x] + int32(i)

			err = el.renderer.PrintCell(cell, pos)
			if err != nil {
				panic(err)
			}
		}
	}

	if el.orderBy > -1 {
		var orderGlyph = types.Cell{
			Sgr:  types.SGR_DEFAULT.Copy(),
			Char: arrowGlyph[el.orderDesc],
		}
		el.renderer.PrintCell(&orderGlyph, &types.XY{
			X: rect.Start.X + el.colOffset[0][el.orderBy] - 1,
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

		for x := range el._stack[y] {

			for col, cell := range el._stack[y][x].cells {
				pos.X = rect.Start.X + el.colOffset[y][x] + int32(col) - lineWrapping.X

				if pos.X >= el.renderer.TermSize().X {
					lineWrapping.X += el.renderer.TermSize().X
					pos.X = 0
					lineWrapping.Y++
					pos.Y++
				}

				if pos.Y > rect.End.Y {
					return
				}

				err = el.renderer.PrintCell(cell, pos)
				if err != nil {
					log.Printf("ERROR: cannot write cell: %s", err.Error())
				}
			}
		}
	}
}

func (el *ElementTable) Close() {
	// clear memory (if required)
	el.db.Close()
}
