package cache

import (
	"encoding/json"
	"errors"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/op/go-logging"
	configs "gitlab.globoi.com/globoid/glive/env_configs"
)

var (
	// ErrNotFound Throwed when not found on cache
	ErrNotFound = errors.New("Object not found")
	// ErrCache Throwed when a general error occurs
	ErrCache = errors.New("Trouble getting object from cache")
	// ErrUnmarshal Throwed when a Unmarshal error occurs
	ErrUnmarshal = errors.New("Trouble unmarshalling object from cache")
)

// Cache interface to access cached data
type Cache interface {
	Get(key string, i interface{}) error
	Set(key string, i interface{}, timeInSeconds time.Duration) error
}

// RedisCache struct for Redis Memory Cache ONLY
type RedisCache struct {
	conn redigo.Conn
	log  *logging.Logger
}

// NewRedisCache return a RedisCache pointer
func NewRedisCache(conn redigo.Conn, log *logging.Logger) *RedisCache {
	return &RedisCache{conn: conn, log: log}
}

// Get Retrieves an interface by desired key
func (r *RedisCache) Get(key string, i interface{}) error {
	jsonString, err := redigo.String(r.conn.Do("GET", key))
	if err != nil {
		r.log.Error("[RedisCache] Fail when trying to get key %v from cache", key)
		return ErrCache
	}
	if len(jsonString) > 0 {
		if err := json.Unmarshal([]byte(jsonString), i); err != nil {
			r.log.Error("[RedisCache] Fail to unmarshal json string %v from cache", jsonString)
			return ErrUnmarshal
		}
		r.log.Debug("[RedisCache] Key %v was found with value %v", key, jsonString)
		return nil
	}
	r.log.Error("[RedisCache] Key %v was not found", key)
	return ErrNotFound
}

// Set store an object
func (r *RedisCache) Set(key string, i interface{}, duration time.Duration) error {
	serviceJSON, err := json.Marshal(i)
	if err != nil {
		r.log.Error("[RedisCache] Fail to marshal object")
		return ErrCache
	}
	_, err = r.conn.Do("SETEX", key, duration.Seconds(), serviceJSON)
	if err != nil {
		r.log.Error("[RedisCache] Fail when trying to set key %v", key)
		return ErrCache
	}
	return nil
}

// NewCache return a new instance of the redis pool
func NewCache(cfg *configs.Config) *redigo.Pool {
	logger := logging.MustGetLogger("logger")
	pool := &redigo.Pool{
		MaxIdle:     int(cfg.Cache.MaxIdle),
		MaxActive:   int(cfg.Cache.MaxActive),
		IdleTimeout: time.Second * time.Duration(cfg.Cache.IdleTimeoutSecs),
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", cfg.Cache.Address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", cfg.Cache.Auth); err != nil {
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
		logger.Fatal("Cannot connect to Cache: ", err)
	}
	return pool
}
