package Broker

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Broker interface {
	Up()
}

type Manager struct {
	OnDoneChan   chan struct{} // 종료 전파용

	_redis *_redis
}

func (m *Manager) Init(wc int32) *Manager {
	m.OnDoneChan = make(chan struct{}, 3)

	return m
}

func (m *Manager) Register(i interface{}) Broker {
	defer PrintPanicStack()

	switch i.(type) {
	case *redis.Client:
		c := i.(*redis.Client)

		// 한번만 만든다고 가정
		m._redis = InitRedis(c)
		return m._redis
	}

	return nil
}

func (m *Manager) Up() {
	fmt.Println("Broker UP")
	defer PrintPanicStack()
}

func (m *Manager) Stop() {
	defer PrintPanicStack()

	m.OnDoneChan <- struct{}{}
	time.Sleep(1*time.Second)

	m._Clear()
}

func (m *Manager) _Clear() {
	defer PrintPanicStack()

	m._redis._Clear()

	for len(m.OnDoneChan) > 0 {
		<- m.OnDoneChan
	}
	close(m.OnDoneChan)
}

