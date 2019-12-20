package gmq

import (
	"github.com/go-redis/redis"
)

type RedisPool interface {
	SetDb(name string, db *redis.Client)
	GetDb(name string) (*redis.Client, error)
	NewRedis(host, port, pwd string, database int) (*redis.Client, error)
}
