package middlewares

import (
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func QueueMiddleware(asynq_queue *asynq.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			c.Set("task_manager", asynq_queue)

			return next(c)
		}
	}
}
