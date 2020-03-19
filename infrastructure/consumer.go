package infrastructure

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tarikbauer/gobot/domain"
)

type binanceConsumer struct {
	url, symbol string
	client http.Client
	interval domain.Interval
}

func NewBinanceConsumer(symbol string, interval domain.Interval, timeout time.Duration) domain.Consumer {
	return &binanceConsumer{
		symbol: symbol,
		interval: interval,
		client: http.Client{Timeout: timeout},
		url: "https://api.binance.com/api/v3/klines?limit=50&symbol=" + symbol + "&interval=" + interval.String(),
	}
}

func (bc *binanceConsumer) GetSymbol() string {
	return bc.symbol
}

func (bc *binanceConsumer) GetInterval() domain.Interval {
	return bc.interval
}

func (bc *binanceConsumer) parseResponse(args []interface{}) (*domain.CandleStick, error) {
	var operation domain.Operation = "infrastructure.consumer.parseResponse"
	candleStick := domain.CandleStick{}
	for index, arg := range args {
		switch index {
		case 0:
			castedArg, ok := arg.(float64)
			if !ok {
				return nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, errors.New("error while casting argument"))
			}
			candleStick.OpenTime = time.Unix(int64(castedArg/1000), 0)
		case 1:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
			}
			candleStick.Open = value
		case 2:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
			}
			candleStick.High = value
		case 3:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
			}
			candleStick.Low = value
		case 4:
			value, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				return nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
			}
			candleStick.Close = value
		}
	}
	return &candleStick, nil
}

func (bc *binanceConsumer) GetData() (time.Duration, []domain.CandleStick, error) {
	var operation domain.Operation = "infrastructure.consumer.GetData"
	response, err := bc.client.Get(bc.url)
	if err != nil {
		return 0, nil, domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	}
	var v [][]interface{}
	var candleSticks []domain.CandleStick
	content, _ := ioutil.ReadAll(response.Body)
	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(content, &v)
		if err != nil {
			return 0, nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
		}
		for _, candlestick := range v {
			if len(candlestick) != 12 {
				return 0, nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, errors.New("not expected binance response"))
			}
			newCandleStick, err := bc.parseResponse(candlestick[:5])
			if err != nil {
				return 0, nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, err)
			}
			candleSticks = append(candleSticks, *newCandleStick)
		}
		return 50 * bc.interval.Duration(), candleSticks, nil
	default:
		return 0, nil, domain.NewError(operation, domain.Internal, logrus.ErrorLevel, errors.New(string(content)))
	}
}
