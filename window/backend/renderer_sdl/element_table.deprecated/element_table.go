package elementTable

import (
	"database/sql"
	"strings"

	"github.com/lmorg/mxtty/types"
)

const (
	_ROW_ID = "___mxapc_row_id"
)

type parametersT struct {
	Name        string
	Format      string
	HeadMissing bool
	HeadOffset  int
}

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
	parameters   parametersT
	db           *sql.DB
	filter       string
	colOffset    [][]int32           // [row][column]
	_colOffset   int32               // counter
	_headOffset  int32               // counter
	_stackTerm   [][]*elTableRecordT // [row][column][cell]
	_stackStruct []rune
	_tableCache  [][]string // [row][column]
	_sqlResult   []int      // [row]
	orderBy      int        // row
	orderDesc    bool       // ASC or DESC
}

var arrowGlyph = map[bool]rune{
	true:  '↓',
	false: '↑',
}

func New(renderer types.Renderer) *ElementTable {
	return &ElementTable{renderer: renderer}
}

func (el *ElementTable) Begin(apc *types.ApcSlice) {
	apc.Parameters(&el.parameters)

	// initialise table
	el.setName()
	el.parameters.Format = strings.ToLower(el.parameters.Format)

	el.orderBy = -1

	if el.parameters.Format == "" {
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
	if el.parameters.Format == "" {
		el.readTerm(cell)
	} else {
		el.readStruct(cell)
	}
}

func (el *ElementTable) End() *types.XY {
	if el.parameters.Format == "" {
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
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot create sqlite3 table: "+err.Error())
		return &types.XY{}
	}

	err = el.runQuery()
	if err != nil {
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot query sqlite3 table: "+err.Error())
		return &types.XY{}
	}

	return nil
}

func (el *ElementTable) Insert(_ *types.ApcSlice) *types.XY {
	// not required for this element
	return nil
}

func (el *ElementTable) Draw(rect *types.Rect) *types.XY {
	if el.size == nil {
		return nil
	}

	if el.parameters.Format == "" {
		return el.drawTerm(rect)
	}

	return el.drawStruct(rect)
}

func (el *ElementTable) Close() {
	// clear memory (if required)
	el.db.Close()
}
