package bManager

import (
	"dancechanlibrary/broker"
	"dancechanlibrary/broker/bRedis"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"net/http"
	_ "net/http"
	"runtime"
	"testing"
	"time"
)

func brokerMock(loggers []broker.BLogger) broker.Broker {
	m := NewBManager(loggers)
	return m
}

func TestManager_Register(t *testing.T) {
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	defer m.Stop()

	//when
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 1, 10, make([]broker.BLogger, 0))

	//then
	assert.IsType(t, &bManager{}, m)
}

func TestManager_Up(t *testing.T) {
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 1, 10, make([]broker.BLogger, 0))

	//when
	c := runtime.NumGoroutine()
	m.Up() // 고루틴으로 동작

	//then
	assert.Equal(t, c+1, runtime.NumGoroutine())
}

func TestManager_PushPopWork(t *testing.T) {
	//given
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 1, 10, make([]broker.BLogger, 0))
	m.Up()

	//when
	pushCommand := bRedis.RPush{
		Key:        "key",
		Value:      "value",
		ResultChan: make(chan broker.CmdResult),
	}
	popCommand := bRedis.BLPop{
		Key:        "key",
		ResultChan: make(chan broker.CmdResult),
	}
	_ = m.Worker(broker.BROKER_WORKER_REDIS).PushWork(&pushCommand)
	_ = m.Worker(broker.BROKER_WORKER_REDIS).PushWork(&popCommand)

	//then
	result := pushCommand.Result()
	assert.NotEqual(t, int64(0), result.Result)
	assert.Nil(t, result.Err)

	result = popCommand.Result()
	assert.Equal(t, "value", result.Result)
	assert.Nil(t, result.Err)
}

func TestBManager_Stop(t *testing.T) {
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 1, 10, make([]broker.BLogger, 0))

	c := runtime.NumGoroutine()
	m.Up() // 고루틴으로 동작

	//when
	m.Stop()

	//then
	assert.Equal(t, c, runtime.NumGoroutine())
}

type sampleLog struct {
}

func (receiver sampleLog) WriteBLog(message string) {
	// BLogger 구현부
	fmt.Println(message)
}

func TestBManager_Logger(t *testing.T) {
	//given    []broker.BLogger{sampleLog{}
	m := brokerMock(
		[]broker.BLogger{sampleLog{}},
	)
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 1, 10, []broker.BLogger{sampleLog{}})

	m.Up()

	//when
	command := bRedis.RPush{
		Key:   "key",
		Value: "value",
	}
	m.Stop() // 에러 유도
	err := m.Worker(broker.BROKER_WORKER_REDIS).PushWork(&command)

	//then
	assert.NotNil(t, err)
}

func TestBManager_RegisterBHttp(t *testing.T) {
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	defer m.Stop()
	_ = m.Register(&http.Client{
		Timeout: time.Second,
	}, 1, 10, make([]broker.BLogger, 0))
	c := runtime.NumGoroutine()

	//when
	m.Up()

	//then
	assert.Equal(t, c+1, runtime.NumGoroutine())
}

func TestBManager_MultiRegister_Up_동작_확인(t *testing.T) {
	//given
	m := brokerMock(make([]broker.BLogger, 0))
	defer m.Stop()

	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 2, 10, make([]broker.BLogger, 0))
	_ = m.Register(&http.Client{
		Timeout: time.Second,
	}, 2, 10, make([]broker.BLogger, 0))

	c := runtime.NumGoroutine()

	//when
	m.Up()

	//then
	assert.Equal(t, c+4, runtime.NumGoroutine())
}

func TestBManager_MultiRegister_Stop_동작_확인(t *testing.T) {
	//given
	m := brokerMock([]broker.BLogger{sampleLog{}})

	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}), 2, 10, []broker.BLogger{sampleLog{}})
	_ = m.Register(&http.Client{
		Timeout: time.Second,
	}, 2, 10, []broker.BLogger{sampleLog{}})

	c := runtime.NumGoroutine()
	m.Up()

	//when
	m.Stop()

	//then
	assert.Equal(t, c, runtime.NumGoroutine())
}
