package model

import (
	"os"
	"time"
)

//StartTime When the Contest Start
func StartTime() time.Time {
	startTime, _ := time.Parse(time.RFC3339, os.Getenv("START_TIME"))
	return startTime
}

//FinishTime When the Contest Finish
func FinishTime() time.Time {
	finishTime, _ := time.Parse(time.RFC3339, os.Getenv("FINISH_TIME"))
	return finishTime
}
