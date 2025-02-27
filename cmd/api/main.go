package main

import (
	delivery "comment-api/internal/delivery/server"
	"comment-api/pkg/handler"
	"comment-api/pkg/repository/postgres"
	"comment-api/pkg/service"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// @title Comment API
// @version 1.0
// @description API Server for CommentApi Application

// @host localhost:8088
// @BasePath /

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		DisableColors:   false,
		ForceQuote:      true,
		PadLevelText:    true,
	})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	logrus.Print(viper.GetString("db.host"))
	logrus.Print(os.Getenv(viper.GetString("db.port")))
	logrus.Print(os.Getenv(viper.GetString("db.password")))
	logrus.Print(viper.GetString("db.dbname"))
	logrus.Print(viper.GetString("db.sslmode"))

	db, err := postgres.NewPostgresDB(
		postgres.Config{
			Host:     viper.GetString("db.host"),
			Port:     os.Getenv(viper.GetString("db.port")),
			Username: os.Getenv(viper.GetString("db.username")),
			Password: os.Getenv(viper.GetString("db.password")),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
	)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	reposPostgres := postgres.NewRepositoryPostgres(db)
	filter := service.NewProfaneFilterService(os.Getenv("PROFANE_WORDS_API"))
	services := service.NewService(reposPostgres, filter)
	newHandler := handler.NewHandler(services)

	srv := new(delivery.Server)
	go func() {
		if err := srv.Run(os.Getenv("COMMENT_API_PORT"), newHandler.InitRoutes()); err != nil {
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
