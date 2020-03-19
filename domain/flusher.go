package domain

type Flusher interface {
	Flush(strategies []StrategyData, candleSticks []CandleStickData) error
	FlushStrategies(strategies []StrategyData) error
	FlushCandleSticks(candleSticks []CandleStickData) error
}
