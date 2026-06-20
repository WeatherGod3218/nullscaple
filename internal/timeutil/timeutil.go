package timeutil

import (
	"time"
)

func GetFormattedTime() string {
	return time.Now().UTC().Format("2006-01-02")
}
