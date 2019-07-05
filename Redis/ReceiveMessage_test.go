package Redis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//redis Init/Register 같은거 수정하고, 테스트케이스 만들고
//Run() 부분 테스트케이스 만들고
//*가능하면 쉬운 버전까지 구현해서 성능 비교

func Test_redis_ReceiveMessage(t *testing.T) {
	r := RedisMock()
	ps := r.Subscribe("topic")
	time.Sleep(10 * time.Millisecond)

	command := ReceiveMessage{
		pubSub:      ps,
		_onDoneChan: make(chan struct{}, 1),
		ResultChan:  make(chan *redis.Message, 10),
	}
	go command.cmd(InitRedis(r))
	time.Sleep(1 * time.Millisecond)

	r.Publish("topic", "message")

	m := <-command.ResultChan
	assert.Equal(t, "message", m.Payload)

	command._onDoneChan <- struct{}{}
}

func Test_redis_ReceiveMessage2(t *testing.T) { // 종료 테스트
	r := RedisMock()
	ps := r.Subscribe("topic3")
	time.Sleep(10 * time.Millisecond)

	command := ReceiveMessage{
		pubSub:      ps,
		_onDoneChan: make(chan struct{}, 1),
		ResultChan:  make(chan *redis.Message, 10),
	}
	go command.cmd(InitRedis(r))
	time.Sleep(1 * time.Millisecond)

	command._onDoneChan <- struct{}{}
	time.Sleep(1 * time.Millisecond)

	assert.Equal(t, 0, len(command._onDoneChan))
}
