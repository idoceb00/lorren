package domain

import (
	"fmt"
	"time"
)

// DailyLog represents a single day's habit tracking entry
type DailyLog struct {
	Date       time.Time
	Reading    bool
	Coding     bool
	Meditation bool
	NoSmoking  bool
	Stretching bool
	SleepHours float64
}

func NewDailyLog(date time.Time, reading, coding, meditation, noSmoking, stretching bool, sleepHours float64) (*DailyLog, error) {
	if sleepHours < 0 || sleepHours > 24 {
		return nil, fmt.Errorf("sleep hours must be between 0 and 24, got %.2f", sleepHours)
	}

	return &DailyLog{
		Date:       date,
		Reading:    reading,
		Coding:     coding,
		Meditation: meditation,
		NoSmoking:  noSmoking,
		Stretching: stretching,
		SleepHours: roundToTwoDecimals(sleepHours),
	}, nil
}

func roundToTwoDecimals(hours float64) float64 {
	return float64(int(hours*100+0.5)) / 100
}
