package domain

import (
	"fmt"
	"strings"
)

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
	for _, s := range sessions {
		key := normalize(s.Name)
		if _, dup := seen[key]; dup {
			return Plan{}, fmt.Errorf("plan %q: duplicate session name %q", id, s.Name)
		}
		seen[key] = struct{}{}
	}

	return Plan{id: id, name: name, sessions: sessions}, nil
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
