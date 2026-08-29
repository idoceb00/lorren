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
		training, reading, coding, meditation, noSmoking, stretching bool
		sleepHoursStr                                                string

		breakfast, lunch, dinner, snacks string

		dayWellSpent                               bool
		whatIDidToday, whatWentWell, whatToImprove string
		quickNotes                                 string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Did you train today?").Value(&training),
			huh.NewConfirm().Title("Did you read today?").Value(&reading),
			huh.NewConfirm().Title("Did you code today?").Value(&coding),
			huh.NewConfirm().Title("Did you meditate today?").Value(&meditation),
			huh.NewConfirm().Title("Did you avoid smoking today?").Value(&noSmoking),
			huh.NewConfirm().Title("Did you stretch today?").Value(&stretching),
			huh.NewInput().Title("How many hours did you sleep?").Value(&sleepHoursStr),
		),
		huh.NewGroup(
			huh.NewText().Title("Breakfast").Value(&breakfast),
			huh.NewText().Title("Lunch").Value(&lunch),
			huh.NewText().Title("Dinner").Value(&dinner),
			huh.NewText().Title("Snacks").Value(&snacks),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("Was today well spent?").Value(&dayWellSpent),
			huh.NewText().Title("What did you do today?").Value(&whatIDidToday),
			huh.NewText().Title("What went well?").Value(&whatWentWell),
			huh.NewText().Title("What could you improve?").Value(&whatToImprove),
			huh.NewText().Title("Quick notes").Value(&quickNotes),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	sleepHours, err := parseSleepHours(sleepHoursStr)
	if err != nil {
		return nil, err
	}

	return domain.NewDailyLog(domain.NewDailyLogInput{
		Date:          time.Now(),
		Training:      training,
		Reading:       reading,
		Coding:        coding,
		Meditation:    meditation,
		NoSmoking:     noSmoking,
		Stretching:    stretching,
		SleepHours:    sleepHours,
		Breakfast:     breakfast,
		Lunch:         lunch,
		Dinner:        dinner,
		Snacks:        snacks,
		DayWellSpent:  dayWellSpent,
		WhatIDidToday: whatIDidToday,
		WhatWentWell:  whatWentWell,
		WhatToImprove: whatToImprove,
		QuickNotes:    quickNotes,
	})
}

func parseSleepHours(s string) (float64, error) {
	var hours float64
	_, err := fmt.Sscanf(s, "%f", &hours)
	if err != nil {
		return 0, fmt.Errorf("invalid sleep hours %q: %w", s, err)
	}

	return hours, nil
}
