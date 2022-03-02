package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisConf struct {
	InsName      string `json:"ins_name"`
	Addr         string `json:"addr"`
	Auth         string `json:"auth"`
	Db           int    `json:"db"`
	ConnTimeout  int    `json:"conn_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	MaxActive    int    `json:"max_active"`
	MinIdleConns int    `json:"min_idle_conns"`
}

func InitRedis(redisConf []RedisConf) {

	for _, conf := range redisConf {

		if conf.InsName == "" {
			panic("init redis err : ins name don't is empty")
		}

		if conf.Addr == "" {
			panic("init redis err : addr don't is empty")
		}

		if _, err := GetRds(conf.InsName); err == nil {
			continue
		}

		options := &redis.Options{
			Network:         "tcp",
			Addr:            conf.Addr,
			Dialer:          nil,
			OnConnect:       nil,
			Password:        conf.Auth,
			DB:              conf.Db,
			MaxRetries:      0,
			MinRetryBackoff: 0,
			MaxRetryBackoff: 0,
			DialTimeout:     time.Millisecond * time.Duration(conf.ConnTimeout),
			ReadTimeout:     time.Millisecond * time.Duration(conf.ReadTimeout),
			WriteTimeout:    time.Millisecond * time.Duration(conf.WriteTimeout),
			PoolSize:        conf.MaxActive,
			MinIdleConns:    conf.MinIdleConns,
			//MaxConnAge:         time.Second * time.Duration(conf.MaxConnAge),
			PoolTimeout: 0,
			//IdleTimeout:        time.Second * time.Duration(conf.IdleTimeout),
			IdleCheckFrequency: 0,
			TLSConfig:          nil,
		}

		redisClient := redis.NewClient(options)

		_, err := redisClient.Ping().Result()
		if err != nil {
			panic("init redis [" + conf.InsName + "] err : ping fail : " + err.Error())
		}

		setRds(conf.InsName, redisClient)
	}
}
