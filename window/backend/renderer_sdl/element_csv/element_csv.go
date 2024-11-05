package elementCsv

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"

	"github.com/lmorg/mxtty/types"
)

type ElementCsv struct {
	renderer   types.Renderer
	size       *types.XY
	headings   []string // columns
	table      []string // rendered rows
	top        string   // rendered headings
	width      []int    // columns
	boundaries []int32  // column lines

	//parameters parametersT

	name   string
	buf    []rune
	lines  int32
	notify types.Notification

	db   *sql.DB
	dbTx *sql.Tx

	filter       string
	orderByIndex int  // row
	orderDesc    bool // ASC or DESC
}

var arrowGlyph = map[bool]rune{
	false: '↑',
	true:  '↓',
}

const notifyLoading = "Loading CSV. Line %d..."

func New(renderer types.Renderer) *ElementCsv {
	el := &ElementCsv{renderer: renderer}

	el.notify = renderer.DisplaySticky(types.NOTIFY_INFO, fmt.Sprintf(notifyLoading, el.lines))

	err := el.createDb()
	if err != nil {
		panic(err)
	}

	return el
}

func (el *ElementCsv) Write(r rune) error {
	el.buf = append(el.buf, r)
	if r == '\n' {
		el.lines++
		el.notify.SetMessage(fmt.Sprintf(notifyLoading, el.lines))
	}
	return nil
}

func (el *ElementCsv) Generate(apc *types.ApcSlice) error {
	defer el.notify.Close()

	buf := bytes.NewBufferString(string(el.buf))
	r := csv.NewReader(buf)
	recs, err := r.ReadAll()
	if err != nil {
		return err
	}

	if len(recs) < 2 {
		return fmt.Errorf("too few rows") // TODO: this shouldn't error
	}

	err = el.createTable(recs[0])
	if err != nil {
		return err
	}
	el.headings = recs[0]

	for row := 1; row < len(recs); row++ {
		err = el.insertRecords(recs[row])
		if err != nil {
			return err
		}
	}

	if el.dbTx.Commit() != nil {
		return err
	}

	el.size = el.renderer.GetTermSize()
	if el.size.Y > el.lines {
		el.size.Y = el.lines
	}

	err = el.runQuery()
	if err != nil {
		return err
	}

	return nil
}

func (el *ElementCsv) Size() *types.XY {
	return el.size
}

func (el *ElementCsv) Rune(pos *types.XY) rune {
	if pos.Y == 0 {
		if int(pos.X) >= len(el.top) {
			return ' '
		}
		return rune(el.top[pos.X])
	}

	if int(pos.X) >= len(el.table[pos.Y-1]) {
		return ' '
	}

	return rune(el.table[pos.Y-1][pos.X])
}

func (el *ElementCsv) Draw(size *types.XY, pos *types.XY) {
	var err error

	cell := &types.Cell{Sgr: &types.Sgr{}}
	cell.Sgr.Reset()
	relPos := &types.XY{X: pos.X, Y: pos.Y}

	cell.Sgr.Bitwise |= types.SGR_INVERT
	for i := range el.top {

		cell.Char = rune(el.top[i])
		err = el.renderer.PrintCell(cell, relPos)
		if err != nil {
			panic(err)
		}
		relPos.X++
	}

	switch el.orderByIndex {
	case 0:
		goto skipOrderGlyph

	case 1:
		relPos.X = 0

	default:
		relPos.X = el.boundaries[el.orderByIndex-2]
	}

	cell.Char = arrowGlyph[el.orderDesc]
	err = el.renderer.PrintCell(cell, relPos)
	if err != nil {
		panic(err)
	}

skipOrderGlyph:

	relPos.Y++

	cell.Sgr.Bitwise ^= types.SGR_INVERT
	for y := int32(0); y < el.size.Y-1 && int(y) < len(el.table); y++ {
		relPos.X = pos.X
		for x := int32(0); x < el.size.X && int(x) < len(el.table[y]); x++ {
			cell.Char = rune(el.table[y][x])
			err = el.renderer.PrintCell(cell, relPos)
			if err != nil {
				panic(err)
			}
			relPos.X++
		}
		relPos.Y++
	}

	el.renderer.DrawTable(pos, int32(len(el.table)), el.boundaries)
}

func (el *ElementCsv) Close() {
	// clear memory (if required)
	el.db.Close()
}
