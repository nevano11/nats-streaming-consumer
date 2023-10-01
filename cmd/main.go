package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"nats-streaming-consumer/internal/handler"
	"nats-streaming-consumer/internal/repository"
	"nats-streaming-consumer/internal/service"
	"net/http"
	"time"
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

	// Repository
	modelRepository := repository.NewPostgresModelRepository(db)
	cachedRepository := repository.NewCachedRepository(modelRepository)

	go func() {
		logrus.Info("Start filling cache from db")
		err := cachedRepository.FillCacheFromRepository()
		if err != nil {
			logrus.Errorf("Failed to load cache from db")
			return
		}
		logrus.Info("End filling cache from db")
	}()

	// Services
	dbService := service.NewDbServiceDb(repository.NewRepository(cachedRepository))
	senderService, err := service.NewNatsSenderService("modelChannel")
	if err != nil {
		logrus.Fatalf("Failed to create NatsSenderService: %s", err.Error())
	}
	consumeService, err := service.NewNatsConsumeService("modelChannel", dbService)
	if err != nil {
		logrus.Fatalf("Failed to create NatsConsumeService: %s", err.Error())
	}
	defer senderService.Close()
	defer consumeService.Close()

	// handler
	han := handler.NewHandler(dbService, senderService)

	// Routes
	routes := han.InitRoutes()

	// Server
	server := createServer(viper.GetString("server.port"), routes)
	logrus.Infof("Server running on http://localhost%s", server.Addr)
	logrus.Infof("Swagger: http://localhost%s/swagger/index.html", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start server: %s", err.Error())
	}
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

// Create server
func createServer(port string, routes *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           routes,
		ReadHeaderTimeout: 2 << 20,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
}
