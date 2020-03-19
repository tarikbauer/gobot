package domain

import "time"

type Interval string

var (
	MINUTE Interval = "1m"
	FIVE_MINUTES Interval = "5m"
	HALF_HOUR Interval = "30m"
	HOUR Interval = "1h"
	FOUR_HOURS Interval = "4h"
	HALF_DAY Interval = "12h"
	DAY Interval = "1d"
	WEEK Interval = "1w"
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
	Low float64
	High float64
	Open float64
	Close float64
	OpenTime time.Time
}
