package cache_test

import (
	"fmt"
	"testing"
	"time"

	"gitlab.globoi.com/globoid/glive/cache"
	"gitlab.globoi.com/globoid/glive/env_configs"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockedConn  *connMock
	logMock     *logging.Logger
	callCache   *cache.RedisCache
	myKey       string
	myInterface *person
)

type connMock struct {
	mock.Mock
}

func (c *connMock) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	params := c.Called(commandName, args)
	return params.Get(0), params.Error(1)
}

func (c *connMock) Close() error                                       { return nil }
func (c *connMock) Err() error                                         { return nil }
func (c *connMock) Send(commandName string, args ...interface{}) error { return nil }
func (c *connMock) Flush() error                                       { return nil }
func (c *connMock) Receive() (reply interface{}, err error)            { return nil, nil }

type person struct {
	Name string
}

func setup() {
	configMock := &configs.Config{Debug: true}
	logMock = configs.GetLogger(configMock)
	mockedConn = new(connMock)
	callCache = cache.NewRedisCache(mockedConn, logMock)
	myKey = "myKey:1"
	myInterface = new(person)
}

func TestGet_ErrorWhenTryToRetrieveAnElementeOnCache(t *testing.T) {
	setup()
	mockedConn.On("Do", "GET", []interface{}{myKey}).Return(nil, fmt.Errorf("Fail to get from redis"))
	err := callCache.Get(myKey, myInterface)
	assert.NotNil(t, err)
	assert.Equal(t, cache.ErrCache, err)
}

func TestGet_ErrorWhenTryUnmarshalAnElementRetrievedFromCache(t *testing.T) {
	setup()
	mockedConn.On("Do", "GET", []interface{}{myKey}).Return(`{ a: }`, nil)
	err := callCache.Get(myKey, myInterface)
	assert.NotNil(t, err)
	assert.Equal(t, cache.ErrUnmarshal, err)
}

func TestGet_ErrorElementNotFound(t *testing.T) {
	setup()
	mockedConn.On("Do", "GET", []interface{}{myKey}).Return("", nil)
	err := callCache.Get(myKey, myInterface)
	assert.NotNil(t, err)
	assert.Equal(t, cache.ErrNotFound, err)
}

func TestGet_WhenElementFoundOnRedis(t *testing.T) {
	setup()
	mockedConn.On("Do", "GET", []interface{}{myKey}).Return(`{ "name": "Fulano" }`, nil)
	err := callCache.Get(myKey, myInterface)
	assert.Nil(t, err)
	assert.Equal(t, myInterface.Name, "Fulano")
}

func Test_RedisCache_Set(t *testing.T) {
	cfg := &configs.Config{
		Cache: configs.CacheConfig{
			Address:         "Address",
			Auth:            "auth",
			DB:              0,
			MaxIdle:         uint(300),
			MaxActive:       uint(300),
			IdleTimeoutSecs: uint(300),
		},
	}
	cachePool := cache.NewCache(cfg)
	rc := cache.NewRedisCache(cachePool.Get(), configs.GetLogger(cfg))
	type testObj struct {
		ID int
	}
	objToCache := testObj{ID: 123}
	err := rc.Set("servico-1", objToCache, 10*time.Second)
	assert.NoError(t, err)
	obj := new(testObj)
	err = rc.Get("servico-1", obj)
	assert.NoError(t, err)
	assert.EqualValues(t, &objToCache, obj)
}
