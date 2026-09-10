package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/idoceb00/lorren/internal/domain"
)

func RecordDay(interviewer domain.DailyInterviewer, repo domain.DailyRepository, date time.Time) (string, error) {
	existing, err := repo.FindByDate(date)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return "", fmt.Errorf("checking for existing daily log: %w", err)
		}
		existing = nil
	}

	log, err := interviewer.AskDailyLog(existing)
	if err != nil {
		return "", fmt.Errorf("running daily log wizard: %w", err)
	}

	path, err := repo.SaveDailyLog(log)
	if err != nil {
		return "", fmt.Errorf("saving daily log: %w", err)
	}

	return path, nil
}
