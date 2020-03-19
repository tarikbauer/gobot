package domain

type Strategy interface {
	String() string
	Append(candleStick CandleStick)
	GetInfo() (result float64, buy bool, sell bool)
}
