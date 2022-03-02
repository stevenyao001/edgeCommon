package redis

import (
	"errors"
	"github.com/go-redis/redis"
)

var redisPool = make(map[string]*redis.Client)

func GetRds(insName string) (*redis.Client, error) {
	rds, ok := redisPool[insName]
	if !ok {
		return rds, errors.New("ins name not found")
	}
	return rds, nil
}

func setRds(insName string, client *redis.Client) {
	redisPool[insName] = client
}
