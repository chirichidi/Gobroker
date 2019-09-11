package bManager

import (
	"dancechanlibrary/broker"
	"dancechanlibrary/broker/bRedis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"runtime"
	"time"
)

type bManager struct {
	commanderChan chan broker.Commander
	//errorChan     chan broker.ErrorCodeType
	onDoneChan chan struct{} // 종료 전파용
	worker     broker.Worker
	ManagerParam
}

type ManagerParam struct {
	WorkChanSize int32
	Loggers      []broker.BLogger
}

func NewBManager(param ManagerParam) broker.Broker {
	m := new(bManager)
	m.init(&param)
	return m
}

func (m *bManager) init(param *ManagerParam) {
	m.commanderChan = make(chan broker.Commander, param.WorkChanSize)
	//m.errorChan = make(chan broker.ErrorCodeType)
	m.onDoneChan = make(chan struct{})
	m.WorkChanSize = param.WorkChanSize
	m.Loggers = param.Loggers
	m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
		ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
		Content:   fmt.Sprintf("init success commanderChan: %d, workerChanSize: %d", param.WorkChanSize, param.WorkChanSize),
	})
}

func (m *bManager) Register(client interface{}) error {
	defer m.printPanicStack()
	if m.worker != nil {
		return errors.New("broker duplicate register")
	}

	switch client.(type) {
	case *redis.Client:
		c := client.(*redis.Client)
		w := bRedis.InitBRedis(c)
		m.worker = w
		m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
			Content:   "redis client register ok",
		})
	default:
		return errors.New("unexpected client type")
	}

	return nil
}

func (m *bManager) Up() {
	defer m.printPanicStack()

	go func() {
		for m.upImpl() {
			break
		}
	}()
}

func (m *bManager) upImpl() bool {
	defer m.printPanicStack()
	isExit := false

LOOP:
	for {
		select {
		case commander := <-m.commanderChan:
			code := commander.Cmd(m.worker)
			if code != broker.BROKER_ERROR_CODE_SUCCESS {
				//m.errorChan <- broker.ErrorCodeType(code)
				m.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{ErrorCode: broker.ErrorCodeType(code)})
			}
			m.writeLog(broker.LOG_LEVEL_DEBUG, broker.Error{ErrorCode: broker.ErrorCodeType(code), Content: "command success"})
		case <-m.onDoneChan:
			isExit = true
			break LOOP
		}
	}

	return isExit
}

func (m *bManager) PushWork(commander broker.Commander) error {
	defer m.printPanicStack()

	if m.commanderChan == nil {
		errCode := broker.BROKER_ERROR_CODE_CHANNEL_CLOSED
		m.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{ErrorCode: broker.ErrorCodeType(errCode), Content: commander})
		return errors.New("broker channel closed")
	}

	select {
	case m.commanderChan <- commander:
	default:
		errCode := broker.BROKER_ERROR_CODE_MANAGER_CHANNEL_OVERFLOW
		m.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{ErrorCode: broker.ErrorCodeType(errCode), Content: commander})
		return errors.New("broker channel overflow")
	}

	return nil
}

func (m *bManager) Stop() {
	defer m.printPanicStack()
	if m.onDoneChan != nil {
		close(m.onDoneChan)
	}

	time.Sleep(1 * time.Second)

	m.clear()
}

func (m *bManager) clear() {
	defer m.printPanicStack()
	if m.worker != nil {
		m.worker.Clear()
	}
	m.commanderChan = nil
	m.onDoneChan = nil
	m.WorkChanSize = 0
}

func (m *bManager) writeLog(level broker.LogLevel, e broker.Error) {
	defer m.printPanicStack()

	if m.Loggers != nil {
		for _, logger := range m.Loggers {
			out, err := json.Marshal(e)
			if err != nil {
				logger.WriteBLog("write log error - json marshal")
				continue
			}

			// TODO: Config에 따라 레벨별로 로그 처리
			logger.WriteBLog(string(out))
		}
	}
}

func (m *bManager) printPanicStack() {
	if x := recover(); x != nil {
		if m.Loggers == nil {
			return
		}
		for _, logger := range m.Loggers {
			logger.WriteBLog(fmt.Sprintf("%v", x))
		}

		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			i++
			msg := fmt.Sprintf("PrintPanicStack. [func]: %s, [file]: %s, [line]: %d\n", runtime.FuncForPC(funcName).Name(), file, line)
			for _, logger := range m.Loggers {
				logger.WriteBLog(msg)
			}
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}
