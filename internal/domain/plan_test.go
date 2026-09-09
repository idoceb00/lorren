package domain

import (
	"errors"
	"testing"
)

func validExercise() Exercise {
	return Exercise{
		Name:    "Front squat",
		Block:   "fuerza",
		Sets:    3,
		RepsMin: 4,
		RepsMax: 6,
	}
}

func validStrengthSession() SessionTemplate {
	return SessionTemplate{
		Name:      "Gym A",
		Kind:      KindStrength,
		Modality:  "fuerza",
		Exercises: []Exercise{validExercise()},
	}
}

func validSportSession() SessionTemplate {
	return SessionTemplate{
		Name:     "Boxeo",
		Kind:     KindSport,
		Modality: "boxeo",
	}
}

func TestNewPlan(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		planName string
		sessions []SessionTemplate
		wantErr  bool
	}{
		{
			name:     "valid plan with all kinds",
			id:       "fuerza-boxeo",
			planName: "Fuerza y boxeo",
			sessions: []SessionTemplate{
				validStrengthSession(),
				validSportSession(),
				{Name: "Zone 2", Kind: KindCardio, Modality: "zona 2"},
			},
		},
		{
			name:     "missing id",
			id:       "  ",
			planName: "Fuerza y boxeo",
			sessions: []SessionTemplate{validStrengthSession()},
			wantErr:  true,
		},
		{
			name:     "missing name",
			id:       "fuerza-boxeo",
			planName: "",
			sessions: []SessionTemplate{validStrengthSession()},
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
			sessions: []SessionTemplate{
				validStrengthSession(),
				func() SessionTemplate {
					s := validStrengthSession()
					s.Name = "  gym a  "
					return s
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPlan(tt.id, tt.planName, tt.sessions)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// This is the validation the whole template design rests on: only strength
// sessions prescribe exercises.
func TestNewPlanValidatesByKind(t *testing.T) {
	tests := []struct {
		name    string
		session SessionTemplate
		wantErr bool
	}{
		{
			name:    "strength with exercises",
			session: validStrengthSession(),
		},
		{
			name: "strength without exercises",
			session: SessionTemplate{
				Name: "Gym A",
				Kind: KindStrength,
			},
			wantErr: true,
		},
		{
			name:    "sport without exercises",
			session: validSportSession(),
		},
		{
			name: "sport with exercises",
			session: SessionTemplate{
				Name:      "Boxeo",
				Kind:      KindSport,
				Exercises: []Exercise{validExercise()},
			},
			wantErr: true,
		},
		{
			name: "cardio with exercises",
			session: SessionTemplate{
				Name:      "Zone 2",
				Kind:      KindCardio,
				Exercises: []Exercise{validExercise()},
			},
			wantErr: true,
		},
		{
			name: "unknown kind",
			session: SessionTemplate{
				Name: "Gym A",
				Kind: Kind("gimnasio"),
			},
			wantErr: true,
		},
		{
			name: "empty kind",
			session: SessionTemplate{
				Name: "Gym A",
			},
			wantErr: true,
		},
		{
			name: "missing session name",
			session: SessionTemplate{
				Kind:      KindStrength,
				Exercises: []Exercise{validExercise()},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPlan("fuerza-boxeo", "Fuerza y boxeo", []SessionTemplate{tt.session})
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPlanValidatesExercises(t *testing.T) {
	tests := []struct {
		name      string
		exercises []Exercise
		wantErr   bool
	}{
		{
			name:      "valid exercise",
			exercises: []Exercise{validExercise()},
		},
		{
			name: "fixed reps",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.Name = "Salto vertical"
				e.RepsMin, e.RepsMax = 3, 3
				e.BodyWeight = true
				return e
			}()},
		},
		{
			name: "missing name",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.Name = " "
				return e
			}()},
			wantErr: true,
		},
		{
			name: "zero sets",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.Sets = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "zero reps_min",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.RepsMin = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "reps_max below reps_min",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.RepsMin, e.RepsMax = 8, 6
				return e
			}()},
			wantErr: true,
		},
		{
			name: "missing reps_max",
			exercises: []Exercise{func() Exercise {
				e := validExercise()
				e.RepsMax = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "duplicate exercise names",
			exercises: []Exercise{
				validExercise(),
				func() Exercise {
					e := validExercise()
					e.Name = "FRONT SQUAT"
					return e
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validStrengthSession()
			s.Exercises = tt.exercises
			_, err := NewPlan("fuerza-boxeo", "Fuerza y boxeo", []SessionTemplate{s})
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanSessionNamesKeepsPlanOrder(t *testing.T) {
	p, err := NewPlan("fuerza-boxeo", "Fuerza y boxeo", []SessionTemplate{
		validStrengthSession(),
		validSportSession(),
		{Name: "Zone 2", Kind: KindCardio},
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
	p, err := NewPlan("fuerza-boxeo", "Fuerza y boxeo", []SessionTemplate{
		validStrengthSession(),
		validSportSession(),
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
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Session() error = %v, want ErrNotFound", err)
		}
	})
}

func TestPlanSessionsIsACopy(t *testing.T) {
	p, err := NewPlan("fuerza-boxeo", "Fuerza y boxeo", []SessionTemplate{validStrengthSession()})
	if err != nil {
		t.Fatalf("NewPlan() unexpected error: %v", err)
	}

	p.Sessions()[0].Name = "mutated"

	if got := p.SessionNames()[0]; got != "Gym A" {
		t.Errorf("plan was mutated through Sessions(): got %q", got)
	}
}
