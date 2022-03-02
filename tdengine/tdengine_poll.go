package tdengine

import "database/sql"

var tdEnginePool map[string]*sql.DB

func init() {
	tdEnginePool = make(map[string]*sql.DB)
}

type TdEngine struct {
	Db        *sql.DB
	InsName   string
	DbName    string
	TableName string
}

func (td *TdEngine) Conn() {
	if td.InsName == "" {
		panic("conn td engine err : ins name don't is empty")
	}

	if td.DbName == "" {
		panic("conn td engine err : db name don't is empty")
	}

	if td.TableName == "" {
		panic("conn td engine err : table name don't is empty")
	}

	if conn, ok := tdEnginePool[td.InsName]; ok {
		td.Db = conn
	} else {
		panic("conn td engine err : ins name not found")
	}
}
