package handlers

import (
	"net/http"

	"github.com/op/go-logging"
	"github.com/viniciusfeitosa/QConSP2016/cache"
)

// UserHandler is a struct
type UserHandler struct {
	cachePool cache.Pool
	log       *logging.Logger
}

// NewUserHandler is a func
func NewUserHandler(cachePool cache.Pool, log *logging.Logger) *UserHandler {
	return &UserHandler{cachePool: cachePool, log: log}
}

// Create is a handler func
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

}
