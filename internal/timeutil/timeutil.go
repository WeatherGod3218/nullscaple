package timeutil

import (
	"time"

	"github.com/WeatherGod3218/nullscaple/internal/logging"
	"github.com/sirupsen/logrus"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("America/New_York")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "DatabaseCookie"}).Fatal("Failed to load timezone!")
	}
}

func GetFormattedTime() string {
	return time.Now().In(loc).Format("2006-01-02")
}
