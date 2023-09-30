package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"nats-streaming-consumer/internal/repository"
)

// @title           Nats-streaming consumer
// @version         1.0
// @description		https://docs.google.com/document/d/1f1Ni6u5mi4If5iyVMLQHjIAZJltDZc0QCGawitSSbxI/edit

// @host      localhost:8082
// @BasePath  /
func main() {
	logrus.Info("nats-streaming consumer start")

	// Init config
	if err := initConfig(); err != nil {
		logrus.Fatalf("Failed to init config: %s", err.Error())
	}

	// Config logger
	if err := configureLogger(viper.GetString("logger.log-level")); err != nil {
		logrus.Fatalf("Failed to configure logger: %s", err.Error())
	}

	// Config db
	db, err := repository.NewPostgresDb(readDbConfig())
	if err != nil {
		logrus.Fatalf("Failed to create db connection: %s", err.Error())
	}
	_ = db

	// Config repository

	// Config service

	// Init and run nats-streaming

	// Init and run httpServer
}

// Config
func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// Logger configuration
func configureLogger(logLevel string) error {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

// Database
func readDbConfig() repository.DbConfig {
	return repository.DbConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DbName:   viper.GetString("database.dbname"),
		SslMode:  viper.GetString("database.sslmode"),
	}
}
