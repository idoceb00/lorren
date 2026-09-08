package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/idoceb00/lorren/internal/domain"
)

func RecordDay(interviewer domain.Interviewer, repo domain.Repository, date time.Time) error {
	existing, err := repo.FindByDate(date)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("checking for existing daily log: %w", err)
		}
		existing = nil
	}

	log, err := interviewer.AskDailyLog(existing)
	if err != nil {
		return fmt.Errorf("running daily log wizard: %w", err)
	}

	if err := repo.SaveDailyLog(log); err != nil {
		return fmt.Errorf("saving daily log: %w", err)
	}

	return nil
}
