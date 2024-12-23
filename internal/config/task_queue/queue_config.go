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
		url := env.GetEnv("REDIS_URL", "")
		redisOpt := asynq.RedisClientOpt{Addr: url}
		client := asynq.NewClient(redisOpt)

		instance = client
	})

	return instance
}
