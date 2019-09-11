package broker

type Broker interface { // ex) Manager
	Register(interface{}) error
	PushWork(Commander) error
	Up()
	Stop()
}

type Worker interface { // ex) bRedis
	Ping() bool
	Clear()
}

type Commander interface { // ex) RPush
	Cmd(Worker) ErrorCodeType
	Result() CmdResult
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
