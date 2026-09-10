package service

import (
	"github.com/idoceb00/lorren/internal/domain"
)

// TrainingRepository writes training session logs and returns the path written.
type TrainingRepository interface {
	SaveTrainingLog(s *domain.TrainingLog) (string, error)
}

// RecordTraining runs the training wizard over the given plan and saves the result, returning the path of the file written.
func RecordTraining(i domain.TrainingInterviewer, r TrainingRepository, plan domain.Plan) (string, error) {
	session, err := i.AskTrainingLog(plan)
	if err != nil {
		return "", err
	}

	return r.SaveTrainingLog(session)
}
