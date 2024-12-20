package middlewares

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DbMiddleware(database *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if database == nil {
				panic("Instância do banco de dados está nula")
			}

			c.Set("db", database)

			return next(c)
		}
	}
}
