package edgeCommon

import (
	"github.com/gin-gonic/gin"
	"github.com/stevenyao001/edgeCommon/config"
	"github.com/stevenyao001/edgeCommon/emqx"
	"github.com/stevenyao001/edgeCommon/http"
	"github.com/stevenyao001/edgeCommon/logger"
	"github.com/stevenyao001/edgeCommon/mqtt"
	"github.com/stevenyao001/edgeCommon/pgsql"
	"github.com/stevenyao001/edgeCommon/redis"
	"github.com/stevenyao001/edgeCommon/tdengine"
)

type EdgeCommon interface {
	RegisterConfig(filePath string, conf interface{})
	RegisterLogger(logPath string)
	RegisterHttp(conf http.Conf, middleware ...gin.HandlerFunc)
	RegisterMqtt(conf []mqtt.Conf, subOpts map[string][]mqtt.SubscribeOpts)
	RegisterEmqx(conf []emqx.Conf) *emqx.Engine
	RegisterPgsql(conf []pgsql.Conf)
	RegisterRedis(conf []redis.Conf)
	RegisterTdEngine(conf []tdengine.Conf)
}

func New() EdgeCommon {
	return new(engine)
}

type engine struct {
}

func (e *engine) RegisterConfig(filePath string, conf interface{}) {
	config.InitConf(filePath, conf)
}

func (e *engine) RegisterLogger(logPath string) {
	logger.InitLog(logPath)
}

func (e *engine) RegisterMqtt(conf []mqtt.Conf, subOpts map[string][]mqtt.SubscribeOpts) {
	mqtt.InitMqtt(conf, subOpts)
}

func (e *engine) RegisterEmqx(conf []emqx.Conf) *emqx.Engine {
	return emqx.InitEmqx(conf)
}

func (e *engine) RegisterPgsql(conf []pgsql.Conf) {
	pgsql.InitPgsql(conf)
}

func (e *engine) RegisterTdEngine(conf []tdengine.Conf) {
	tdengine.InitTdEngine(conf)
}

func (e *engine) RegisterHttp(conf http.Conf, middleware ...gin.HandlerFunc) {
	http.InitHttp(conf, middleware...)
}

func (e *engine) RegisterRedis(conf []redis.Conf) {
	redis.InitRedis(conf)
}
