package main

import (
	"sync"

	"github.com/op/go-logging"
	configs "gitlab.globoi.com/globoid/glive_consumer/env_configs"
	"gitlab.globoi.com/globoid/glive_consumer/redis"
)

func main() {
	conf := configs.LoadConfig()
	log := configs.GetLogger(conf)
	redisPool := redis.NewPool(conf)

	// omitted code
	go runLogEvents(log, conf, redisPool)
	// omitted code
}

func runLogEvents(log *logging.Logger, config *configs.Config, redisPool redis.Pool) {
	var wg sync.WaitGroup
	for i := 0; i < config.NumWorkers.NumLogEventsWorkers; i++ {
		wg.Add(1)
		go func(redisPool redis.Pool, config *configs.Config, log *logging.Logger, id int) {
			logEvents := newLogEvents(redisPool, config, log, i)
			logEvents.process(i)
			defer wg.Done()
		}(redisPool, config, log, i)
	}
	wg.Wait()
}

func newLogEvents(redisPool redis.Pool, config *configs.Config, log *logging.Logger, id int) *LogEvents {
	return &LogEvents{redisPool: redisPool, config: config, log: log, id: id}
}

type LogEvents struct {
	id        int
	redisPool redis.Pool
	config    *configs.Config
	log       *logging.Logger
}

func (l *LogEvents) process(id int) {
	for {
		// omitted code
	}
}
