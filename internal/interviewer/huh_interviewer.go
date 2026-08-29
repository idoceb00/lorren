package interviewer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/idoceb00/lorren/internal/domain"
)

type HuhInterviewer struct{}

func NewHuhInterviewer() *HuhInterviewer {
	return &HuhInterviewer{}
}

func (h *HuhInterviewer) AskDailyLog() (*domain.DailyLog, error) {
	var (
		reading, coding, meditation, noSmoking, stretching bool
		sleepHoursStr                                      string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Did you read today?").Value(&reading),
			huh.NewConfirm().Title("Did you code today?").Value(&coding),
			huh.NewConfirm().Title("Did you meditate today?").Value(&meditation),
			huh.NewConfirm().Title("Did you avoid smoking today?").Value(&noSmoking),
			huh.NewConfirm().Title("Did you stretch today?").Value(&stretching),
			huh.NewInput().
				Title("How many hours did you sleep?").
				Value(&sleepHoursStr),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	sleepHours, err := parseSleepHours(sleepHoursStr)
	if err != nil {
		return nil, err
	}

	return domain.NewDailyLog(time.Now(), reading, coding, meditation, noSmoking, stretching, sleepHours)
}

func parseSleepHours(s string) (float64, error) {
	var hours float64
	_, err := fmt.Sscanf(s, "%f", &hours)
	if err != nil {
		return 0, fmt.Errorf("invalid sleep hours %q: %w", s, err)
	}

	return hours, nil
}
