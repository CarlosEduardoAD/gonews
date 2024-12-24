package task_queue

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/hibiken/asynq"
)

var (
	once     sync.Once
	instance *asynq.Client
)

func GenerateAsynqClient() *asynq.Client {
	once.Do(func() {
		host := env.GetEnv("REDIS_HOST", "gonews-redis")
		password := env.GetEnv("REDIS_PASSWORD", "redis-password")

		redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:6379", host), Password: password, TLSConfig: &tls.Config{}, DialTimeout: 10 * time.Second, // Aumente este valor se necessário
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second}
		client := asynq.NewClient(redisOpt)

		instance = client
	})

	return instance
}
