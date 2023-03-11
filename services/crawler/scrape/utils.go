package scrape

import (
	"time"
)

func timeToStr(t time.Time) string {
	return t.Format("20060102_150405")
}
