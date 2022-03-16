package td

import (
	"bytes"
	"database/sql"
)

type TableAction interface {
	Create()
	Query()
	Exec()
}

type TDGroup struct {
	Clients []*sql.DB
}

type listData struct {
	Name     string
	ListType string
}

func (group *TDGroup) Create(table string, data ...listData) error {
	buffer := bytes.Buffer{}
	buffer.WriteString("create stable ")
	buffer.WriteString(table)
	buffer.WriteString(" (ts timestamp")
	for _, v := range data {
		buffer.WriteString(", ")
		buffer.WriteString(v.Name)
		buffer.WriteString(" ")
		buffer.WriteString(v.ListType)
	}
	buffer.WriteString(" tags (name nchar(20))")

	for _, db := range group.Clients {
		_, err := db.Exec(buffer.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (group *TDGroup) Query(sql string, scan func(rows *sql.Rows)) error {
	for _, db := range group.Clients {
		rows, err := db.Query(sql)
		if err != nil {
			return err
		}
		for rows.Next() {
			scan(rows)
		}
		rows.Close()
	}
	return nil
}

func (group *TDGroup) Exec(sql string) error {
	for _, db := range group.Clients {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}
