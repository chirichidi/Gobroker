package bManager

import (
	"dancechanlibrary/broker"
	"dancechanlibrary/broker/bRedis"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func brokerMock(param ManagerParam) broker.Broker {
	m := NewBManager(param)
	return m
}

func TestManager_Register(t *testing.T) {
	//given
	m := brokerMock(ManagerParam{
		WorkChanSize: 10,
		Loggers:      nil,
	})
	defer m.Stop()

	//when
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))

	//then
	assert.IsType(t, &bManager{}, m)
}

func TestManager_Up(t *testing.T) {
	//given
	m := brokerMock(ManagerParam{
		WorkChanSize: 10,
		Loggers:      nil,
	})
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))

	//when
	c := runtime.NumGoroutine()
	m.Up() // 고루틴으로 동작

	//then
	assert.Equal(t, c+1, runtime.NumGoroutine())
}

func TestManager_PushPopWork(t *testing.T) {
	//given
	m := brokerMock(ManagerParam{
		WorkChanSize: 10,
		Loggers:      nil,
	})
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))
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
	_ = m.PushWork(&pushCommand)
	_ = m.PushWork(&popCommand)

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
	m := brokerMock(ManagerParam{
		WorkChanSize: 10,
		Loggers:      nil,
	})
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))

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
	//given
	m := brokerMock(ManagerParam{
		WorkChanSize: 10,
		Loggers:      []broker.BLogger{sampleLog{}},
	})
	defer m.Stop()
	_ = m.Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))
	m.Up()

	//when
	command := bRedis.RPush{
		Key:   "key",
		Value: "value",
	}
	m.Stop() // 에러 유도
	err := m.PushWork(&command)

	//then
	assert.NotNil(t, err)
}
