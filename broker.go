package Broker

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Manager struct {
	CommandChan     chan *RedisParam // define 정의
	CommandChanSize int32
	OnDoneChan      chan struct{} // 종료 전파용

	redisWrapper *_redis
}

func (m *Manager) Init(commandChannelSize int32) *Manager {
	m.OnDoneChan = make(chan struct{}, 3)
	m.CommandChan = make(chan *RedisParam, commandChannelSize)

	m.CommandChanSize = commandChannelSize
	return m
}

func (m *Manager) Register(i interface{}) interface{} {
	defer PrintPanicStack()

	switch i.(type) {
	case *redis.Client:
		c := i.(*redis.Client)
		m.redisWrapper = InitRedis(c)
		return m.redisWrapper
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
				m.redisWrapper.RPush(p.Key, p.Values)
			case COMMAND_REDIS_BRPOP:
				p.ResultChan <- m.redisWrapper.BRPop(p.Timeout, p.Keys)
			case COMMAND_REDIS_BLPOP:
				p.ResultChan <- m.redisWrapper.BLPop(p.Timeout, p.Keys)
			case COMMAND_REDIS_PUBLISH:
				m.redisWrapper.Publish(p.Channel, p.Message)
			case COMMAND_REDIS_SUBSCRIBE:
				p.ResultChan <- m.redisWrapper.Subscribe(p.Channels)
			case COMMAND_REDIS_UNSUBSCRIBE:
				m.redisWrapper.UnSubscribe(p.PubSub, p.Channels)
			case COMMAND_REDIS_RECEIVEMESSAGE:
				p.ResultChan <- m.redisWrapper.ReceiveMessage(p.PubSub, p.OnDoneChan)
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

	if int32(len(m.CommandChan)) >= m.CommandChanSize {
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

	m.redisWrapper._Clear()

	for len(m.OnDoneChan) > 0 {
		<-m.OnDoneChan
	}
	close(m.OnDoneChan)

	m.CommandChanSize = 0
}
