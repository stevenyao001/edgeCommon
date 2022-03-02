package tdengine

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v2/taosRestful"
	//_ "github.com/taosdata/driver-go/v2/taosSql"
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

func InitTdEngine(tdConf []Conf) {
	//   sddfgdfgsddsdfsdsdfsdfs
	//url := "root:taosdata@/tcp(localhost:6030)/test2"
	//db, err := sql.Open("taosSql", url)
	//var taosDSN = "root:taosdata@http(localhost:6041)/test2"
	//db, err := sql.Open("taosRestful", taosDSN)

	for _, conf := range tdConf {

		url := fmt.Sprintf("%s:%s@/%s(%s:%d)/%s", conf.Username, conf.Password, conf.Network, conf.Addr, conf.Port, conf.Db)
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
}
