package model

import (
	"os"
	"time"
)

//StartTime When the Contest Start
func StartTime() time.Time {
	startTime, _ := time.Parse("2006-01-02 15:04:05 MST", os.Getenv("START_TIME"))
	return startTime
}
