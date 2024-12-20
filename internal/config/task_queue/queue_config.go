package task_queue

import (
	"sync"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/hibiken/asynq"
)

var (
	once     sync.Once
	instance *asynq.Client
)

func GenerateAsynqClient() *asynq.Client {
	once.Do(func() {
		password := env.GetEnv("REDIS_PASSWORD", "redis-password")
		client := asynq.NewClient(asynq.RedisClientOpt{Addr: "gonews-redis:6379", Password: password})

		instance = client
	})

	return instance
}
