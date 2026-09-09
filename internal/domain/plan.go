package domain

import (
	"fmt"
	"strings"
)

// Kind is the structural type of a training session. It decides which detail a
// session log carries and which wizard strategy collects it.
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

// Plan is the set of sessions currently available to log.
type Plan struct {
	id       string
	name     string
	sessions []SessionTemplate
}

func NewPlan(id, name string, sessions []SessionTemplate) (Plan, error) {
	if strings.TrimSpace(id) == "" {
		return Plan{}, fmt.Errorf("plan id is required")
	}
	if strings.TrimSpace(name) == "" {
		return Plan{}, fmt.Errorf("plan %q: name is required", id)
	}
	if len(sessions) == 0 {
		return Plan{}, fmt.Errorf("plan %q has no sessions", id)
	}

	seen := make(map[string]struct{}, len(sessions))
	for i, s := range sessions {
		if err := validateSession(s); err != nil {
			return Plan{}, fmt.Errorf("plan %q, session %d: %w", id, i+1, err)
		}
		key := normalize(s.Name)
		if _, dup := seen[key]; dup {
			return Plan{}, fmt.Errorf("plan %q: duplicate session name %q", id, s.Name)
		}
		seen[key] = struct{}{}
	}

	return Plan{id: id, name: name, sessions: sessions}, nil
}

func validateSession(s SessionTemplate) error {
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("session name is required")
	}
	if !s.Kind.Valid() {
		return fmt.Errorf("session %q: unknown kind %q", s.Name, s.Kind)
	}

	if s.Kind == KindStrength {
		if len(s.Exercises) == 0 {
			return fmt.Errorf("strength session %q has no exercises", s.Name)
		}
	} else if len(s.Exercises) > 0 {
		return fmt.Errorf("session %q is %s and cannot prescribe exercises", s.Name, s.Kind)
	}

	seen := make(map[string]struct{}, len(s.Exercises))
	for i, e := range s.Exercises {
		if err := validateExercise(e); err != nil {
			return fmt.Errorf("session %q, exercise %d: %w", s.Name, i+1, err)
		}
		key := normalize(e.Name)
		if _, dup := seen[key]; dup {
			return fmt.Errorf("session %q: duplicate exercise name %q", s.Name, e.Name)
		}
		seen[key] = struct{}{}
	}

	return nil
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

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (p Plan) ID() string   { return p.id }
func (p Plan) Name() string { return p.name }

// Sessions returns a copy so callers cannot mutate the validated plan.
func (p Plan) Sessions() []SessionTemplate {
	out := make([]SessionTemplate, len(p.sessions))
	copy(out, p.sessions)
	return out
}

// Returns the session names in plan order, ready to feed a picker.
func (p Plan) SessionNames() []string {
	names := make([]string, 0, len(p.sessions))
	for _, s := range p.sessions {
		names = append(names, s.Name)
	}
	return names
}

// Session looks up a session by name, case insensitively
func (p Plan) Session(name string) (SessionTemplate, error) {
	target := normalize(name)
	for _, s := range p.sessions {
		if normalize(s.Name) == target {
			return s, nil
		}
	}
	return SessionTemplate{}, fmt.Errorf("session %q: %w", name, ErrNotFound)
}
