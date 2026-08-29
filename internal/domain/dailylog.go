package domain

import (
	"fmt"
	"time"
)

// DailyLog represents a single day's habit tracking entry
type DailyLog struct {
	Date time.Time

	// Habits
	Training   bool
	Reading    bool
	Coding     bool
	Meditation bool
	NoSmoking  bool
	Stretching bool
	SleepHours float64

	// Meals (free text, any of these may be empty)
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    string

	// Day evaluation
	DayWellSpent  bool
	WhatIDidToday string
	WhatWentWell  string
	WhatToImprove string

	//Free-form journaling
	QuickNotes string
}

type NewDailyLogInput struct {
	Date time.Time

	Training   bool
	Reading    bool
	Coding     bool
	Meditation bool
	NoSmoking  bool
	Stretching bool
	SleepHours float64

	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    string

	DayWellSpent  bool
	WhatIDidToday string
	WhatWentWell  string
	WhatToImprove string

	QuickNotes string
}

func NewDailyLog(input NewDailyLogInput) (*DailyLog, error) {
	if input.SleepHours < 0 || input.SleepHours > 24 {
		return nil, fmt.Errorf("sleep hours must be between 0 and 24, got %.2f", input.SleepHours)
	}

	return &DailyLog{
		Date:          input.Date,
		Training:      input.Training,
		Reading:       input.Reading,
		Coding:        input.Coding,
		Meditation:    input.Meditation,
		NoSmoking:     input.NoSmoking,
		Stretching:    input.Stretching,
		SleepHours:    roundToTwoDecimals(input.SleepHours),
		Breakfast:     input.Breakfast,
		Lunch:         input.Lunch,
		Dinner:        input.Dinner,
		Snacks:        input.Snacks,
		DayWellSpent:  input.DayWellSpent,
		WhatIDidToday: input.WhatIDidToday,
		WhatWentWell:  input.WhatWentWell,
		WhatToImprove: input.WhatToImprove,
		QuickNotes:    input.QuickNotes,
	}, nil
}

func roundToTwoDecimals(hours float64) float64 {
	return float64(int(hours*100+0.5)) / 100
}
