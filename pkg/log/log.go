package log

import (
	"os"
	"time"

	"charm.land/log/v2"
)

var Logger *log.Logger

func LogSetting() {
	Logger = log.New(os.Stderr) //DateTime
	Logger.SetTimeFormat(time.DateTime)
}
