package http_server

import (
	"net/http"

	"github.com/CarlosEduardoAD/go-news/internal/controllers/email_controllers"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UseEmailRoutes(group *echo.Group) {
	group.POST("/check-in", CheckInRoute)
	group.GET("/authorize", AuthorizationRoute)
	group.GET("/dismiss", DismissRoute)
	group.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "aight!")
	})
}

func CheckInRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	e := &emailmodel.SaveEmailModelDTO{}

	if err := c.Bind(e); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, shared.GenerateError(err))
	}

	email := emailmodel.NewEmailModel(e.Email)
	controller := email_controllers.NewEmailController(db, nil)
	url, err := controller.CheckInEmail(email)

	if err != nil && err.Error() == "invalid email" {
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	// TODO: Criar objeto de response próprio
	return c.JSON(http.StatusCreated, map[string]interface{}{"url": url})
}

func AuthorizationRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	task_manager := c.Get("task_manager").(*asynq.Client)
	token := c.QueryParam("token")
	controller := email_controllers.NewEmailController(db, task_manager)
	err := controller.AuthorizeEmail(token)

	if err != nil && err.Error() == "invalid-token" {
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil && err.Error() == "email not found" {
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": "ACCEPTED",
	})
}

func DismissRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	// FIXME: isso daqui não pode ficar assim
	token := c.QueryParam("token")

	controller := email_controllers.NewEmailController(db, nil)
	err := controller.DismissEmail(token)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}
