package domain

import "time"

type Consumer interface {
	GetSymbol() string
	GetInterval() Interval
	GetData() (sleep time.Duration, candleStick []CandleStick, err error)
}
