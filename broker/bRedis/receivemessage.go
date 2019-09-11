package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type ReceiveMessage struct { // PubSub.ReceiveMessage() 가 리슨(blocking) 함수이므로 application 에서 고루틴으로 실행할 것
	PubSub     *redis.PubSub
	OnDoneChan chan struct{}
	ResultChan chan broker.CmdResult `json:"-"`
}

// 사용하는쪽에서 고루틴 회수 관리 잘해야함
func (rm *ReceiveMessage) Cmd(worker broker.Worker) broker.ErrorCodeType {
	go rm.exitSignal()
	go rm.receiveProcess()

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (rm *ReceiveMessage) Result() broker.CmdResult {
	return <-rm.ResultChan
}

func (rm *ReceiveMessage) receiveProcess() {
	for {
		if rm.receiveProcessImpl() {
			return
		}
	}
}

func (rm *ReceiveMessage) receiveProcessImpl() bool {
	isWantedTermination := false

	for {
		if rm.PubSub == nil {
			isWantedTermination = true
			return isWantedTermination
		}
		m, err := rm.PubSub.ReceiveMessage() //blocking method
		if m == nil {
			rm.ResultChan <- broker.CmdResult{Result: nil, Err: err}
		} else {
			rm.ResultChan <- broker.CmdResult{Result: m.Payload, Err: err}
		}
	}
}

func (rm *ReceiveMessage) exitSignal() {
	select {
	case <-rm.OnDoneChan:
		_ = rm.PubSub.Close()
		rm.PubSub = nil
	}
}
