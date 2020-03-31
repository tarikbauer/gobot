package domain

import "time"

type StrategyData struct {
	Buy bool `json:"buy"`
	Sell bool `json:"sell"`
	Name string `json:"name"`
	Date time.Time `json:"date"`
	Result float64 `json:"result"`
}

type Strategy interface {
	String() string
	Append(candleStick CandleStick)
	GetInfo() (result float64, buy bool, sell bool)
}
