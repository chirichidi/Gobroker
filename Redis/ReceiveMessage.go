package Redis

import (
	"Broker"
	"github.com/go-redis/redis"
)

type ReceiveMessage struct { // pubSub.ReceiveMessage() 가 블락킹 함수이므로, application 에서 고루틴으로 실행해야함 (리슨용 함수)
	pubSub *redis.PubSub
	_onDoneChan chan struct{}
	ResultChan chan *redis.Message
}

func (rm *ReceiveMessage) cmd(worker Broker.Worker) bool {
	defer Broker.PrintPanicStack()
	defer func() {
		if rm.pubSub != nil {
			rm.pubSub.Close()
		}
	}()

	if len(rm.ResultChan) >= cap(rm.ResultChan) {
		return false
	}

	아아아아아 여기 동작 이상해~~~~~~~



	for rm.exitSignal() != false {
		go func() bool {
			m, err := rm.pubSub.ReceiveMessage()
			if m == nil {
				return false
			}
			if err != nil {
				panic(err.Error())
				return false
			}
			rm.ResultChan <- m
			return true
		}()
	}

	return true
}

func (rm *ReceiveMessage) exitSignal() bool {
	go func() bool {
	LOOP:
		for {
			select {
			case _ = <-rm._onDoneChan:
				if rm.pubSub != nil {
					rm.pubSub.Close()
					rm.pubSub = nil
					break LOOP
				}
			}
		}
		return true
	}()

	return false
}