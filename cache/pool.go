package cache

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/op/go-logging"
	configs "github.com/viniciusfeitosa/QConSP2016/env_configs"
)

// Pool is the interface to pool of redis
type Pool interface {
	Get() redigo.Conn
}

// NewPool return a new instance of the redis pool
func NewPool(cfg *configs.Config) *redigo.Pool {
	logger := logging.MustGetLogger("logger")
	pool := &redigo.Pool{
		MaxIdle:     int(cfg.Redis.MaxIdle),
		MaxActive:   int(cfg.Redis.MaxActive),
		IdleTimeout: time.Second * time.Duration(cfg.Redis.IdleTimeoutSecs),
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", cfg.Redis.Address)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", cfg.Redis.Auth); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			if _, err := c.Do("SELECT", cfg.Redis.DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	c := pool.Get() // Test connection during init
	if _, err := c.Do("PING"); err != nil {
		logger.Fatal("Cannot connect to Redis: ", err)
	}

	return pool
}
