package domain_test

import (
	"testing"

	"github.com/idoceb00/lorren/internal/domain"
)

// This is the validation the whole template design rests on: only strength
// sessions prescribe exercises.
func TestNewSessionTemplateValidatesByKind(t *testing.T) {
	tests := []struct {
		name      string
		session   string
		kind      domain.Kind
		exercises []domain.Exercise
		wantErr   bool
	}{
		{
			name:      "strength with exercises",
			session:   "Gym A",
			kind:      domain.KindStrength,
			exercises: []domain.Exercise{validExercise()},
		},
		{
			name:    "strength without exercises",
			session: "Gym A",
			kind:    domain.KindStrength,
			wantErr: true,
		},
		{
			name:    "sport without exercises",
			session: "Boxeo",
			kind:    domain.KindSport,
		},
		{
			name:      "sport with exercises",
			session:   "Boxeo",
			kind:      domain.KindSport,
			exercises: []domain.Exercise{validExercise()},
			wantErr:   true,
		},
		{
			name:    "cardio without exercises",
			session: "Zone 2",
			kind:    domain.KindCardio,
		},
		{
			name:      "cardio with exercises",
			session:   "Zone 2",
			kind:      domain.KindCardio,
			exercises: []domain.Exercise{validExercise()},
			wantErr:   true,
		},
		{
			name:    "unknown kind",
			session: "Gym A",
			kind:    domain.Kind("gimnasio"),
			wantErr: true,
		},
		{
			name:    "empty kind",
			session: "Gym A",
			wantErr: true,
		},
		{
			name:      "missing session name",
			session:   "  ",
			kind:      domain.KindStrength,
			exercises: []domain.Exercise{validExercise()},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewSession(tt.session, tt.kind, "fuerza", tt.exercises)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewSessionTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSessionTemplateValidatesExercises(t *testing.T) {
	tests := []struct {
		name      string
		exercises []domain.Exercise
		wantErr   bool
	}{
		{
			name:      "valid exercise",
			exercises: []domain.Exercise{validExercise()},
		},
		{
			name: "fixed reps",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.Name = "Salto vertical"
				e.RepsMin, e.RepsMax = 3, 3
				e.BodyWeight = true
				return e
			}()},
		},
		{
			name: "missing name",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.Name = " "
				return e
			}()},
			wantErr: true,
		},
		{
			name: "zero sets",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.Sets = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "zero reps_min",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.RepsMin = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "reps_max below reps_min",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.RepsMin, e.RepsMax = 8, 6
				return e
			}()},
			wantErr: true,
		},
		{
			name: "missing reps_max",
			exercises: []domain.Exercise{func() domain.Exercise {
				e := validExercise()
				e.RepsMax = 0
				return e
			}()},
			wantErr: true,
		},
		{
			name: "duplicate exercise names",
			exercises: []domain.Exercise{
				validExercise(),
				func() domain.Exercise {
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
			_, err := domain.NewSession("Gym A", domain.KindStrength, "fuerza", tt.exercises)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewSessionTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
