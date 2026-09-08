package domain

import "time"

type Interviewer interface {
	AskDailyLog(existing *DailyLog) (*DailyLog, error)
}

type Repository interface {
	SaveDailyLog(log *DailyLog) error
	FindByDate(date time.Time) (*DailyLog, error)
}
