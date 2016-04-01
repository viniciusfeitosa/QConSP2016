package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	"github.com/viniciusfeitosa/QConSP2016/env_configs"
	"github.com/viniciusfeitosa/QConSP2016/handlers"
	"github.com/viniciusfeitosa/QConSP2016/mongo"
)

func main() {
	conf := configs.LoadConfig()
	log := configs.GetLogger(conf)
	cachePool := cache.NewPool(conf)
	db := mongo.StartDB(conf, log)

	router := pat.New()

	router.Get("/healthcheck", handlers.NewHealthCheck(cachePool, conf, log))

	userHandler := handlers.NewUserHandler(cachePool, log)
	router.Post("/user/create", http.HandlerFunc(userHandler.Create))
}
