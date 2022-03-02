package edgeCommon

import (
	"github.com/gin-gonic/gin"
	"github.com/stevenyao001/edgeCommon/config"
	"github.com/stevenyao001/edgeCommon/http"
	logger2 "github.com/stevenyao001/edgeCommon/logger"
	"github.com/stevenyao001/edgeCommon/mqtt"
	"github.com/stevenyao001/edgeCommon/pgsql"
	"github.com/stevenyao001/edgeCommon/tdengine"
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
