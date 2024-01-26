package main

import (
	"context"
	"github.com/fredmayer/go-rest-api-template/internal/app"
	"github.com/fredmayer/go-rest-api-template/internal/config"
	"github.com/fredmayer/go-rest-api-template/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	//init config
	cfg := config.MustLoad()

	//init logger
	l := logging.Init(cfg.LogLevel)

	//load application
	application := app.New(app.DbOptions{cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName})

	//run server
	go func() {
		s := &http.Server{
			Addr:         cfg.HTTPAddr,
			ReadTimeout:  1 * time.Minute,
			WriteTimeout: 1 * time.Minute,
		}
		if err := application.EchoServer.StartServer(s); err != http.ErrServerClosed {
			l.Panicln(err)
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	//Завершаем все!
	application.EchoServer.Shutdown(ctx)
	application.Storage.Stop()
}
