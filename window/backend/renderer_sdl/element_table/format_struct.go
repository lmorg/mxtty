package elementTable

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/lmorg/mxtty/types"
)

func (el *ElementTable) beginStruct() {
}

func (el *ElementTable) readStruct(cell *types.Cell) {
	if cell == nil {
		el._stackStruct = append(el._stackStruct, '\n')
		return
	}
	el._stackStruct = append(el._stackStruct, cell.Char)
}

func (el *ElementTable) endStruct() {
	var err error
	switch el._paramFormat {
	case "csv":
		err = el.parseCsv()
	}

	if err != nil {
		//log.Printf("ERROR: cannot parse %s: %s", el._paramFormat, err.Error())
		el.renderer.DisplayNotification(types.NOTIFY_ERROR,
			fmt.Sprintf("Cannot parse %s: %s", el._paramFormat, err.Error()))
	}

	return
}

func (el *ElementTable) parseCsv() error {
	r := strings.NewReader(string(el._stackStruct))
	csvReader := csv.NewReader(r)

	table, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if el.apc.Parameter(_KEY_HAS_HEADINGS) == "false" {
		//table = append([][]string{},table...)
		panic("TODO")
	}

	table[0] = append([]string{_ROW_ID}, table[0]...)
	for i := 1; i < len(table); i++ {
		table[i] = append([]string{strconv.Itoa(i)}, table[i]...)
	}

	el._tableCache = table

	return nil
}

func (el *ElementTable) drawStruct(rect *types.Rect) *types.XY {
	el.colOffset = [][]int32{{}}

	var err error
	var b []byte
	buf := bytes.NewBuffer(b)

	w := tabwriter.NewWriter(buf, 0, 0, 4, 0, 0)
	for i := range el._tableCache {
		_, err = w.Write([]byte(strings.Join(el._tableCache[i], "\t")))
		if err != nil {
			el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot output table: "+err.Error())
			return nil
		}
	}
	w.Flush()

	var (
		s         = buf.String()
		pos       = &types.XY{X: -1, Y: rect.Start.Y}
		last      rune
		heading   = true
		sortGlyph = &types.XY{X: -1, Y: rect.Start.Y}
		cell      = &types.Cell{Sgr: types.SGR_DEFAULT.Copy()}
	)

	for _, r := range s {
		pos.X++
		if pos.X >= el.renderer.TermSize().X {
			pos.X = -1
			pos.Y++
		}
		if pos.Y > rect.End.Y {
			break
		}

		switch r {
		case '\n':
			pos.X = -1
			pos.Y++
			heading = false

		case 0:
			if last != 0 && len(el.colOffset[0]) == el.orderBy {
				sortGlyph.X, sortGlyph.Y = pos.X, pos.Y
			}

		default:
			if last == 0 && heading {
				el.colOffset[0] = append(el.colOffset[0], pos.X)
			}
			pos.X++
			cell.Char = r
			el.renderer.PrintCell(cell, pos)
		}

		last = r
	}

	if el.orderBy > -1 {
		cell.Char = arrowGlyph[el.orderDesc]
		el.renderer.PrintCell(cell, sortGlyph)
	}

	return nil
}
