package task_queue

import (
	"crypto/tls"
	"fmt"
	"log"
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
		host := env.GetEnv("REDIS_HOST", "gonews-redis")
		password := env.GetEnv("REDIS_PASSWORD", "redis-password")

		log.Println(host)
		log.Println(password)

		redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:6379", host), Password: password, TLSConfig: &tls.Config{}}
		client := asynq.NewClient(redisOpt)

		instance = client
	})

	return instance
}
