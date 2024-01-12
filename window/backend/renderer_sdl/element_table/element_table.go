package elementTable

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lmorg/murex/utils/json"
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
	size        *types.XY
	name        string
	db          *sql.DB
	colOffset   []int32             // [column]
	recOffset   []int32             // [column]
	_colOffset  int32               // counter
	headOffset  int32               // row
	_headOffset int32               // counter
	_stack      [][]*elTableRecordT // [row][column][cell]
	_tableCache [][]string          // [row][column]
	_sqlResult  []int               // [row]
}

func New(renderer types.Renderer) *ElementTable {
	return &ElementTable{renderer: renderer}
}

func (el *ElementTable) Begin(parameters types.ApcSlice) {
	// initialise table
	el.name = "TODO"
	el.hashName()

	el._colOffset = -1
	el._stack = make([][]*elTableRecordT, 1)
	el._stack[0] = make([]*elTableRecordT, 1)
	el._stack[0][0] = new(elTableRecordT)

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
		//switch {
		//case len(el._stack[y]) == 1 && len(el._stack[y][0].cells) == 0:
		//	return

		//default:
		el._stack = append(el._stack, []*elTableRecordT{new(elTableRecordT)})
		el._colOffset = -1
		return
		//}
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

	switch y {
	case 0:
		if len(el._stack[y][x].cells) == 0 {
			el.colOffset = append(el.colOffset, el._colOffset)
		}
	case 1:
		if len(el._stack[y][x].cells) == 0 {
			el.recOffset = append(el.recOffset, el._colOffset)
		}
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
	el._tableCache[0][0] = "ROW_ID"
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

		fmt.Printf("---\n")
		fmt.Printf("%s (%d)\n", json.LazyLogging(el._stack[y]), len(el._stack[y][0].cells))
		fmt.Printf("%s\n", json.LazyLogging(el._tableCache[y]))
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

	el._sqlResult, err = el.runQuery("")
	if err != nil {
		panic(err)
	}

}

func (el *ElementTable) Draw(offset *types.XY) {
	if el.size == nil {
		return
	}

	var err error
	for x := range el._stack[el._sqlResult[0]] {
		for i, cell := range el._stack[0][x].cells {
			err = el.renderer.PrintRuneColour(
				cell.Char,
				offset.X+el.colOffset[x]+int32(i),
				offset.Y,
				types.SGR_COLOUR_CYAN,
				nil,
				cell.Sgr.Bitwise,
			)
			if err != nil {
				panic(err)
			}
		}
	}

	for i := 0; i < len(el._sqlResult); i++ {
		y := el._sqlResult[i]
		for x := range el._stack[y] {
			for col, cell := range el._stack[y][x].cells {
				err = el.renderer.PrintRuneColour(
					cell.Char,
					offset.X+el.recOffset[x]+int32(col),
					offset.Y+int32(i)+1,
					types.SGR_COLOUR_CYAN,
					nil,
					cell.Sgr.Bitwise,
				)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}

func (el *ElementTable) Close() {
	// clear memory (if required)
	el.db.Close()
}
