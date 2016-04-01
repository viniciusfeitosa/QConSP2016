package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/op/go-logging"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	"github.com/viniciusfeitosa/QConSP2016/models"
	"github.com/viniciusfeitosa/poc/dao"
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
	userModel := &models.UserModel{}
	conn := u.cachePool.Get()
	defer conn.Close()

	if err := json.NewDecoder(r.Body).Decode(userModel); err != nil {
		u.log.Errorf("[UserHandler Create] Error to parse body error: %s", err.Error())
		http.Error(w, "Error to parse body", http.StatusBadRequest)
		return
	}

	dao := dao.NewUserRedisDAO(u.cachePool, u.log)
	if err := dao.Create(userModel); err != nil {
		u.log.Errorf("[UserHandler CreateUser] Error to create new user status: %d error: %s", http.StatusBadRequest, err.Error())
		http.Error(w, "Error to create new user", http.StatusBadRequest)
		return
	}

	output, err := userModel.ToJSONString()
	if err != nil {
		u.log.Errorf("[UserHandler CreateUser] Error to generate output user status: %d error: %s", http.StatusBadRequest, err.Error())
		http.Error(w, "Error to create new user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(output))
}

// FindByID return an user from mongo
func (u *UserHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	conn := u.cachePool.Get()
	defer conn.Close()

	id := r.URL.Query().Get(":id")

	dao := dao.NewUserRedisDAO(u.cachePool, u.log)
	value, err := dao.FindByID(id)
	if err != nil {
		u.log.Errorf("[UserHandler CreateUser] Error to get user: %s error: %s", id, err.Error())
		http.Error(w, "Error to get user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// FindAll return all users from mongo
func (u *UserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	conn := u.cachePool.Get()
	defer conn.Close()

	dao := dao.NewUserRedisDAO(u.cachePool, u.log)
	values, err := dao.FindAll()
	if err != nil {
		u.log.Errorf("[UserHandler CreateUser] Error to get all users error: %s", err.Error())
		http.Error(w, "Error to get all users", http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(values)
	if err != nil {
		u.log.Errorf("[UserHandler CreateUser] Error to generate output error: %s", err.Error())
		http.Error(w, "Error to get all users", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
