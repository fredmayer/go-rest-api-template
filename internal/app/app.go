package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	EchoServer *echo.Echo
}

type DbOptions struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func New(options DbOptions) *App {
	//ctx, cancelMain := context.WithCancel(context.Background())

	//init databases clickhouse and mongo

	//TODO init storage and DB connections

	e := echo.New()
	e.Use(middleware.Recover())

	//TODO register routers

	return &App{
		EchoServer: e,
	}

}
