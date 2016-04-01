package handlers

import (
	"fmt"
	"net/http"

	"github.com/op/go-logging"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	config "github.com/viniciusfeitosa/QConSP2016/env_configs"
)

// HealthCheck is a struct
type HealthCheck struct {
	cachePool cache.Pool
	conf      *config.Config
	log       *logging.Logger
}

// NewHealthCheck is a func
func NewHealthCheck(cachePool cache.Pool, conf *config.Config, log *logging.Logger) *HealthCheck {
	return &HealthCheck{
		cachePool: cachePool,
		conf:      conf,
		log:       log,
	}
}

func (hc *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := hc.cachePool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		fmt.Fprint(w, "FAIL\n")
	} else {
		fmt.Fprintf(w, "WORKING")
	}
}
