package domain

import "time"

type Interval string

var (
	Minute Interval = "1m"
	FiveMinutes Interval = "5m"
	HalfHour Interval = "30m"
	Hour Interval = "1h"
	FourHours Interval = "4h"
	HalfDay Interval = "12h"
	Day Interval = "1d"
	Week Interval = "1w"
)

func (i Interval) String() string {
	return string(i)
}

func (i Interval) Duration() time.Duration {
	switch i.String() {
	case "1m":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "12h":
		return 12 * time.Hour
	case "1d":
		return 24 * time.Hour
	case "1w":
		return 7 * 24 * time.Hour
	default:
		return 0
	}
}

type CandleStick struct {
	Low float64 `json:"low"`
	High float64 `json:"high"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	OpenTime time.Time `json:"open_time"`
}
