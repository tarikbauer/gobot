package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"github.com/tarikbauer/gobot/app"
	"github.com/tarikbauer/gobot/application"
	"github.com/tarikbauer/gobot/domain"
	"github.com/tarikbauer/gobot/infrastructure/consumer"
	"github.com/tarikbauer/gobot/infrastructure/repository"
	"github.com/tarikbauer/gobot/infrastructure/strategies"
)

func getClient() *redis.Client {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		DB:       db,
		Password: os.Getenv("REDIS_PASSWORD"),
		Addr:    fmt.Sprint(os.Getenv("REDIS_HOST"), ":", os.Getenv("REDIS_PORT")),
	})
	return client
}

func getLogger() *logrus.Logger {
	level, _ := strconv.Atoi(os.Getenv("LEVEL"))
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.AllLevels[level])
	return logger
}

func getConsumer() domain.Consumer {
	client := http.Client{}
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	return consumer.NewBinanceConsumer(os.Getenv("SYMBOL"), domain.Minute, limit, &client)
}

func getRepository(client *redis.Client) domain.Repository {
	return repository.NewRedisRepository(client)
}

func getStrategies() []*domain.Strategy {
	var strategyList []*domain.Strategy
	strategy := strategies.NewAverageStrategy("10", 10)
	strategyList = append(strategyList, &strategy)
	return strategyList
}

func getService(client *redis.Client) application.Bot {
	return application.NewBotService(getConsumer(), getRepository(client), getStrategies()...)
}

func main() {
	logger := getLogger()
	client := getClient()
	service := getService(client)
	//cServer := make(chan error)
	cWorker := make(chan error)
	//port, _ := strconv.Atoi(os.Getenv("PORT"))
	//go app.RunServer(port, *logger, service, cServer)
	go app.RunWorker(*logger, service, cWorker)
	select {
	case <- cWorker:
		logger.Fatal("server error")
	}
}
