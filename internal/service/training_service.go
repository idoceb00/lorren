package service

import (
	"github.com/idoceb00/lorren/internal/domain"
)

// RecordTraining runs the training wizard over the given plan and saves the result, returning the path of the file written.
func RecordTraining(i domain.TrainingInterviewer, r domain.TrainingRepository, plan domain.Plan) (string, error) {
	session, err := i.AskTrainingLog(plan)
	if err != nil {
		return "", err
	}

	return r.SaveTrainingLog(session)
}
