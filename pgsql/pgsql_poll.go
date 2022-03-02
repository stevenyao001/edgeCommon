package pgsql

import "database/sql"

var pgPool map[string]*sql.DB

type Postgres struct {
	Db        *sql.DB
	InsName   string
	DbName    string
	TableName string
}

func (pg *Postgres) Conn() {
	if pg.InsName == "" {
		panic("conn pgsql err : ins name don't is empty")
	}

	if pg.DbName == "" {
		panic("conn pgsql err : db name don't is empty")
	}

	if pg.TableName == "" {
		panic("conn pgsql err : table name don't is empty")
	}

	if conn, ok := pgPool[pg.InsName]; ok {
		pg.Db = conn
	} else {
		panic("conn pgsql err : ins name not found")
	}
}
