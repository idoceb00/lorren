package domain

import "time"

type Interviewer interface {
	AskDailyLog() (*DailyLog, error)
}

type Repository interface {
	SaveDailyLog(log *DailyLog) error
	FindByDate(date time.Time) (*DailyLog, error)
}
