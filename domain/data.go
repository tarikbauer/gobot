package domain

type Data interface {
	GetID() string
}

type StrategyData struct {
	Buy bool `json:"buy"`
	Sell bool `json:"sell"`
	Date string `json:"date"`
	Result float64 `json:"result"`
}

func (s StrategyData) GetID() string {
	return s.Date
}

type CandleStickData struct {
	Symbol string `json:"symbol"`
	Low float64 `json:"low"`
	High float64 `json:"high"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	Interval string `json:"interval"`
	OpenTime string `json:"openTime"`
}

func (c CandleStickData) GetID() string {
	return c.OpenTime
}
