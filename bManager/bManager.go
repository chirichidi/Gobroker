package bManager

import (
	"dancechanlibrary/broker"
	"dancechanlibrary/broker/bCurl"
	"dancechanlibrary/broker/bRedis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

type bManager struct {
	workerMap map[broker.BROKER_WORKER_TYPE]*bWorker
	Loggers   []broker.BLogger
}

func (m *bManager) String() string {
	return "BROKER_bMANAGER"
}

type bWorker struct {
	workerCount uint32

	commanderChan     chan broker.Commander
	commanderChanSize uint32
	onDoneChan        chan struct{} // 종료 전파용

	Loggers []broker.BLogger
	broker.Worker
	broker.Pusher
}

func (m *bManager) Worker(workerType broker.BROKER_WORKER_TYPE) broker.Pusher {
	if bWorker, exist := m.workerMap[workerType]; exist == false {
		m.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_NOT_EXIST_WORKER,
			Source:    m.String(),
			Content:   fmt.Sprintf("workerType param: %d", workerType),
		})
		return nil
	} else {
		return bWorker
	}
}

func NewBManager(loggers []broker.BLogger) broker.Broker {
	m := new(bManager)
	m.init(loggers)
	return m
}

func (m *bManager) init(loggers []broker.BLogger) {
	//m.errorChan = make(chan broker.ErrorCodeType)
	m.workerMap = make(map[broker.BROKER_WORKER_TYPE]*bWorker)
	m.Loggers = loggers
	m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
		ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
		Source:    m.String(),
		Content:   fmt.Sprintf("init bManager success"),
	})
}

func workerType(client interface{}) broker.BROKER_WORKER_TYPE {
	switch client.(type) {
	case *redis.Client:
		return broker.BROKER_WORKER_REDIS
	case *http.Client:
		return broker.BROKER_WORKER_CURL
	default:
		return broker.BROKER_WORKER_NONE
	}
}

func (m *bManager) Register(client interface{}, workerCount uint32, workChanSize uint32, loggers []broker.BLogger) error {
	defer m.printPanicStack()

	// 0. check
	if workerType(client) == broker.BROKER_WORKER_NONE {
		m.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_NOT_EXIST_WORKER,
			Source:    m.String(),
			Content:   "invalid workerType",
		})
		return errors.New("invalid workerType")
	}
	if _, exist := m.workerMap[workerType(client)]; exist == true {
		return errors.New("broker duplicate register")
	}

	switch client.(type) {
	case *redis.Client:
		c := client.(*redis.Client)
		w := bRedis.InitBRedis(c)
		m.workerMap[broker.BROKER_WORKER_REDIS] = &bWorker{
			workerCount:       workerCount,
			commanderChan:     make(chan broker.Commander, workChanSize),
			commanderChanSize: workChanSize,
			Loggers:           loggers,
			Worker:            w,
			onDoneChan:        make(chan struct{}, 1),
		}
		m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
			Source:    m.String(),
			Content:   "redis client register ok",
		})
	case *http.Client:
		c := client.(*http.Client)
		w := bCurl.InitBHttp(c)
		m.workerMap[broker.BROKER_WORKER_CURL] = &bWorker{
			workerCount:       workerCount,
			commanderChan:     make(chan broker.Commander, workChanSize),
			commanderChanSize: workChanSize,
			Loggers:           loggers,
			Worker:            w,
			onDoneChan:        make(chan struct{}, 1),
		}
		m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
			Source:    m.String(),
			Content:   "http client register ok",
		})
	default:
		return errors.New("unexpected client type")
	}

	return nil
}

func (m *bManager) Up() {
	defer m.printPanicStack()

	for workerType, worker := range m.workerMap {
		atomicValue := uint32(0)
		for i := uint32(0); i < worker.workerCount; i++ {
			go func() {
				var workerTypeBytes []byte
				copy(workerTypeBytes, workerType.String())
				//TODO workerType 이거 잘 안되는데

				goRoutineIdx := atomic.LoadUint32(&atomicValue)
				atomic.AddUint32(&atomicValue, 1)

				m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
					ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
					Source:    m.String(),
					Content:   fmt.Sprintf("%s worker's goRoutine-%d upImpl", string(workerTypeBytes), goRoutineIdx),
				})
				for worker.upImpl() {
					break
				}
				m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
					ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
					Source:    m.String(),
					Content:   fmt.Sprintf("%s worker's goRoutine-%d upImpl break", string(workerTypeBytes), goRoutineIdx),
				})
			}()
		}
	}
}

func (w *bWorker) upImpl() bool {
	defer w.printPanicStack()
	isExit := false

LOOP:
	for {
		select {
		case commander := <-w.commanderChan:
			code := commander.Cmd(w.Worker)
			if code != broker.BROKER_ERROR_CODE_SUCCESS {
				w.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{ErrorCode: code, Source: workerType(w.Worker)})
			}
			w.writeLog(broker.LOG_LEVEL_DEBUG, broker.Error{ErrorCode: code, Source: workerType(w.Worker), Content: "command success"})
		case <-w.onDoneChan:
			isExit = true
			break LOOP
		}
	}

	return isExit
}

func (w *bWorker) PushWork(commander broker.Commander) error {
	defer w.printPanicStack()

	if w.commanderChan == nil {
		errCode := broker.BROKER_ERROR_CODE_CHANNEL_CLOSED
		w.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{
			ErrorCode: broker.ErrorCodeType(errCode),
			Source:    workerType(w.Worker),
			Content:   commander})
		return errors.New("broker channel closed")
	}

	select {
	case w.commanderChan <- commander:
		if _, ok := commander.(broker.Initiator); ok == true {
			commander.(broker.Initiator).Init()
		}
	default:
		errCode := broker.BROKER_ERROR_CODE_MANAGER_CHANNEL_OVERFLOW
		w.writeLog(broker.LOG_LEVEL_ERROR, broker.Error{
			ErrorCode: broker.ErrorCodeType(errCode),
			Source:    workerType(w.Worker),
			Content:   commander})
		return errors.New("broker channel overflow")
	}

	return nil
}

func (m *bManager) Stop() {
	defer m.printPanicStack()
	for _, worker := range m.workerMap {
		if worker.onDoneChan != nil {
			worker.onDoneChan <- struct{}{}
		}
	}
	time.Sleep(1 * time.Second)
	m.clear()
}

func (m *bManager) clear() {
	defer m.printPanicStack()

	for workerType, worker := range m.workerMap {
		worker.Clear()
		worker.commanderChan = nil
		worker.commanderChanSize = 0
		if worker.onDoneChan != nil {
			close(worker.onDoneChan)
		}

		time.Sleep(1 * time.Second)
		worker.onDoneChan = nil

		m.writeLog(broker.LOG_LEVEL_INFO, broker.Error{
			ErrorCode: broker.BROKER_ERROR_CODE_SUCCESS,
			Source:    m.String(),
			Content:   fmt.Sprintf("%s worker clear", workerType.String()),
		})
	}
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

func (w *bWorker) writeLog(level broker.LogLevel, e broker.Error) {
	defer w.printPanicStack()

	if w.Loggers != nil {
		for _, logger := range w.Loggers {
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

func (w *bWorker) printPanicStack() {
	if x := recover(); x != nil {
		if w.Loggers == nil {
			return
		}
		for _, logger := range w.Loggers {
			logger.WriteBLog(fmt.Sprintf("%v", x))
		}

		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			i++
			msg := fmt.Sprintf("PrintPanicStack. [func]: %s, [file]: %s, [line]: %d\n", runtime.FuncForPC(funcName).Name(), file, line)
			for _, logger := range w.Loggers {
				logger.WriteBLog(msg)
			}
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}
