package td

import (
	"database/sql"
)

type DataAction interface {
	SubData()
	PubData()
}

type TDGroup struct {
	Clients map[string]*sql.DB
}

func (group *TDGroup) Find(insName string, sql string, dest ...interface{}) error {
	db := group.Clients[insName]
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(dest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (group *TDGroup) Insert(insName string, sql string) error {
	db := group.Clients[insName]
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
