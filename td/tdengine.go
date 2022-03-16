package td

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v2/taosRestful"
)

type Engine struct {
	TDGroup
}

func CollectorTD(conf CollectorConf) *Engine {
	conf = Default(conf)
	url := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", conf.Username, conf.Password, conf.Network, conf.Addr, conf.Port, conf.Db)
	db, err := sql.Open(conf.Driver, url)
	if err != nil {
		panic("init td engine err : " + err.Error())
	}

	db.SetMaxIdleConns(conf.MaxIdleConns)
	//db.SetConnMaxIdleTime(time.Duration(pgsqlConfig.MaxIdleTime) * time.Second)
	//db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	engine := &Engine{TDGroup{DB: db}}
	return engine
}
