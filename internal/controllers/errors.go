package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func ErrBadRequest(e echo.Context, field string, msg string) error {
	return e.JSON(http.StatusBadRequest, response{
		false,
		field + ": " + msg,
	})
}

func NotFound(e echo.Context) error {
	return e.JSON(http.StatusBadRequest, response{
		false,
		"Not Found",
	})
}
