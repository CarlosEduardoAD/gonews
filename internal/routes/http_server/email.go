package http_server

import (
	"errors"
	"log"
	"net/http"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/controllers/email_controllers"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UseEmailRoutes(group *echo.Group) {
	group.POST("/check-in", CheckInRoute)
	group.GET("/authorize", AuthorizationRoute)
	group.GET("/dismiss", DismissRoute)
	group.POST("/resend", ResendRoute)
	group.GET("/verify", VerifyRoute)
}

type CheckInRequest struct {
	Email string `json:"email"`
}

func CheckInRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	e := &CheckInRequest{}

	if err := c.Bind(e); err != nil {
		log.Println("err: ", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, shared.GenerateError(err))
	}

	controller := email_controllers.NewEmailController(db)

	select_email := emailmodel.EmailModel{}
	select_email.Email = e.Email
	email, err := select_email.SelectOneByEmail(db, e.Email)

	if err != nil && err.Error() != "record not found" {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	if email != nil {
		return echo.NewHTTPError(http.StatusConflict, shared.GenerateError(errors.New("email already exists")))
	}

	token, err := controller.CheckInEmail(e.Email)

	if err != nil && err.Error() == "invalid email" {
		log.Println("err: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil {
		log.Println("err: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	// TODO: Criar objeto de response próprio
	return c.JSON(http.StatusCreated, map[string]interface{}{"token": token})
}

func AuthorizationRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	token := c.QueryParam("token")
	controller := email_controllers.NewEmailController(db)
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

	next_front_url := env.GetEnv("NEXT_FRONT_URL", "http://localhost:3001")

	return c.Redirect(http.StatusMovedPermanently, next_front_url+"/accepted?token="+token)
}

func VerifyRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	token := c.QueryParam("token")
	controller := email_controllers.NewEmailController(db)
	err := controller.VerifyEmail(token)

	if err != nil {
		log.Println("err: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

func ResendRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	token := c.QueryParam("token")

	controller := email_controllers.NewEmailController(db)
	err := controller.ResendEmail(token)

	if err != nil && err.Error() == "invalid-token" {
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil && err.Error() == "email not found" {
		return echo.NewHTTPError(http.StatusBadRequest, shared.GenerateError(err))
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

func DismissRoute(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	// FIXME: isso daqui não pode ficar assim
	token := c.QueryParam("token")

	controller := email_controllers.NewEmailController(db)
	err := controller.DismissEmail(token)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.GenerateError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}
