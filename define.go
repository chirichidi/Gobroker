package Broker

type Broker interface { // ex) Manager
	pushJob(Worker, Commander)
}

type Worker interface { // ex) redis
	Ping() bool
	PushChan(chan WorkerCommand, Commander)
}

type Commander interface { // ex) RPush
	cmd(Worker) bool
}

type WorkerCommand struct {
	Worker    Worker
	Commander Commander
}
