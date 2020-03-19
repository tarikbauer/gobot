package application

import (
	"github.com/tarikbauer/gobot/infrastructure/flusher"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tarikbauer/gobot/domain"
)

type Bot interface {
	Run()
}

type botService struct {
	domain.Flusher
	domain.Consumer
	logger logrus.Logger
	strategies []*domain.Strategy
}

func NewBotService(
	logger logrus.Logger,
	flusher domain.Flusher,
	consumer domain.Consumer,
	strategies ... *domain.Strategy,
) Bot {
	return &botService{
		logger:     logger,
		Flusher:    flusher,
		Consumer:   consumer,
		strategies: strategies,
	}
}

func (b *botService) flushStrategies(strategies []domain.StrategyData, wg *sync.WaitGroup, c chan<- error) {
	defer wg.Done()
	var operation domain.Operation = "application.bot.flushStrategy"
	err := b.FlushStrategies(strategies)
	if err != nil {
		c <- domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	} else {
		c <- nil
	}
}

func (b *botService) flushCandleSticks(candleSticks []domain.CandleStickData, wg *sync.WaitGroup, c chan<- error) {
	defer wg.Done()
	var operation domain.Operation = "application.bot.flushData"
	err := b.FlushCandleSticks(candleSticks)
	if err != nil {
		c <- domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	} else {
		c <- nil
	}

}

func (b *botService) appendStrategy(strategies *[]domain.StrategyData, candleStick domain.CandleStick, strategy *domain.Strategy) {
	s := *strategy
	s.Append(candleStick)
	result, buy, sell := s.GetInfo()
	*strategies = append(*strategies, domain.StrategyData{
		Buy:    buy,
		Sell:   sell,
		Result: result,
		Date:   candleStick.OpenTime.Format(time.RFC3339),
	})
}

func (b *botService) appendCandleStick(candleSticksData *[]domain.CandleStickData, candleStick domain.CandleStick) {
	*candleSticksData = append(*candleSticksData, domain.CandleStickData{
		Symbol:   b.GetSymbol(),
		Low:      candleStick.Low,
		High:     candleStick.High,
		Open:     candleStick.Open,
		Close:    candleStick.Close,
		Interval: b.GetInterval().String(),
		OpenTime: candleStick.OpenTime.Format(time.RFC3339),
	})
}

// TODO: improve performance
func (b *botService) Run() {
	var lastAt time.Time
	var operation domain.Operation = "application.bot.Run"
	log := b.logger.WithFields(logrus.Fields{"symbol": b.GetSymbol(), "interval": b.GetInterval()})
	log.Info("starting bot")
	for {
		sleep, candleSticks, err := b.GetData()
		t := time.Now()
		if err != nil {
			log.Error(domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err))
			break
		}
		var strategies []domain.StrategyData
		var candleSticksData []domain.CandleStickData
		for _, candleStick := range candleSticks {
			if lastAt.After(candleStick.OpenTime) || lastAt == candleStick.OpenTime {
				continue
			}
			lastAt = candleStick.OpenTime
			b.appendCandleStick(&candleSticksData, candleStick)
			for _, strategy := range b.strategies {
				b.appendStrategy(&strategies, candleStick, strategy)
			}
		}
		_ = b.Flush(strategies, candleSticksData)
		log.Info("strategies and data flushed")
		flusher.Render()
		log.Info("chart rendered")
		sleep -= time.Since(t)
		if sleep > 0 {
			log.Info("sleeping for ", sleep)
			time.Sleep(sleep)
		}
	}
	log.Info("stopping bot")
}
