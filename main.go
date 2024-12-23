package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/config/logging"
	"github.com/CarlosEduardoAD/go-news/internal/config/task_queue"
	"github.com/CarlosEduardoAD/go-news/internal/middlewares"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	http_server "github.com/CarlosEduardoAD/go-news/internal/routes/http_server"
	"github.com/CarlosEduardoAD/go-news/internal/routes/queue"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup
var logger *logrus.Logger

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger = logging.GenerateLogrus()

	serverShutdown := make(chan struct{})

	wg.Add(2)

	go func() {
		defer wg.Done()
		runEchoServer(ctx, serverShutdown)
	}()

	go func() {
		defer wg.Done()
		runAsynqServer(ctx)
	}()

	<-ctx.Done()
	logger.Warningln("Iniciando graceful shutdown...")

	close(serverShutdown)

	wg.Wait()
	logger.Warningln("Aplicação encerrada com sucesso.")
}

func runEchoServer(ctx context.Context, shutdown chan struct{}) {

	go func() {
		e := echo.New()
		db := db.GenereateDB()
		db.AutoMigrate(emailmodel.EmailModel{})
		queue_client := task_queue.GenerateAsynqClient()

		e.Use(middlewares.DbMiddleware(db))
		e.Use(middlewares.QueueMiddleware(queue_client))
		e.Use(middlewares.LogrusMiddleware)

		route_group := e.Group("/api/v1/emails")
		http_server.UseEmailRoutes(route_group)

		e.Logger.Fatal(e.Start(":3000"))
	}()

	<-shutdown

	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := echo.New().Shutdown(ctxShutdown); err != nil {
		logger.Fatalf("Erro ao encerrar servidor HTTP: %v", err)
	}
	logger.Info("Servidor HTTP encerrado com sucesso")
}

func runAsynqServer(ctx context.Context) {
	url := env.GetEnv("REDIS_URL", "")
	redisOpt := asynq.RedisClientOpt{Addr: url}

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"default": 1,
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("send_email", queue.HandleEmailDelivery)

	errChan := make(chan error, 1)

	go func() {
		logger.Infoln("Servidor Asynq rodando...")
		errChan <- srv.Run(mux)
	}()

	select {
	case <-ctx.Done():
		logger.Warningln("Encerrando servidor Asynq...")
		srv.Shutdown()
	case err := <-errChan:
		if err != nil {
			log.Printf("Erro no servidor Asynq: %v", err)
			logger.Fatalf("Erro no servidor Asynq: %v", err)
		}
	}

	logger.Infoln("Servidor Asynq encerrado com sucesso")
}
