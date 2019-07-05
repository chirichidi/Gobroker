package Redis

import "Broker"

type RPush struct {
	key   string
	value string // TODO 성능때문에 single parameter 사용
}

func (rp *RPush) cmd(worker Broker.Worker) bool {
	defer Broker.PrintPanicStack()
	r := worker.(*_redis)

	_, err := r.Client.RPush(rp.key, rp.value).Result()
	if err != nil {
		panic(err.Error())
		return false
	}

	return true
}
