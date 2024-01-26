package app

import (
	"github.com/fredmayer/go-rest-api-template/internal/storage/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	EchoServer *echo.Echo
	Storage    *mysql.Storage
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

	storage := mysql.New(options.DbHost, options.DbPort, options.DbUser, options.DbPassword, options.DbName)

	e := echo.New()
	e.Use(middleware.Recover())

	//TODO register routers

	return &App{
		EchoServer: e,
		Storage:    storage,
	}

}
