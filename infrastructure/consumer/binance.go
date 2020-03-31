package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/tarikbauer/gobot/domain"
)

type binanceConsumer struct {
	limit int
	url, symbol string
	client *http.Client
	interval domain.Interval
}

func NewBinanceConsumer(symbol string, interval domain.Interval, limit int, client *http.Client) domain.Consumer {
	return &binanceConsumer{
		limit: limit,
		client: client,
		symbol: symbol,
		interval: interval,
		url: fmt.Sprint("https://api.binance.com/api/v3/klines?limit=", limit, "&symbol=", symbol, "&interval=", interval.String()),
	}
}

func (bc *binanceConsumer) GetLimit() int {
	return bc.limit
}

func (bc *binanceConsumer) GetSymbol() string {
	return bc.symbol
}

func (bc *binanceConsumer) GetInterval() domain.Interval {
	return bc.interval
}

func (bc *binanceConsumer) parseResponse(args []interface{}) (*domain.CandleStick, error) {
	candleStick := domain.CandleStick{}
	for index, arg := range args {
		switch index {
		case 0:
			castedArg, ok := arg.(float64)
			if !ok {
				return nil, domain.UnexpectedResponseFormat
			}
			candleStick.OpenTime = time.Unix(int64(castedArg/1000), 0)
		case 1:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.UnexpectedResponseFormat
			}
			candleStick.Open = value
		case 2:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.UnexpectedResponseFormat
			}
			candleStick.High = value
		case 3:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.UnexpectedResponseFormat
			}
			candleStick.Low = value
		case 4:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.UnexpectedResponseFormat
			}
			candleStick.Close = value
		}
	}
	return &candleStick, nil
}

func (bc *binanceConsumer) GetData() ([]domain.CandleStick, error) {
	response, err := bc.client.Get(bc.url)
	if err != nil {
		return nil, err
	}
	var v [][]interface{}
	var candleSticks []domain.CandleStick
	content, _ := ioutil.ReadAll(response.Body)
	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(content, &v)
		if err != nil {
			return nil, domain.InvalidConsumedContent
		}
		for _, candlestick := range v {
			if len(candlestick) != 12 {
				return nil, domain.UnexpectedResponseFormat
			}
			newCandleStick, err := bc.parseResponse(candlestick[:5])
			if err != nil {
				return nil, err
			}
			candleSticks = append(candleSticks, *newCandleStick)
		}
		return candleSticks, nil
	default:
		return nil, errors.New(string(content))
	}
}
