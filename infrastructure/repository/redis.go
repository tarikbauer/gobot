package repository

import (
	"encoding/json"
	"net"
	"time"

	"github.com/go-redis/redis"

	"github.com/tarikbauer/gobot/domain"
)

type redisRepository struct {
	hash string
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.Repository {
	return &redisRepository{hash: "candlesticks", client: client}
}

func (rf *redisRepository) FlushStrategy(strategy domain.StrategyData) error {
	resp := rf.client.HSet(strategy.Name, strategy.Date.Format(time.RFC3339), strategy)
	_, ok := resp.Err().(*net.OpError)
	if ok {
		return resp.Err()
	}
	return nil
}

func (rf *redisRepository) FlushCandleSticks(candleSticks []domain.CandleStick) error {
	var data []interface{}
	for _, candleStick := range candleSticks {
		data = append(data, candleStick)
	}
	resp := rf.client.LPush(rf.hash, data...)
	_, ok := resp.Err().(*net.OpError)
	if ok {
		return resp.Err()
	}
	return nil
}

func (rf *redisRepository) RetrieveCandleSticks() ([]domain.CandleStick, error) {
	var candleSticks []domain.CandleStick
	resp := rf.client.LRange(rf.hash, 0, -1)
	_, ok := resp.Err().(*net.OpError)
	if ok {
		return nil, resp.Err()
	}
	for _, value := range resp.Val() {
		var currentCandleStick domain.CandleStick
		err := json.Unmarshal([]byte(value), &currentCandleStick)
		if err != nil {
			return nil, err
		}
		candleSticks = append(candleSticks, currentCandleStick)
	}
	return candleSticks, nil
}

func (rf *redisRepository) RetrieveStrategy(name string, date time.Time) (domain.StrategyData, error) {
	var strategy domain.StrategyData
	resp := rf.client.HGet(name, date.Format(time.RFC3339))
	_, ok := resp.Err().(*net.OpError)
	if ok {
		return strategy, resp.Err()
	}
	err := json.Unmarshal([]byte(resp.Val()), &strategy)
	return strategy, err
}
