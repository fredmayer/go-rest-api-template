package app

import (
	"context"
	dummyController "github.com/fredmayer/go-rest-api-template/internal/controllers/dummy"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(ctx context.Context, e *echo.Echo) {

	//TODO подключить сервис
	//service := reports.NewService()
	handlers := dummyController.NewHandler(ctx, service)

	rc := e.Group("/dummy")

	//Создание отчета
	rc.GET(
		"/:report_id",
		handlers.FindDummy,
	)
}
