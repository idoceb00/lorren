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

func (h *HuhInterviewer) AskDailyLog(existing *domain.DailyLog) (*domain.DailyLog, error) {
	state := newFormState(existing)

	if err := buildForm(state).Run(); err != nil {
		return nil, err
	}

	input, err := toDailyLogInput(state)
	if err != nil {
		return nil, err
	}

	return domain.NewDailyLog(input)
}

// formState holds every field the wizard reads from and writes to.
// huh binds each widget to a field's address, so this struct is the
// single source of truth across seeding, form input, and reconstruction.
type formState struct {
	Date time.Time

	Training, Reading, Coding, Meditation, NoSmoking, Stretching bool
	SleepHoursStr                                                string

	Breakfast, Lunch, Dinner, Snacks string

	DayWellSpent                               bool
	WhatIDidToday, WhatWentWell, WhatToImprove string
	QuickNotes                                 string
}

// newFormState seeds a formState from an existing log, or returns
// zero-value defaults (and today's date) when existing is nil.
func newFormState(existing *domain.DailyLog) *formState {
	state := &formState{Date: time.Now()}

	if existing == nil {
		return state
	}

	state.Date = existing.Date
	state.Training = existing.Training
	state.Reading = existing.Reading
	state.Coding = existing.Coding
	state.Meditation = existing.Meditation
	state.NoSmoking = existing.NoSmoking
	state.Stretching = existing.Stretching
	state.SleepHoursStr = fmt.Sprintf("%.2f", existing.SleepHours)
	state.Breakfast = existing.Breakfast
	state.Lunch = existing.Lunch
	state.Dinner = existing.Dinner
	state.Snacks = existing.Snacks
	state.DayWellSpent = existing.DayWellSpent
	state.WhatIDidToday = existing.WhatIDidToday
	state.WhatWentWell = existing.WhatWentWell
	state.WhatToImprove = existing.WhatToImprove
	state.QuickNotes = existing.QuickNotes

	return state
}

// buildForm wires the widgets to state's fields. huh reads state's
// current values as defaults and writes the user's answers back
// into the same fields on form.Run().
func buildForm(state *formState) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Did you train today?").Value(&state.Training),
			huh.NewConfirm().Title("Did you read today?").Value(&state.Reading),
			huh.NewConfirm().Title("Did you code today?").Value(&state.Coding),
			huh.NewConfirm().Title("Did you meditate today?").Value(&state.Meditation),
			huh.NewConfirm().Title("Did you avoid smoking today?").Value(&state.NoSmoking),
			huh.NewConfirm().Title("Did you stretch today?").Value(&state.Stretching),
			huh.NewInput().Title("How many hours did you sleep?").Value(&state.SleepHoursStr),
		),
		huh.NewGroup(
			huh.NewText().Title("Breakfast").Value(&state.Breakfast),
			huh.NewText().Title("Lunch").Value(&state.Lunch),
			huh.NewText().Title("Dinner").Value(&state.Dinner),
			huh.NewText().Title("Snacks").Value(&state.Snacks),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("Was today well spent?").Value(&state.DayWellSpent),
			huh.NewText().Title("What did you do today?").Value(&state.WhatIDidToday),
			huh.NewText().Title("What went well?").Value(&state.WhatWentWell),
			huh.NewText().Title("What could you improve?").Value(&state.WhatToImprove),
			huh.NewText().Title("Quick notes").Value(&state.QuickNotes),
		),
	)
}

// toDailyLogInput translates a completed formState into the domain's
// construction type, parsing the one field (sleep hours) whose wire
// format differs from the domain's.
func toDailyLogInput(state *formState) (domain.NewDailyLogInput, error) {
	sleepHours, err := parseSleepHours(state.SleepHoursStr)
	if err != nil {
		return domain.NewDailyLogInput{}, err
	}

	return domain.NewDailyLogInput{
		Date:          state.Date,
		Training:      state.Training,
		Reading:       state.Reading,
		Coding:        state.Coding,
		Meditation:    state.Meditation,
		NoSmoking:     state.NoSmoking,
		Stretching:    state.Stretching,
		SleepHours:    sleepHours,
		Breakfast:     state.Breakfast,
		Lunch:         state.Lunch,
		Dinner:        state.Dinner,
		Snacks:        state.Snacks,
		DayWellSpent:  state.DayWellSpent,
		WhatIDidToday: state.WhatIDidToday,
		WhatWentWell:  state.WhatWentWell,
		WhatToImprove: state.WhatToImprove,
		QuickNotes:    state.QuickNotes,
	}, nil
}

func parseSleepHours(s string) (float64, error) {
	var hours float64
	_, err := fmt.Sscanf(s, "%f", &hours)
	if err != nil {
		return 0, fmt.Errorf("invalid sleep hours %q: %w", s, err)
	}

	return hours, nil
}
