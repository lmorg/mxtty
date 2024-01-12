package elementTable

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	sqlCreateTable  = `CREATE TABLE IF NOT EXISTS '%s' (%s);`
	sqlInsertRecord = `INSERT INTO '%s' VALUES (%s);`
	sqlSelect       = `SELECT ROW_ID from '%s';`
)

func (el *ElementTable) hashName() {
	if el.name == "" {
		el.name = time.Now().String()
	}

	el.name = hex.EncodeToString([]byte(el.name))
}

func (el *ElementTable) createDb() error {
	var err error
	el.db, err = sql.Open(driverName, ":memory:")
	if err != nil {
		return fmt.Errorf("could not open database: %s", err.Error())
	}

	return nil
}

func (el *ElementTable) _openTable(headings []string) (*sql.Tx, error) {
	var err error

	if len(headings) == 0 {
		return nil, fmt.Errorf("cannot create table '%s': no titles supplied", el.name)
	}

	var sHeadings string
	for i := range headings {
		sHeadings += fmt.Sprintf(`"%s" NUMERIC,`, headings[i])
	}
	sHeadings = sHeadings[:len(sHeadings)-1]

	query := fmt.Sprintf(sqlCreateTable, el.name, sHeadings)
	_, err = el.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("could not create table '%s': %s\n%s", el.name, err.Error(), query)
	}

	tx, err := el.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %s", err.Error())
	}

	return tx, nil
}

func (el *ElementTable) _insertRecords(tx *sql.Tx, records []interface{}) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to insert into transaction on table %s", el.name)
	}

	values, err := createValues(len(records))
	if err != nil {
		return fmt.Errorf("cannot insert records into transaction on table %s: %s", el.name, err.Error())
	}

	_, err = tx.Exec(fmt.Sprintf(sqlInsertRecord, el.name, values), records...)
	if err != nil {
		return fmt.Errorf("cannot insert records into transaction on table %s: %s", el.name, err.Error())
	}

	return nil
}

func createValues(length int) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("no records to insert")
	}

	values := strings.Repeat("?,", length)
	values = values[:len(values)-1]

	return values, nil
}
