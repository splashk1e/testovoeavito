package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"testovoe.com/bootstrap"
	"testovoe.com/internal/handler"
	"testovoe.com/internal/repository"
	"testovoe.com/internal/server"
	"testovoe.com/internal/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	app := bootstrap.App()
	env := app.Env
	fmt.Print(env.PostgresConn)
	fmt.Print("hello")
	db := app.Postgres

	repository := repository.NewRepository(db)
	services := service.NewService(repository)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(env.ServerPort, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
	var str string
	go func() {
		fmt.Scan(&str)
		if str == "stop" {
			logrus.Print("App Shutting Down")
			if err := srv.Shutdown(context.Background()); err != nil {
				logrus.Errorf("error occured on server shutting down: %s", err.Error())

			}
		}
	}()
}
