package domain

import (
	"fmt"
	"strings"
)

// Kind is the structural type of a training session. It decides which detail a session log carries and which wizard strategy collects it.
type Kind string

const (
	KindStrength Kind = "strength"
	KindCardio   Kind = "cardio"
	KindSport    Kind = "sport"
)

func (k Kind) Valid() bool {
	switch k {
	case KindStrength, KindCardio, KindSport:
		return true
	default:
		return false
	}
}

// Exercise is a precribed exercise inside a session template
type Exercise struct {
	Name       string
	Block      string
	Sets       int
	RepsMin    int
	RepsMax    int
	BodyWeight bool
}

// SessionTemplate is one of the sessions a plan makes available
type SessionTemplate struct {
	Name      string
	Kind      Kind
	Modality  string
	Exercises []Exercise
}

// NewSessionTemplate validates a session template. Only strength sessions
// prescribe exercises: cardio activity is chosen when logging, and sport
// sessions are not prescribed at all.
func NewSessionTemplate(name string, kind Kind, modality string, exercises []Exercise) (SessionTemplate, error) {
	if strings.TrimSpace(name) == "" {
		return SessionTemplate{}, fmt.Errorf("session name is required")
	}
	if !kind.Valid() {
		return SessionTemplate{}, fmt.Errorf("session %q: unknown kind %q", name, kind)
	}

	if kind == KindStrength {
		if len(exercises) == 0 {
			return SessionTemplate{}, fmt.Errorf("strength session %q has no exercises", name)
		}
	} else if len(exercises) > 0 {
		return SessionTemplate{}, fmt.Errorf("session %q is %s and cannot prescribe exercises", name, kind)
	}

	seen := make(map[string]struct{}, len(exercises))
	for i, e := range exercises {
		if err := validateExercise(e); err != nil {
			return SessionTemplate{}, fmt.Errorf("session %q, exercise %d: %w", name, i+1, err)
		}
		key := normalize(e.Name)
		if _, dup := seen[key]; dup {
			return SessionTemplate{}, fmt.Errorf("session %q: duplicate exercise name %q", name, e.Name)
		}
		seen[key] = struct{}{}
	}

	return SessionTemplate{
		Name:      strings.TrimSpace(name),
		Kind:      kind,
		Modality:  strings.TrimSpace(modality),
		Exercises: exercises,
	}, nil
}

func validateExercise(e Exercise) error {
	if strings.TrimSpace(e.Name) == "" {
		return fmt.Errorf("exercise name is required")
	}
	if e.Sets < 1 {
		return fmt.Errorf("exercise %q: sets must be at least 1, got %d", e.Name, e.Sets)
	}
	if e.RepsMin < 1 {
		return fmt.Errorf("exercise %q: reps_min must be at least 1, got %d", e.Name, e.RepsMin)
	}
	if e.RepsMax < e.RepsMin {
		return fmt.Errorf("exercise %q: reps_max (%d) is below reps_min (%d)", e.Name, e.RepsMax, e.RepsMin)
	}
	return nil
}
