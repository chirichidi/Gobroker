package broker

type Broker interface { // ex) Manager
	Register(client interface{}, workerCount uint32, workChanSize uint32, loggers []BLogger) error
	Up()
	Worker(BROKER_WORKER_TYPE) Pusher
	Stop()
}

type BROKER_WORKER_TYPE int

const (
	BROKER_WORKER_NONE = BROKER_WORKER_TYPE(iota)
	BROKER_WORKER_REDIS
	BROKER_WORKER_CURL
)

func (r BROKER_WORKER_TYPE) String() string {
	switch r {
	case BROKER_WORKER_REDIS:
		return "BROKER_WORKER_REDIS"
	case BROKER_WORKER_CURL:
		return "BROKER_WORKER_CURL"
	default:
		return "BROKER_WORKER_NONE"
	}
}

type Pusher interface {
	PushWork(commander Commander) error
}

type Worker interface { // ex) bRedis
	Ping() bool
	Clear()
}

type Commander interface { // ex) RPush
	Cmd(Worker) ErrorCodeType
}

type Initiator interface { // Implement and use preconditions if necessary
	Init()
}

type CmdResult struct {
	Result interface{}
	Err    error
}

type LogLevel int

const (
	LOG_LEVEL_DEBUG = LogLevel(iota + 1)
	LOG_LEVEL_INFO
	LOG_LEVEL_ERROR
)

type BLogger interface { // 어플리케이션에서 BLogger 를 구현해서 사용할 것 (ex - scribe 로그)
	WriteBLog(string)
}
