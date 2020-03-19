package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tarikbauer/gobot/application"
	"github.com/tarikbauer/gobot/domain"
	"github.com/tarikbauer/gobot/infrastructure"
	"github.com/tarikbauer/gobot/infrastructure/flusher"
	"os"
	"time"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	c := infrastructure.NewBinanceConsumer("LTCBTC", domain.MINUTE, 30*time.Second)
	//f := flusher.NewElasticsearchFlusher("http://127.0.0.1:9200", 30*time.Second)
	f := flusher.NewLocalFlusher("/Users/tarikbauer/go/src/github.com/tarikbauer/gobot/asd.html")
	s := infrastructure.NewAverageStrategy("50MAS", 50)
	b := application.NewBotService(*log, f, c, &s)
	b.Run()
}


