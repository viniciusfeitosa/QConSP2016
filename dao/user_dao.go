package dao

import (
	"github.com/garyburd/redigo/redis"
	"github.com/op/go-logging"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	"github.com/viniciusfeitosa/QConSP2016/models"
	"github.com/viniciusfeitosa/QConSP2016/mongo"
)

const (
	usersQueue      = "USERS"
	collectionUsers = "users"
)

// UserRedisDAO is the struct with fields of the dao
type UserRedisDAO struct {
	cachePool cache.Pool
	log       *logging.Logger
}

// NewUserRedisDAO generate an instance of userRedisDAO
func NewUserRedisDAO(cachePool cache.Pool, log *logging.Logger) *UserRedisDAO {
	return &UserRedisDAO{cachePool: cachePool, log: log}
}

// Create is a method to create user
func (dao *UserRedisDAO) Create(user *models.UserModel) error {
	conn := dao.cachePool.Get()
	defer conn.Close()

	userModel, err := user.GetUserModel()
	if err != nil {
		dao.log.Errorf("[UserRedisDAO Create] error to get userModel from user %s, error: %s", user.ID.Hex(), err.Error())
		return err
	}

	userJSON, err := userModel.ToJSONString()
	if err != nil {
		dao.log.Errorf("[UserRedisDAO Create] error to get a json from user %s, error: %s", user.ID.Hex(), err.Error())
		return err
	}

	if _, err := conn.Do("SET", user.ID.Hex(), userJSON); err != nil {
		dao.log.Errorf("[UserRedisDAO Create] error to set JSON data for redis with user %s, error: %s", user.ID.Hex(), err.Error())
		return err
	}
	if _, err := conn.Do("RPUSH", usersQueue, user.ID.Hex()); err != nil {
		dao.log.Errorf("[UserRedisDAO Create] error to set for redis queue with user %s, error: %s", user.ID.Hex(), err.Error())
		return err
	}
	return nil
}

// FindByID return an user from cache
func (dao *UserRedisDAO) FindByID(id string) (string, error) {
	conn := dao.cachePool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", id))
	if err != nil {
		dao.log.Errorf("[UserRedisDAO FindByID] %s", err.Error())
		return "", err
	}

	return value, nil
}

// FindAll return all users from cache
func (dao *UserRedisDAO) FindAll() ([]models.UserModel, error) {
	conn := dao.cachePool.Get()
	defer conn.Close()
	var userList []models.UserModel
	userModel := &models.UserModel{}

	values, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		dao.log.Errorf("[UserRedisDAO FindAll] %s", err.Error())
		return userList, err
	}

	for _, v := range values {

		user, err := dao.FindByID(v)
		if err != nil {
			dao.log.Errorf("[UserRedisDAO FindAll] %s", err.Error())
			return userList, err
		}

		model, err := userModel.FromJSONToObject(user)
		if err != nil {
			dao.log.Errorf("[UserRedisDAO FindAll] %s", err.Error())
			return userList, err
		}
		userList = append(userList, model)
	}
	return userList, nil
}

// UserMongoDAO is the struct with fields of the dao
type UserMongoDAO struct {
	db  *mongo.DB
	log *logging.Logger
}

// NewUserMongoDAO generate an instance of NewUserMongoDAO
func NewUserMongoDAO(db *mongo.DB, log *logging.Logger) *UserMongoDAO {
	return &UserMongoDAO{db: db, log: log}
}

// Create is the method to create users on mongo
func (u *UserMongoDAO) Create(user *models.UserModel) error {
	if err := u.db.Insert(collectionUsers, user); err != nil {
		u.log.Errorf("[UserMongoDAO Create] error to insert user in mongo %s error: %s", user.ID.Hex(), err.Error())
		return err
	}
	u.log.Infof("[UserMongoDAO Create] User created in mongo %s", user.ID.Hex())

	return nil
}
