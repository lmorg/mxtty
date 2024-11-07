package elementCsv

import (
	"fmt"
	"strings"
)

const _ROW_ID = "rowid"

func (el *ElementCsv) runQuery() error {
	where := el.filter
	if where != "" {
		where = "WHERE " + where
	}

	orderBy := _ROW_ID
	var sql string
	if el.orderByIndex > 0 {
		orderBy = string(el.headings[el.orderByIndex-1])
		sql = sqlSelect[el.isNumber[el.orderByIndex-1]]
	} else {
		sql = sqlSelect[selectNumeric]
	}

	query := fmt.Sprintf(sql, el.name, where, orderBy, orderByStr[el.orderDesc], el.size.Y-1)

	dbRows, err := el.db.Query(query)
	if err != nil {
		return fmt.Errorf("cannot query table: %s\nSQL: %s", err.Error(), query)
	}

	var (
		table []string
		width = make([]int, len(el.headings))
		rows  [][]string
		l     = len(el.headings)
	)

	for dbRows.Next() {
		row := make([]string, l)
		slice := _strToAnyPtr(&row, l)

		err = dbRows.Scan(slice...)
		if err != nil {
			return err
		}

		for i := range row {
			if len([]rune(row[i])) > width[i] {
				width[i] = len([]rune(row[i]))
			}
		}

		rows = append(rows, row)
	}

	boundaries := make([]int32, len(el.headings))
	var boundaryPos int32
	// check if rows are smaller than headings
	// also lets do the boundaries for the table lines
	for i := range el.headings {
		if len(el.headings[i]) > width[i] {
			width[i] = len(el.headings[i])
		}
		boundaryPos += int32(width[i]) + 2
		boundaries[i] = boundaryPos
	}

	for _, row := range rows {
		var s string
		for i := range row {
			s += fmt.Sprintf(" %s%s ", row[i], strings.Repeat(" ", width[i]-len([]rune(row[i]))))
		}

		table = append(table, s)
	}

	var top string
	for i := range el.headings {
		top += fmt.Sprintf(" %s%s ", string(el.headings[i]), strings.Repeat(" ", width[i]-len(el.headings[i])))
	}

	if err = dbRows.Err(); err != nil {
		return fmt.Errorf("cannot retrieve rows: %s", err.Error())
	}

	err = dbRows.Close()
	if err != nil {
		return err
	}

	el.table = make([][]rune, len(table))
	for i := range table {
		el.table[i] = []rune(table[i])
	}
	el.top = []rune(top)
	el.width = width
	el.boundaries = boundaries

	return nil
}

func _strToAnyPtr(s *[]string, max int) []any {
	slice := make([]interface{}, max)
	for i := range slice {
		slice[i] = &(*s)[i]
	}

	return slice
}
