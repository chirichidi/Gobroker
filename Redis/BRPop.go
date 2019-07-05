package Redis

import (
	"Broker"
	"time"
)

type BRPop struct {
	time       time.Duration
	key        string // TODO 성능때문에 single parameter 사용
	ResultChan chan string
}

func (bp *BRPop) cmd(worker Broker.Worker) bool {
	defer Broker.PrintPanicStack()
	if len(bp.ResultChan) >= cap(bp.ResultChan) {
		return false
	}

	r := worker.(*_redis)
	brPop := r.Client.BRPop(bp.time, bp.key)
	val := brPop.Val()

	bp.ResultChan <- val[1]
	return true
}
