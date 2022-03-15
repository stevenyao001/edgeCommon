package td

import (
	"database/sql"
	"fmt"
)

type Conf struct {
	InsName      string `json:"ins_name"`
	Driver       string `json:"driver"`
	Network      string `json:"network"`
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

type Engine struct {
	TDGroup
}

func InitTd(tdConf []Conf) *Engine {
	tdEnginePool := make(map[string]*sql.DB)

	for _, conf := range tdConf {

		url := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", conf.Username, conf.Password, conf.Network, conf.Addr, conf.Port, conf.Db)
		db, err := sql.Open(conf.Driver, url)
		if err != nil {
			panic("init td engine err : " + err.Error())
		}

		db.SetMaxIdleConns(conf.MaxIdleConns)
		//db.SetConnMaxIdleTime(time.Duration(pgsqlConfig.MaxIdleTime) * time.Second)
		//db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)
		db.SetMaxOpenConns(conf.MaxOpenConns)

		tdEnginePool[conf.InsName] = db
	}
	engin := &Engine{TDGroup{Clients: tdEnginePool}}
	return engin
}
