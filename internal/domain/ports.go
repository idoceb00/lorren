package domain

type Interviewer interface {
	AskDailyLog() (*DailyLog, error)
}

type Repository interface {
	SaveDailyLog(log *DailyLog) error
}
