package application

import (
	"bytes"
	"context"
	"html/template"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tarikbauer/gobot/domain"
)

type Bot interface {
	Run(ctx context.Context) error
	Render(ctx context.Context) ([]byte, error)
}

type botService struct {
	domain.Repository
	domain.Consumer
	strategies []*domain.Strategy
}

func NewBotService(
	consumer domain.Consumer,
	repository domain.Repository,
	strategies ... *domain.Strategy,
) Bot {
	return &botService{
		Repository: repository,
		Consumer:   consumer,
		strategies: strategies,
	}
}

func (b *botService) flushStrategy(strategy domain.StrategyData, wg *sync.WaitGroup, c chan<- error) {
	defer wg.Done()
	c <- b.FlushStrategy(strategy)

}

func (b *botService) flushCandleSticks(candleSticks []domain.CandleStick, wg *sync.WaitGroup, c chan<- error) {
	defer wg.Done()
	c <- b.FlushCandleSticks(candleSticks)
}

func (b *botService) Render(ctx context.Context) ([]byte, error) {
	var contents []data
	logger := ctx.Value("logger").(logrus.Logger)
	tmpl := template.New("CandleStick Chart")
	tmpl, _ = tmpl.Parse(html)
	buf := bytes.NewBuffer([]byte{})
	logger.Info("getting candlesticks")
	candleSticks, err := b.RetrieveCandleSticks()
	if err != nil {
		return []byte{}, err
	}
	for _, candleStick := range candleSticks {
		var results []float64
		for _, strategy := range b.strategies {
			s := *strategy
			logger.WithFields(logrus.Fields{"strategy": s.String()}).Info("getting strategy")
			strategyData, err := b.RetrieveStrategy(s.String(), candleStick.OpenTime)
			if err != nil {
				return []byte{}, err
			}
			results = append(results, strategyData.Result)
		}
		contents = append(contents, data{
			ID:     "",
			Low:    candleStick.Low,
			Open:   candleStick.Open,
			Close:  candleStick.Close,
			High:   candleStick.High,
			Results: results,
		})
	}
	_ = tmpl.Execute(buf, contents)
	return buf.Bytes(), nil
}

func (b *botService) Run(ctx context.Context) (err error) {
	logger := ctx.Value("logger").(logrus.Logger)
	var lastAt time.Time
	for {
		c := make(chan error)
		wg := sync.WaitGroup{}
		logger.Info("fetching data")
		candleSticks, err := b.GetData()
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Info("data fetched")
		wg.Add(1)
		go b.flushCandleSticks(candleSticks, &wg, c)
		for _, candleStick := range candleSticks {
			if lastAt.After(candleStick.OpenTime) || lastAt == candleStick.OpenTime {
				continue
			}
			lastAt = candleStick.OpenTime
			for _, strategy := range b.strategies {
				wg.Add(1)
				s := *strategy
				result, buy, sell := s.GetInfo()
				go b.flushStrategy(domain.StrategyData{
					Buy:    buy,
					Sell:   sell,
					Result: result,
					Name: s.String(),
					Date:   candleStick.OpenTime,
				}, &wg, c)
			}
		}
		wg.Wait()
		close(c)
		for err = range c {
			if err != nil {
				logger.Error(err)
				return err
			}
		}
		logger.Info("data flushed")
	}
}
