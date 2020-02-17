package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type bRedis struct {
	client *redis.Client
}

func InitBRedis(client *redis.Client) *bRedis {
	// 0. check
	_, err := client.Ping().Result()
	if err != nil {
		panic(broker.BROKER_ERROR_CODE_CLIENT_DISCONNECT)
	}

	// 1. register
	_client := &bRedis{
		client: client,
	}
	return _client
}

func (r *bRedis) Clear() {
	if r.client != nil {
		_ = r.client.Close()
	}
	r.client = nil
}

func (r *bRedis) Ping() bool {
	p := r.client.Ping()
	if p.Err() != nil {
		panic(broker.BROKER_ERROR_CODE_CLIENT_DISCONNECT)
	}

	return true
}
