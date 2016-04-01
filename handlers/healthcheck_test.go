package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	redis "github.com/viniciusfeitosa/QConSP2016/cache/mock"
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

func TestHealthCheckHandlerOK(t *testing.T) {
	setup(t)
	redisMock.EXPECT().Close()
	redisMock.EXPECT().Do("PING").Return(nil, nil).Times(1)
	poolRedisMock.EXPECT().Get().Return(redisMock).Times(1)
	expected := "WORKING"
	handler := NewHealthCheck(poolRedisMock, confMock, logMock)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://blabla/healthcheck", bytes.NewBufferString(""))
	assert.Nil(t, err)
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, expected, recorder.Body.String())
}
