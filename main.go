package main

import (
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/bmizerany/pat"
	"github.com/codegangsta/negroni"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	"github.com/viniciusfeitosa/QConSP2016/env_configs"
	"github.com/viniciusfeitosa/QConSP2016/handlers"
	"github.com/viniciusfeitosa/QConSP2016/mongo"
	"github.com/viniciusfeitosa/QConSP2016/workers"
)

func main() {
	conf := configs.LoadConfig()
	log := configs.GetLogger(conf)
	cachePool := cache.NewPool(conf)
	db := mongo.StartDB(conf, log)

	go workers.UsersToMongo(log, conf, cachePool, db)

	router := pat.New()

	router.Get("/healthcheck", handlers.NewHealthCheck(cachePool, conf, log))

	userHandler := handlers.NewUserHandler(cachePool, log)
	router.Post("/user/create", http.HandlerFunc(userHandler.Create))
	router.Get("/user/all", http.HandlerFunc(userHandler.FindAll))
	router.Get("/user/:id", http.HandlerFunc(userHandler.FindByID))

	n := negroni.Classic()
	n.UseHandler(router)
	log.Infof("App started on port: %s", conf.Address)
	graceful.Run(conf.Address, 30*time.Second, n)
}
