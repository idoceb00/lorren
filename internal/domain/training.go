package domain

import (
	"fmt"
	"strings"
	"time"
)

// SessionDetail is the kind-specific part of a training session. Each Kind has
// exactly one implementation.
type SessionDetail interface {
	// Kind reports which session kind this detail belongs to, so callers can
	// check it matches the session it is attached to.
	Kind() Kind

	validate() error
}

// PerformedExercise is one exercise as actually performed. The prescription is
// copied from the template so a log stays readable years later, even if the
// plan has changed since.
type PerformedExercise struct {
	Name    string
	Block   string
	Sets    int
	RepsMin int
	RepsMax int

	// BodyWeight mirrors the template: true when no load is recorded.
	BodyWeight bool

	// Done is false when the exercise was skipped. Bodyweight exercises are
	// logged with this alone; loaded ones are skipped by leaving Weight unset.
	Done bool

	// Weight in kilograms, nil for bodyweight or skipped exercises.
	Weight *float64

	// Reps actually performed, a single value for the whole exercise.
	Reps int
}

// StrengthDetail is the detail of a gym session.
type StrengthDetail struct {
	Exercises []PerformedExercise
}

func (StrengthDetail) Kind() Kind { return KindStrength }

func (d StrengthDetail) validate() error {
	if len(d.Exercises) == 0 {
		return fmt.Errorf("strength session has no exercises")
	}
	for i, e := range d.Exercises {
		if strings.TrimSpace(e.Name) == "" {
			return fmt.Errorf("exercise %d: name is required", i+1)
		}
		if !e.Done {
			continue
		}
		if e.Reps < 1 {
			return fmt.Errorf("exercise %q: reps must be at least 1, got %d", e.Name, e.Reps)
		}
		if e.BodyWeight && e.Weight != nil {
			return fmt.Errorf("exercise %q: bodyweight exercises carry no weight", e.Name)
		}
		if !e.BodyWeight {
			if e.Weight == nil {
				return fmt.Errorf("exercise %q: weight is required", e.Name)
			}
			if *e.Weight <= 0 {
				return fmt.Errorf("exercise %q: weight must be positive, got %g", e.Name, *e.Weight)
			}
		}
	}
	return nil
}

// CardioDetail is the detail of a cardio session. Activity is the medium used
// that day (running, bike, rower); the intended effort lives in the session
// modality.
type CardioDetail struct {
	Activity string
}

func (CardioDetail) Kind() Kind { return KindCardio }

func (d CardioDetail) validate() error {
	if strings.TrimSpace(d.Activity) == "" {
		return fmt.Errorf("cardio session has no activity")
	}
	return nil
}

// SportDetail is the detail of a sport session. It carries nothing of its own
// yet: the session name, modality and duration say enough.
type SportDetail struct{}

func (SportDetail) Kind() Kind { return KindSport }

func (SportDetail) validate() error { return nil }

// TrainingSession is one logged training session.
type TrainingSession struct {
	Date        time.Time
	SessionName string
	Modality    string
	Duration    time.Duration
	Notes       string
	Detail      SessionDetail
}

type NewTrainingSessionInput struct {
	Date        time.Time
	SessionName string
	Modality    string
	Duration    time.Duration
	Notes       string
	Detail      SessionDetail
}

func NewTrainingSession(in NewTrainingSessionInput) (*TrainingSession, error) {
	if in.Date.IsZero() {
		return nil, fmt.Errorf("date is required")
	}
	if strings.TrimSpace(in.SessionName) == "" {
		return nil, fmt.Errorf("session name is required")
	}
	if in.Detail == nil {
		return nil, fmt.Errorf("session %q: detail is required", in.SessionName)
	}
	if !in.Detail.Kind().Valid() {
		return nil, fmt.Errorf("session %q: unknown kind %q", in.SessionName, in.Detail.Kind())
	}
	if in.Duration < 0 {
		return nil, fmt.Errorf("session %q: duration cannot be negative", in.SessionName)
	}
	if err := in.Detail.validate(); err != nil {
		return nil, fmt.Errorf("session %q: %w", in.SessionName, err)
	}

	return &TrainingSession{
		Date:        in.Date,
		SessionName: strings.TrimSpace(in.SessionName),
		Modality:    strings.TrimSpace(in.Modality),
		Duration:    in.Duration,
		Notes:       in.Notes,
		Detail:      in.Detail,
	}, nil
}

// Kind is the kind of the session, taken from its detail. There is no separate
// field, so the two can never disagree.
func (s *TrainingSession) Kind() Kind { return s.Detail.Kind() }
