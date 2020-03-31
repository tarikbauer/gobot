package domain

import "time"

type Repository interface {
	FlushStrategy(strategies StrategyData) error
	RetrieveCandleSticks() ([]CandleStick, error)
	FlushCandleSticks(candleSticks []CandleStick) error
	RetrieveStrategy(name string, date time.Time) (StrategyData, error)
}
