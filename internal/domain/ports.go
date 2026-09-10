package domain

import "time"

type DailyInterviewer interface {
	AskDailyLog(existing *DailyLog) (*DailyLog, error)
}

type TrainingInterviewer interface {
	AskTrainingSession(plan Plan) (*TrainingSession, error)
}

type Repository interface {
	SaveDailyLog(log *DailyLog) error
	FindByDate(date time.Time) (*DailyLog, error)
}

// PlanReader gives read-only access to the training plan templates the user maintains by hand
type PlanReader interface {
	LoadPlan(id string) (Plan, error)
	ListPlans() ([]string, error)
}
