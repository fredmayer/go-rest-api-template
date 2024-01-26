package app

import (
	"context"
	dummyController "github.com/fredmayer/go-rest-api-template/internal/controllers/dummy"
	"github.com/fredmayer/go-rest-api-template/internal/services/dummy"
	"github.com/fredmayer/go-rest-api-template/internal/storage/mysql"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(ctx context.Context, e *echo.Echo, storage *mysql.Storage) {

	//TODO подключить сервис
	serviceDummy := dummy.New(storage.Dummy)
	handlers := dummyController.NewHandler(ctx, serviceDummy)

	rc := e.Group("/dummy")

	//Создание отчета
	rc.GET(
		"/:report_id",
		handlers.FindDummy,
	)
}
