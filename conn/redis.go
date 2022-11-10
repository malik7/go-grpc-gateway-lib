package conn

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

/*
RedisConfig holds per database config
Host				: address (usualy IP) of redis. Example: 127.0.0.1
Port				: port of redis. Example: 6379
MaxActive			: maximum number of connections allocated by the pool. When zero, there is no limit on the number of connections in the pool. Default value: 0
Wait				: if wait is true and the pool is at the MaxActive limit, then Get() waits for the connection to be returned to the pool. Default value: false
Database			: database index has same function as prefix. It encapsulate the data for each index
*/
type RedisConfig struct {
	Host      string `envconfig:"HOST"`
	Port      int    `envconfig:"PORT"`
	MaxIdle   int    `envconfig:"MAX_IDLE"`
	MaxActive int    `envconfig:"MAX_ACTIVE"`
	Wait      bool   `envconfig:"WAIT"`
	Database  int    `envconfig:"DATABASE"`
}

// RedisConn holds all redis connections which exist in this repository
type RedisConn struct {
	Core *redis.Pool
}

// redisDial is used to mock redis.Dial fsunction
var redisDial = redis.Dial

// InitRedis init the redis from config to redis connection
func InitRedis(cfg RedisConfig) (*redis.Pool, error) {
	if cfg.MaxIdle == 0 {
		cfg.MaxIdle = cfg.MaxActive
	}

	var (
		address = cfg.Host + ":" + strconv.Itoa(cfg.Port)
		pool    = &redis.Pool{
			IdleTimeout:     2 * time.Second,
			MaxConnLifetime: 10 * time.Second,
			MaxActive:       cfg.MaxActive,
			MaxIdle:         cfg.MaxIdle,
			Wait:            cfg.Wait,
			Dial: func() (redis.Conn, error) {
				return redisDial("tcp", address, redis.DialDatabase(cfg.Database))
			},
		}
	)

	err := examineRedisConn(pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func examineRedisConn(pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}
