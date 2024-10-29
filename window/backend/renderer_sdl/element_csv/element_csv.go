package elementCsv

import (
	"database/sql"

	"github.com/lmorg/mxtty/types"
)

type ElementCsv struct {
	renderer types.Renderer
	size     *types.XY
	//parameters parametersT
	db *sql.DB
	//filter       string
	//orderBy      int        // row
	//orderDesc    bool       // ASC or DESC
}

var arrowGlyph = map[bool]rune{
	true:  '↓',
	false: '↑',
}

func New(renderer types.Renderer) *ElementCsv {
	return &ElementCsv{renderer: renderer}
}

func (el *ElementCsv) Begin(apc *types.ApcSlice) {
	//apc.Parameters(&el.parameters)

	// initialise table
	//el.setName()
	//el.parameters.Format = strings.ToLower(el.parameters.Format)

	//el.orderBy = -1

	/*if el.parameters.Format == "" {
		el.beginTerm()
	} else {
		el.beginStruct()
	}*/

	err := el.createDb()
	if err != nil {
		panic(err)
	}
}

func (el *ElementCsv) End() *types.XY {
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
	err := el.ElementCsv(confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings)
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

func (el *ElementCsv) Insert(_ *types.ApcSlice) *types.XY {
	// not required for this element
	return nil
}

func (el *ElementCsv) Draw(rect *types.Rect) *types.XY {
	if el.size == nil {
		return nil
	}

	if el.parameters.Format == "" {
		return el.drawTerm(rect)
	}

	return el.drawStruct(rect)
}

func (el *ElementCsv) Close() {
	// clear memory (if required)
	el.db.Close()
}
