package domain

import "errors"

var (
	InvalidConsumedContent = errors.New("invalid consumed content")
	UnexpectedResponseFormat = errors.New("unexpected response format")
)

type Consumer interface {
	GetLimit() int
	GetSymbol() string
	GetInterval() Interval
	GetData() (candleStick []CandleStick, err error)
}
