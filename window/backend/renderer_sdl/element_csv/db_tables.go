package elementCsv

/*
func (el *ElementTable) runQuery() error {
	where := el.filter
	if where != "" {
		where = "WHERE " + where
	}

	orderBy := _ROW_ID
	if el.orderBy > -1 {
		orderBy = el._stackTerm[0][el.orderBy].String
	}

	query := fmt.Sprintf(sqlSelect, _ROW_ID, el.name, where, orderBy, orderByStr[el.orderDesc])

	//log.Printf("DEBUG: SQL query = %s", query)

	rows, err := el.db.Query(query)
	if err != nil {
		return fmt.Errorf("cannot query table: %s\nSQL: %s", err.Error(), query)
	}

	var (
		table []int
		s     string
		i     int
	)

	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			return fmt.Errorf("cannot retrieve rows: %s", err.Error())
		}

		i, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("cannot retrieve rows: %s", err.Error())
		}

		table = append(table, i)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("cannot retrieve rows: %s", err.Error())
	}

	//log.Printf("DEBUG: %s", json.LazyLogging(table))

	rows.Close()

	el._sqlResult = table
	return nil
}
*/
