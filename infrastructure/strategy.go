package infrastructure

import "github.com/tarikbauer/gobot/domain"

type trend struct {
	uptrendCount, downtrendCount int
	lastLow, lastHigh, lastValue float64
	uptrend, downtrend, higherLow, lowerHigh bool
}

func (t *trend) evaluate(candleStick domain.CandleStick) {
	if candleStick.Close < candleStick.Open && candleStick.Open > t.lastValue {
		if candleStick.Open < t.lastHigh {
			t.lowerHigh = true
		} else {
			t.lowerHigh = false
		}
		t.lastHigh = candleStick.Open
	}
	if candleStick.Close > candleStick.Open && candleStick.Open < t.lastValue {
		if candleStick.Open > t.lastLow {
			t.higherLow = true
		} else {
			t.higherLow = false
		}
		t.lastLow = candleStick.Open
	}
	t.lastValue = candleStick.Open
	if ! t.lowerHigh && t.higherLow {
		t.uptrend = true
		t.downtrend = false
		t.uptrendCount++
		t.downtrendCount = 0
	}
	if t.lowerHigh && ! t.higherLow {
		t.uptrend = false
		t.downtrend = true
		t.downtrendCount++
		t.uptrendCount = 0
	}
}

type averageStrategy struct {
	length int
	name string
	trend *trend
	result float64
	buy, sell bool
	results []float64
}

func NewAverageStrategy(name string, length int) domain.Strategy {
	return &averageStrategy{
		name:    name,
		length:  length,
		trend: &trend{},
	}
}

func (as *averageStrategy) String() string {
	return as.name
}

func (as *averageStrategy) GetInfo() (float64, bool, bool){
	return as.result, as.buy, as.sell
}

func (as *averageStrategy) buyCondition(value float64) {
	as.buy = as.result < value * 0.98 && as.trend.downtrendCount > 3
}

func (as *averageStrategy) sellCondition(value float64) {
	as.sell = as.result > value * 1.02 && as.trend.uptrendCount > 3
}

func (as *averageStrategy) Append(candleStick domain.CandleStick) {
	var result float64
	length := len(as.results)
	if length == 0 {
		result = candleStick.Close
	} else {
		result = candleStick.Close + as.results[length - 1]
	}
	if length == as.length {
		as.results = append(as.results[1:], result)
	} else {
		as.results = append(as.results, result)
	}
	as.result = result / float64(len(as.results))
	as.trend.evaluate(candleStick)
	as.buyCondition(candleStick.Close)
	as.sellCondition(candleStick.Close)
}
