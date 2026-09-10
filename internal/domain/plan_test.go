package domain_test

import (
	"errors"
	"testing"

	"github.com/idoceb00/lorren/internal/domain"
)

func validExercise() domain.Exercise {
	return domain.Exercise{
		Name:    "Front squat",
		Block:   "fuerza",
		Sets:    3,
		RepsMin: 4,
		RepsMax: 6,
	}
}

func validStrengthSession(t *testing.T) domain.SessionTemplate {
	t.Helper()

	s, err := domain.NewSessionTemplate("Gym A", domain.KindStrength, "fuerza", []domain.Exercise{validExercise()})
	if err != nil {
		t.Fatalf("building a valid strength session: %v", err)
	}
	return s
}

func validSportSession(t *testing.T) domain.SessionTemplate {
	t.Helper()

	s, err := domain.NewSessionTemplate("Boxeo", domain.KindSport, "boxeo", nil)
	if err != nil {
		t.Fatalf("building a valid sport session: %v", err)
	}
	return s
}

func validCardioSession(t *testing.T) domain.SessionTemplate {
	t.Helper()

	s, err := domain.NewSessionTemplate("Zone 2", domain.KindCardio, "zona 2", nil)
	if err != nil {
		t.Fatalf("building a valid cardio session: %v", err)
	}
	return s
}

func TestNewPlan(t *testing.T) {
	renamed := validStrengthSession(t)
	renamed.Name = "  gym a  "

	tests := []struct {
		name     string
		id       string
		planName string
		sessions []domain.SessionTemplate
		wantErr  bool
	}{
		{
			name:     "valid plan with all kinds",
			id:       "fuerza-boxeo",
			planName: "Fuerza y boxeo",
			sessions: []domain.SessionTemplate{
				validStrengthSession(t),
				validSportSession(t),
				validCardioSession(t),
			},
		},
		{
			name:     "missing id",
			id:       "  ",
			planName: "Fuerza y boxeo",
			sessions: []domain.SessionTemplate{validStrengthSession(t)},
			wantErr:  true,
		},
		{
			name:     "missing name",
			id:       "fuerza-boxeo",
			planName: "",
			sessions: []domain.SessionTemplate{validStrengthSession(t)},
			wantErr:  true,
		},
		{
			name:     "no sessions",
			id:       "fuerza-boxeo",
			planName: "Fuerza y boxeo",
			sessions: nil,
			wantErr:  true,
		},
		{
			name:     "duplicate session names ignoring case and spaces",
			id:       "fuerza-boxeo",
			planName: "Fuerza y boxeo",
			sessions: []domain.SessionTemplate{
				validStrengthSession(t),
				renamed,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewPlan(tt.id, tt.planName, tt.sessions)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanSessionNamesKeepsPlanOrder(t *testing.T) {
	p, err := domain.NewPlan("fuerza-boxeo", "Fuerza y boxeo", []domain.SessionTemplate{
		validStrengthSession(t),
		validSportSession(t),
		validCardioSession(t),
	})
	if err != nil {
		t.Fatalf("NewPlan() unexpected error: %v", err)
	}

	want := []string{"Gym A", "Boxeo", "Zone 2"}
	got := p.SessionNames()
	if len(got) != len(want) {
		t.Fatalf("SessionNames() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SessionNames()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestPlanSession(t *testing.T) {
	p, err := domain.NewPlan("fuerza-boxeo", "Fuerza y boxeo", []domain.SessionTemplate{
		validStrengthSession(t),
		validSportSession(t),
	})
	if err != nil {
		t.Fatalf("NewPlan() unexpected error: %v", err)
	}

	t.Run("finds session ignoring case and spaces", func(t *testing.T) {
		s, err := p.Session("  gym a ")
		if err != nil {
			t.Fatalf("Session() unexpected error: %v", err)
		}
		if s.Name != "Gym A" {
			t.Errorf("Session().Name = %q, want %q", s.Name, "Gym A")
		}
	})

	t.Run("returns ErrNotFound for unknown session", func(t *testing.T) {
		_, err := p.Session("Natación")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Session() error = %v, want ErrNotFound", err)
		}
	})
}

func TestPlanSessionsIsACopy(t *testing.T) {
	p, err := domain.NewPlan("fuerza-boxeo", "Fuerza y boxeo", []domain.SessionTemplate{validStrengthSession(t)})
	if err != nil {
		t.Fatalf("NewPlan() unexpected error: %v", err)
	}

	p.Sessions()[0].Name = "mutated"

	if got := p.SessionNames()[0]; got != "Gym A" {
		t.Errorf("plan was mutated through Sessions(): got %q", got)
	}
}
