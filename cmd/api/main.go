package main

import (
	delivery "comment-api/internal/delivery/server"
	"comment-api/pkg/handler"
	"comment-api/pkg/repository"
	"comment-api/pkg/service"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	filePath := "example.json"
	file, err := os.Create(filePath)
	if err != nil {
		logrus.Fatalf("failed to initialize %s: %s", filePath, err.Error())
	}
	defer file.Close()

	fmt.Println("Файл создан успешно.")

	localFile := &repository.LocalFile{Path: filePath}

	repos := repository.NewRepository(localFile)
	services := service.NewService(repos)
	newHandler := handler.NewHandler(services)

	srv := new(delivery.Server)
	go func() {
		if err := srv.Run("8088", newHandler.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running http sever: %s", err.Error())
		}
	}()

	logrus.Print("Comment-API Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Comment-API Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shugging down: %s", err.Error())
	}
}
