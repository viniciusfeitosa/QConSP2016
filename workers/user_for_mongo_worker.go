package workers

import (
	"sync"

	"github.com/viniciusfeitosa/QConSP2016/dao"
	"github.com/viniciusfeitosa/QConSP2016/models"

	"github.com/viniciusfeitosa/QConSP2016/mongo"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/op/go-logging"
	"github.com/viniciusfeitosa/QConSP2016/cache"
	"github.com/viniciusfeitosa/QConSP2016/env_configs"
)

const (
	usersToMongoQueue = "USERS"
)

// UsersMongo is the struct with parameters
// to events in mongo
type UsersMongo struct {
	id        int
	cachePool cache.Pool
	config    *configs.Config
	log       *logging.Logger
	db        *mongo.DB
}

// UsersToMongo is the func to start workers to log event in GLive
func UsersToMongo(log *logging.Logger, config *configs.Config, cachePool cache.Pool, db *mongo.DB) {
	var wg sync.WaitGroup
	for i := 0; i < config.NumWorkers.NumUsersMongoWorkers; i++ {
		wg.Add(1)
		go func(cachePool cache.Pool, config *configs.Config, db *mongo.DB, log *logging.Logger, id int) {
			UsersMongo := newUsersMongo(cachePool, config, db, log, i)
			UsersMongo.process(i)
			defer wg.Done()
		}(cachePool, config, db, log, i)
	}
	wg.Wait()
}

// newUsersMongo generate a instance of UsersMongo
func newUsersMongo(cachePool cache.Pool, config *configs.Config, db *mongo.DB, log *logging.Logger, id int) *UsersMongo {
	return &UsersMongo{cachePool: cachePool, config: config, db: db, log: log, id: id}
}

func (l *UsersMongo) process(id int) {
	for {
		conn := l.cachePool.Get()
		var channel, uuid string
		if reply, err := redigo.Values(conn.Do("BLPOP", usersToMongoQueue, 30+id)); err == nil {

			if _, err := redigo.Scan(reply, &channel, &uuid); err != nil {
				l.log.Errorf("[UsersMongo - process] error in USERS pulling: %s", err.Error())
				l.requeue(uuid)
				continue
			}

			values, err := redigo.String(conn.Do("GET", uuid))
			if err != nil {
				l.log.Errorf("[UsersMongo - process] error in USERS pulling: %s", err.Error())
				l.requeue(uuid)
				continue
			}

			user := new(models.UserModel)
			userObject, err := user.FromJSONToObject(values)
			if err != nil {
				l.log.Errorf("[UsersMongo - process] error to unmarshal USERS pulling: %s", err.Error())
				l.requeue(uuid)
				continue
			}

			dao := dao.NewUserMongoDAO(l.db, l.log)
			if err := dao.Create(&userObject); err != nil {
				l.log.Errorf("[UsersMongo - process] error to insert USERS in mongo: %s", err.Error())
				l.requeue(uuid)
				continue
			}

		} else if err != redigo.ErrNil {
			l.log.Errorf("[UsersMongo - process] error runing BLPOP on USERS: %s", err.Error())
		}
		conn.Close()
	}
}

func (l *UsersMongo) requeue(uuid string) {
	for {
		conn := l.cachePool.Get()
		defer conn.Close()
		_, err := conn.Do("RPUSH", usersToMongoQueue, uuid)
		if err == nil {
			return
		}
		l.log.Errorf("[UsersMongo - requeue] error for requeue to %s with user %s error: %s", usersToMongoQueue, uuid, err.Error())
	}
}
