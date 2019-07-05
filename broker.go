package Broker

import (
	"github.com/go-redis/redis"
	"time"
)

type Manager struct {
	WorkChan   chan WorkerCommand
	OnDoneChan chan struct{} // 종료 전파용

	workChanSize int32
}

func (m *Manager) Init(workChanSIze int32) *Manager {
	m.WorkChan = make(chan WorkerCommand, workChanSIze)
	m.OnDoneChan = make(chan struct{}, 3)
	m.workChanSize = workChanSIze
	return m
}

func (m *Manager) Register(i interface{}) Broker {
	defer PrintPanicStack()

	switch i.(type) {
	case *redis.Client:
		//c := i.(*redis.Client)

		// 한번만 만든다고 가정
		//m._redis = InitRedis(c)
		//return m._redis
	}

	return nil
}

func (m *Manager) Run() {
	defer PrintPanicStack()

	go func() {
		for m._RunImpl() {
			break
		}
	}()
}

func (m *Manager) _RunImpl() bool {
	defer PrintPanicStack()
	isExit := false

LOOP:
	for {
		select {
		case wc := <-m.WorkChan:
			wc.Commander.cmd(wc.Worker)
		case <-m.OnDoneChan:
			isExit = true
			break LOOP
		}
	}

	return isExit
}

func (m *Manager) pushWork(cc WorkerCommand) {
	defer PrintPanicStack()

	if int32(len(m.WorkChan)) >= m.workChanSize {
		panic("workChan Overflow..")
	} else {
		cc.Worker.PushChan(m.WorkChan, cc.Commander)
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

	m.workChanSize = 0

	if len(m.WorkChan) > 0 {
		time.Sleep(5 * time.Second)

		for len(m.WorkChan) > 0 {
			<-m.WorkChan
		}
		close(m.WorkChan)
	}

	for len(m.OnDoneChan) > 0 {
		<-m.OnDoneChan
	}
	close(m.OnDoneChan)
}
