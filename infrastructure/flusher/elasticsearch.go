package flusher

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tarikbauer/gobot/domain"
)

type header struct{
	ID string `json:"_id"`
	Index string `json:"_index"`
}

type index struct{
	Index header `json:"index"`
}

type elasticsearchFlusher struct {
	url string
	client http.Client
}

func NewElasticsearchFlusher(host string, timeout time.Duration) domain.Flusher {
	return &elasticsearchFlusher{url: host, client: http.Client{Timeout: timeout}}
}

func (ef *elasticsearchFlusher) flush(data []byte) error {
	var operation domain.Operation = "infrastructure.flusher.flush"
	response, err := ef.client.Post(ef.url + "/_bulk", "application/x-ndjson", bytes.NewBuffer(data))
	if err != nil {
		return domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	}
	content, _ := ioutil.ReadAll(response.Body)
	switch response.StatusCode {
	case 200:
		return nil
	default:
		return domain.NewError(operation, domain.Internal, logrus.ErrorLevel, errors.New(string(content)))
	}
}

func (ef *elasticsearchFlusher) appendData(_index string, data *[]byte, object domain.Data) {
	header := index{header{
		Index: _index,
		ID:    object.GetID(),
	}}
	content, _ := json.Marshal(object)
	headerContent, _ := json.Marshal(header)
	*data = append(*data, headerContent...)
	*data = append(*data, []byte("\n")...)
	*data = append(*data, content...)
	*data = append(*data, []byte("\n")...)
}

func (ef *elasticsearchFlusher) FlushStrategies(strategies []domain.StrategyData) error {
	var data []byte
	var operation domain.Operation = "infrastructure.flusher.FlushStrategies"
	for _, object := range strategies {
		ef.appendData("strategies", &data, object)
	}
	err := ef.flush(data)
	if err != nil {
		return domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	}
	return nil
}

func (ef *elasticsearchFlusher) FlushCandleSticks(candleSticks []domain.CandleStickData) error {
	var data []byte
	var operation domain.Operation = "infrastructure.flusher.FlushCandleSticks"
	for _, object := range candleSticks {
		ef.appendData("candlesticks", &data, object)
	}
	err := ef.flush(data)
	if err != nil {
		return domain.NewError(operation, domain.Unavailable, logrus.ErrorLevel, err)
	}
	return nil
}

func (ef *elasticsearchFlusher) Flush(strategies []domain.StrategyData, candleSticks []domain.CandleStickData) error {
	return nil
}
