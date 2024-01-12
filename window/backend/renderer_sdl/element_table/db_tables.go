package elementTable

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lmorg/murex/utils/humannumbers"
)

func (el *ElementTable) createTable(confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings bool) error {
	return el._createTable_SliceSliceString(confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings)
	/*switch v := v.(type) {
	case [][]string:
		return el.createTable_SliceSliceString(v, confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings)

	case []interface{}:
		table := make([][]string, len(v)+1)
		i := 1
		err := types.MapToTable(v, func(s []string) error {
			table[i] = s
			i++
			return nil
		})
		if err != nil {
			return err
		}
		return el.createTable_SliceSliceString(table, confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings)

	default:
		return fmt.Errorf("unable to convert the following data structure into a table '%s': %T", el.name, v)
	}*/
}

func (el *ElementTable) _createTable_SliceSliceString(confFailColMismatch, confMergeTrailingColumns, confTableIncHeadings bool) error {
	var (
		tx       *sql.Tx
		err      error
		headings []string
		nRow     int
	)

	if confTableIncHeadings {
		headings = make([]string, len(el._tableCache[0]))
		for i := range headings {
			headings[i] = fmt.Sprint(el._tableCache[0][i])
		}
		tx, err = el._openTable(headings)
		if err != nil {
			return err
		}
		nRow = 1

	} else {
		headings = make([]string, len(el._tableCache[0]))
		for i := range headings {
			headings[i] = humannumbers.ColumnLetter(i)
		}
		tx, err = el._openTable(headings)
		if err != nil {
			return err
		}

		slice := stringToInterfaceTrim(el._tableCache[0], len(el._tableCache))
		err = el._insertRecords(tx, slice)
		if err != nil {
			return fmt.Errorf("unable to insert headings into sqlite3: %s", err.Error())
		}
		nRow = 1
	}

	for ; nRow < len(el._tableCache); nRow++ {
		if len(el._tableCache[nRow]) != len(headings) && confFailColMismatch {
			return fmt.Errorf("table rows contain a different number of columns to table headings\n%d: %s", nRow, el._tableCache[nRow])
		}

		if confMergeTrailingColumns {
			slice := stringToInterfaceMerge(el._tableCache[nRow], len(headings))
			err = el._insertRecords(tx, slice)
			if err != nil {
				return fmt.Errorf("%s\n%d: %s", err.Error(), nRow, el._tableCache[nRow])
			}
		} else {
			slice := stringToInterfaceTrim(el._tableCache[nRow], len(headings))
			err = el._insertRecords(tx, slice)
			if err != nil {
				return fmt.Errorf("%s\n%d: %s", err.Error(), nRow, el._tableCache[nRow][:len(headings)-1])
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit sqlite3 transaction: %s", err.Error())
	}

	return nil
}

func (el *ElementTable) runQuery(parameters string) ([]int, error) {
	query := fmt.Sprintf(sqlSelect, el.name)

	rows, err := el.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("cannot query table: %s\nSQL: %s", err.Error(), query)
	}

	var (
		table []int
		s     string
		i     int
	)

	for rows.Next() {

		err = rows.Scan(&s)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve rows: %s", err.Error())
		}

		i, err = strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve rows: %s", err.Error())
		}

		table = append(table, i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot retrieve rows: %s", err.Error())
	}

	return table, nil
}
