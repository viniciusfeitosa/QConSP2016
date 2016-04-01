package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/op/go-logging"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/viniciusfeitosa/QConSP2016/env_configs"
)

type mockRedis struct {
	testifyMock.Mock
}

var (
	mock          *gomock.Controller
	poolRedisMock *redis.MockPool
	redisMock     *redis.MockConn
	confMock      *configs.Config
	logMock       *logging.Logger
)

func setup(t *testing.T) {
	mock = gomock.NewController(t)
	poolRedisMock = redis.NewMockPool(mock)
	redisMock = redis.NewMockConn(mock)
	confMock = &configs.Config{Debug: true}
	logMock = configs.GetLogger(confMock)
}