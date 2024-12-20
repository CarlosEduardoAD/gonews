package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Request().URL.Path,
			"ip":     c.RealIP(),
		}).Info("Incoming request")
		return next(c)
	}
}
