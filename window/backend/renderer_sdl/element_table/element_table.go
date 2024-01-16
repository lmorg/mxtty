package elementTable

import (
	"database/sql"
	"strings"

	"github.com/lmorg/mxtty/types"
)

const (
	_KEY_TABLE_NAME  = "name"
	_KEY_FORMAT      = "format"
	_KEY_NO_HEADINGS = "add_headings"
	_KEY_HEAD_OFFSET = "col_offset"
	_ROW_ID          = "___mxapc_row_id"
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
	renderer     types.Renderer
	name         string
	size         *types.XY
	apc          *types.ApcSlice
	db           *sql.DB
	filter       string
	_paramFormat string
	colOffset    [][]int32           // [row][column]
	_colOffset   int32               // counter
	headOffset   int32               // row
	_headOffset  int32               // counter
	_stackTerm   [][]*elTableRecordT // [row][column][cell]
	_stackStruct []rune
	_tableCache  [][]string // [row][column]
	_sqlResult   []int      // [row]
	orderBy      int        // row
	orderDesc    bool       // ASC or DESC
}

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
	el._paramFormat = strings.ToLower(el.apc.Parameter(_KEY_FORMAT))

	el.orderBy = -1

	if el._paramFormat == "" {
		el.beginTerm()
	} else {
		el.beginStruct()
	}

	err := el.createDb()
	if err != nil {
		panic(err)
	}
}

func (el *ElementTable) ReadCell(cell *types.Cell) {
	if el._paramFormat == "" {
		el.readTerm(cell)
	} else {
		el.readStruct(cell)
	}
}

func (el *ElementTable) End() {
	if el._paramFormat == "" {
		el.endTerm()
	} else {
		el.endStruct()
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

	if el._paramFormat == "" {
		el.drawTerm(rect)
	} else {
		el.drawStruct(rect)
	}
}

func (el *ElementTable) Close() {
	// clear memory (if required)
	el.db.Close()
}
