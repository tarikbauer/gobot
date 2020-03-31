package application

type data struct {
	ID string
	Low float64
	Open float64
	Close float64
	High float64
	Results []float64
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
                ['Date',     'CandleStick', '', '', '', 'Average'],
				{{range .}}
				['',  {{.Low}}, {{.Open}}, {{.Close}}, {{.High}}, {{range .Results}} {{.}}, {{ end }}],
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
</html>`
