package Broker

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Manager struct {
	CommandChan      chan *RedisParam // define 정의
	_CommandChanSize int32
	OnDoneChan       chan struct{} // 종료 전파용

	_redis *_redis
}

func (m *Manager) Init(commandChannelSize int32) *Manager {
	m.OnDoneChan = make(chan struct{}, 3)
	m.CommandChan = make(chan *RedisParam, commandChannelSize)

	m._CommandChanSize = commandChannelSize
	return m
}

func (m *Manager) Register(i interface{}) interface{} {
	defer PrintPanicStack()

	switch i.(type) {
	case *redis.Client:
		c := i.(*redis.Client)
		m._redis = InitRedis(c)
		return m._redis
	}

	panic("broker registering failed..")
	return nil
}

func (m *Manager) Up() {
	defer PrintPanicStack()
	fmt.Println("Broker UP")

	go func() {
		for {
			if m._UpImpl() {
				break
			}
		}
	}()
	fmt.Println("Broker DOWN")
}

func (m *Manager) _UpImpl() bool {
	defer PrintPanicStack()
	isExit := false

LOOP:
	for {
		select {
		case p := <-m.CommandChan:

			switch p.command {
			case COMMAND_REDIS_RPUSH:
				m._redis.RPush(p.Key, p.Values)
			case COMMAND_REDIS_BRPOP:
				p.ResultChan <- m._redis.BRPop(p.Timeout, p.Keys)
			case COMMAND_REDIS_BLPOP:
				p.ResultChan <- m._redis.BLPop(p.Timeout, p.Keys)
			case COMMAND_REDIS_PUBLISH:
				m._redis.Publish(p.Channel, p.Message)
			case COMMAND_REDIS_SUBSCRIBE:
				p.ResultChan <- m._redis.Subscribe(p.Channels)
			case COMMAND_REDIS_UNSUBSCRIBE:
				m._redis.UnSubscribe(p.PubSub, p.Channels)
			case COMMAND_REDIS_RECEIVEMESSAGE:
				p.ResultChan <- m._redis.ReceiveMessage(p.PubSub, p.OnDoneChan)
			default:
				panic("broker command is wrong")
			}
		case <-m.OnDoneChan:
			isExit = true
			break LOOP
		}
	}

	return isExit
}

func (m *Manager) pushCommand(param *RedisParam) {
	defer PrintPanicStack()

	if int32(len(m.CommandChan)) >= m._CommandChanSize {
		panic("CommandChan Overflow...")
	} else {
		m.CommandChan <- param
	}
}

func (m *Manager) Stop() {
	defer PrintPanicStack()

	m.OnDoneChan <- struct{}{}
	time.Sleep(1 * time.Second)

	m._Clear()
}

func (m *Manager) _Clear() {
	defer PrintPanicStack()

	m._redis._Clear()

	for len(m.OnDoneChan) > 0 {
		<-m.OnDoneChan
	}
	close(m.OnDoneChan)

	m._CommandChanSize = 0
}
