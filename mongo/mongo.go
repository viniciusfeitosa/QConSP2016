package mongo

import (
	"time"

	logging "github.com/op/go-logging"
	configs "github.com/viniciusfeitosa/QConSP2016/env_configs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo is the interface with methods to return a mongo session
type Mongo interface {
	Insert(collection string, obj interface{}) error
	Remove(collection string, query bson.M) error
	FindOne(collection string, query bson.M, obj interface{}) error
	ExistsInDB(collection string, query bson.M) (bool, error)
	Update(collection string, key, change bson.M) error
}

// DB is the Struct that implements MongoInterface
type DB struct {
	conf    *configs.Config
	log     *logging.Logger
	session *mgo.Session
}

// StartDB return a session of mongo database
func StartDB(conf *configs.Config, log *logging.Logger) *DB {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs: []string{
			conf.Mongo.Address1,
		},
		Database:  conf.Mongo.DatabaseName,
		Timeout:   time.Second * 60,
		PoolLimit: 100,
	}

	s, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	s.SetSafe(&mgo.Safe{WMode: "majority"})

	if err := s.Ping(); err != nil {
		log.Error(err.Error())
		return nil
	}

	return &DB{conf: conf, log: log, session: s}
}

// FindOne return one instance of the find
func (db *DB) FindOne(collection string, query bson.M, obj interface{}) error {
	session := db.session.Copy()
	defer session.Close()
	if err := session.DB(db.conf.Mongo.DatabaseName).C(collection).Find(query).One(obj); err != nil {
		return err
	}
	return nil
}

// Insert a new object in collection
func (db *DB) Insert(collection string, obj interface{}) error {
	session := db.session.Copy()
	defer session.Close()
	if err := session.DB(db.conf.Mongo.DatabaseName).C(collection).Insert(obj); err != nil {
		return err
	}
	return nil
}

// Update change an exists object in my collection
func (db *DB) Update(collection string, key, change bson.M) error {
	session := db.session.Copy()
	defer session.Close()
	if err := session.DB(db.conf.Mongo.DatabaseName).C(collection).Update(key, change); err != nil {
		return err
	}
	return nil
}

// Upsert change an exists object in my collection
func (db *DB) Upsert(collection string, key, change bson.M) error {
	session := db.session.Copy()
	defer session.Close()

	if _, err := session.DB(db.conf.Mongo.DatabaseName).C(collection).Upsert(key, change); err != nil {
		return err
	}
	return nil
}

// Remove an object in collection
func (db *DB) Remove(collection string, query bson.M) error {
	session := db.session.Copy()
	defer session.Close()
	if _, err := session.DB(db.conf.Mongo.DatabaseName).C(collection).RemoveAll(query); err != nil {
		return err
	}
	return nil
}

// Close is a method to close session of the database
func (db *DB) Close() {
	db.session.Close()
}
