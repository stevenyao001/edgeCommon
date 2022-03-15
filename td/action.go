package td

import (
	"database/sql"
)

type DataAction interface {
	Query()
	Scan()
	Exec()
}

type TDGroup struct {
	Clients []*sql.DB
	Rows    []*sql.Rows
}

func (group *TDGroup) Query(sql string) error {
	for _, db := range group.Clients {
		rows, err := db.Query(sql)
		if err != nil {
			return err
		}
		group.Rows = append(group.Rows, rows)
	}
	return nil
}

func (group *TDGroup) Scan(dest ...interface{}) error {
	for _, rows := range group.Rows {
		defer rows.Close()
		for rows.Next() {
			return rows.Scan(dest)
		}
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
