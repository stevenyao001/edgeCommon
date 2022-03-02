package edgeCommon

import (
	"github.com/gd1024/edge_common/config"
	"github.com/gd1024/edge_common/http"
	logger2 "github.com/gd1024/edge_common/logger"
	"github.com/gd1024/edge_common/mqtt"
	"github.com/gd1024/edge_common/pgsql"
	"github.com/gd1024/edge_common/tdengine"
	"github.com/gin-gonic/gin"
)

type engine struct {
}

func New() *engine {
	return &engine{}
}

func (e *engine) RegisterConfig(filePath string, conf interface{}) {
	config.InitConf(filePath, conf)
}

func (e *engine) RegisterLogger(logPath string) {
	logger2.InitLog(logPath)
}

func (e *engine) RegisterMqtt(confs []mqtt.Conf, subOpts map[string][]mqtt.SubscribeOpts) {
	mqtt.InitMqtt(confs, subOpts)
}

func (e *engine) RegisterPgsql(pgConf []pgsql.Conf) {
	pgsql.InitPgsql(pgConf)
}

func (e *engine) RegisterTdEngine(tdConf []tdengine.Conf) {
	tdengine.InitTdEngine(tdConf)
}

func (e *engine) RegisterHttp(httpConf http.Conf, middleware ...gin.HandlerFunc) {
	http.InitHttp(httpConf, middleware...)
}
