package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/config/logging"
	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	"github.com/CarlosEduardoAD/go-news/internal/middlewares"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	http_server "github.com/CarlosEduardoAD/go-news/internal/routes/http_server"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup
var logger *logrus.Logger

type Context struct{}

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
		runWorkerServer(ctx)
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

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowCredentials: true,
			AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowCredentials, echo.HeaderAccessControlAllowHeaders},
		}))

		e.Use(middlewares.DbMiddleware(db))
		e.Use(middlewares.LogrusMiddleware)
		e.Static("/images", "internal/views/images")

		route_group := e.Group("/api/v1/emails")
		http_server.UseEmailRoutes(route_group)

		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Up and running!")
		})

		e.HEAD("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Up and running!")
		})

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

func runWorkerServer(ctx context.Context) {
	host := env.GetEnv("REDIS_HOST", "gonews-redis")
	password := env.GetEnv("REDIS_PASSWORD", "Carloseduardo08#")

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				fmt.Sprintf("%s:6379", host),
				redis.DialPassword(password))
		},
	}

	pool := work.NewWorkerPool(Context{}, 10, "go_news_namespace", redisPool)

	pool.Job("send_email", (*Context).HandleEmailDelivery)

	pool.Start()

	log.Println("Worker server iniciado com sucesso")

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// Stop the pool
	pool.Stop()

	logger.Infoln("Servidor Worker encerrado com sucesso")
}

func (c *Context) HandleEmailDelivery(work *work.Job) error {

	controller := jobcontrollers.NewJobController()
	err := controller.ExecuteTask(work)

	if err != nil {
		log.Println("err: ", err)
		logrus.Error(shared.GenerateError(err))
		return err
	}

	return nil
}
