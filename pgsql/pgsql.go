package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)



type Conf struct {
	InsName      string `json:"ins_name"`
	Addr         string `json:"addr"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Db           string `json:"db"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxIdleTime  int    `json:"max_idle_time"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxOpenConns int    `json:"max_open_conns"`
}

func InitPgsql(pgConf []Conf) {

	for _, pgsqlConfig := range pgConf {

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			pgsqlConfig.Addr, pgsqlConfig.Port, pgsqlConfig.Username, pgsqlConfig.Password, pgsqlConfig.Db)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic("init postgres open err :" + pgsqlConfig.InsName + ":" + err.Error())
		}

		db.SetMaxIdleConns(pgsqlConfig.MaxIdleConns)
		//db.SetConnMaxIdleTime(time.Duration(pgsqlConfig.MaxIdleTime) * time.Second)
		//db.SetConnMaxLifetime(time.Duration(pgsqlConfig.MaxLifeTime) * time.Second)
		db.SetMaxOpenConns(pgsqlConfig.MaxOpenConns)

		err = db.Ping()
		if err != nil {
			panic("init postgres ping err :" + pgsqlConfig.InsName + ":" + err.Error())
		}
		pgPool[pgsqlConfig.InsName] = db
	}
}
