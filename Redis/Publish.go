package Redis

import "Broker"

type Publish struct {
	channel    string
	message    string
	ResultChan chan int64 //optional
}

func (p *Publish) cmd(worker Broker.Worker) bool {
	defer Broker.PrintPanicStack()
	r := worker.(*_redis)

	n, err := r.Client.Publish(p.channel, p.message).Result()
	if err != nil {
		panic(err.Error())
		return false
	}

	if p.ResultChan != nil {
		if len(p.ResultChan) >= cap(p.ResultChan) {
			return false
		}
		p.ResultChan <- n
	}

	return true
}
