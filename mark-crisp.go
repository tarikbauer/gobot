package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CandleStick struct {
	Open float64
	High float64
	Low float64
	Close float64
}

type status string

const (
	Wait status = "wait"
	Buy status = "buy"
	Sell status = "sell"
	Finish status = "finish"
	Cancel status = "cancel"
)

type state struct {
	status status
	value float64
	values [3]float64
}

type buy state

type sell state

type stock struct {
	buy buy
	sell sell
	uptrend bool
	downtrend bool
	upThreshHold float64
	downThreshHold float64
	values []CandleStick
}

func (s *stock) Append(stick CandleStick) {
	s.values = append(s.values, stick)
	if len(s.values) == 1 {
		return
	}
	lastStick := s.values[len(s.values)-2]
	if stick.Close > lastStick.Close {
		if s.downtrend {
			if s.buy.values[0] == -1 {
				s.buy.values[0] = lastStick.Low
			} else if s.buy.values[1] != -1 && s.buy.values[2] == -1 {
				s.buy.values[2] = lastStick.Low
			}
			if stick.High < s.buy.value {
				s.buy.value = stick.High
			}
		}
		s.uptrend = true
		s.downtrend = false
	} else if stick.Close < lastStick.Close {
		if s.uptrend {
			if s.buy.values[0] != -1 && s.buy.values[1] == -1 {
				s.buy.status = Buy
				s.buy.value = lastStick.High
				s.buy.values[1] = lastStick.High
			}
		}
		s.uptrend = false
		s.downtrend = true
	}
	if s.buy.status != Sell {
		if stick.Low < s.buy.values[0] {
			s.reset(Cancel)
		}
		if s.buy.values[2] != -1 && s.buy.value < stick.High {
			s.buy.status = Sell
			s.downThreshHold = s.buy.values[2]
			s.upThreshHold = s.buy.value*(s.buy.values[1]/s.buy.values[0])
		}
	}
	if s.buy.status == Sell && stick.High > s.upThreshHold || stick.Low < s.downThreshHold {
		s.reset(Finish)
	}
}

func (s *stock) reset(status status) {
	s.buy.value = -1
	s.upThreshHold = -1
	s.downThreshHold = -1
	s.buy.status = status
	s.buy.values = [3]float64{-1, -1, -1}
}

func NewStock() *stock {
	return &stock{
		uptrend: false,
		downtrend: false,
		upThreshHold: -1,
		downThreshHold: -1,
		buy: buy{Wait, -1,  [3]float64{-1, -1, -1}},
		sell: sell{Wait, -1,  [3]float64{-1, -1, -1}},
	}
}

var currentStock *stock

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("<html><h1>Page Not Found</h1></html>"))
}

func notAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("<html><h1>Method Not Allowed</h1></html>"))
}

func get(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.New("CandleStick Chart")
	tmpl, _ = tmpl.Parse(html)
	buf := bytes.NewBuffer([]byte{})
	asd := []CandleStick{{1, 2, 0.5, 1.5}, {1.5, 2.1, 0.7, 1.8}}
	for _, zxc := range asd {
		currentStock.Append(zxc)
	}
	_ = tmpl.Execute(buf, currentStock.values)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	currentStock = NewStock()
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").HandlerFunc(get).Methods("GET")
	router.Path("/").HandlerFunc(post).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	err := http.ListenAndServe(fmt.Sprint(":", 8088), router)
	if err != nil {
		log.Fatal(err)
	}
}

var html = `<html>
<head>
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
        google.charts.load('current', {'packages':['corechart']});
        google.charts.setOnLoadCallback(drawVisualization);
        function drawVisualization() {
            // Some raw data (not necessarily accurate)
            var data = google.visualization.arrayToDataTable([
                ['Date',     'CandleStick', '', '', ''],
				{{range .}}
				['',  {{.Low}}, {{.Open}}, {{.Close}}, {{.High}}],
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
            };
            var chart = new google.visualization.ComboChart(document.getElementById('chart_div'));
            chart.draw(data, options);
        }
    </script>
</head>
<body>
<div id="chart_div" style="height: 900px;"></div>
</body>
</html>`
