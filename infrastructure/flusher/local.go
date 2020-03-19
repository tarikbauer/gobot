package flusher

import (
	"bytes"
	"github.com/tarikbauer/gobot/domain"
	"html/template"
	"io/ioutil"
)

type localFlusher struct {
	path string
}

type data struct {
	ID string
	Low float64
	Open float64
	Close float64
	High float64
	Result float64
}

func Render() {
	tmpl := template.New("CandleStick Chart")
	tmpl, _ = tmpl.Parse(`<html>
<head>
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
        google.charts.load('current', {'packages':['corechart']});
        google.charts.setOnLoadCallback(drawVisualization);

        function drawVisualization() {
            // Some raw data (not necessarily accurate)
            var data = google.visualization.arrayToDataTable([
                ['Date',     'CandleStick', '', '', '', 'Average'],
				{{range .}}
				['',  {{.Low}}, {{.Open}}, {{.Close}}, {{.High}}, {{.Result}}],
				{{ end }}
                
            ]);

            var options = {
                title : 'CandleSticks',
                animation: {startup: true},
                candlestick: {
                    risingColor: { strokeWidth: 0, fill: '#00e676'},
                    fallingColor: { strokeWidth: 0, fill: '#ff1744'}
                },
                vAxis: {title: 'Values'},
                hAxis: {title: 'Date'},
                seriesType: 'candlesticks',
                series: {
                    1: {
                        type: 'line',
                        curveType: 'function'
                    }
                }
            };

            var chart = new google.visualization.ComboChart(document.getElementById('chart_div'));
            chart.draw(data, options);
        }
    </script>
</head>
<body>
<div id="chart_div" style="height: 900px;"></div>
</body>
</html>`)
	buf := bytes.NewBuffer([]byte{})
	tmpl.Execute(buf, contents)
	ioutil.WriteFile("/Users/tarikbauer/go/src/github.com/tarikbauer/gobot/asd.html", buf.Bytes(), 0777)
}

func NewLocalFlusher(path string) domain.Flusher {
	return &localFlusher{path}
}

func (lf *localFlusher) FlushStrategies(strategies []domain.StrategyData) error {
	return nil
}

var contents []data

func (lf *localFlusher) FlushCandleSticks(candleSticks []domain.CandleStickData) error {
	return nil
}

func(lf *localFlusher) Flush(strategies []domain.StrategyData, candleSticks []domain.CandleStickData) error {
	for index, object := range strategies {
		contents = append(contents, data{
			ID:     object.Date,
			Low:    candleSticks[index].Low,
			Open:   candleSticks[index].Open,
			Close:  candleSticks[index].Close,
			High:   candleSticks[index].High,
			Result: object.Result,
		})
	}
	return nil
}

