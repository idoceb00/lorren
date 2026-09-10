package repository_test

import (
	"errors"
	"testing"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/repository"
)

func TestPlanReaderLoadPlan(t *testing.T) {
	r := repository.NewPlanReader("testdata")

	plan, err := r.LoadPlan("valid")
	if err != nil {
		t.Fatalf("LoadPlan() unexpected error: %v", err)
	}

	if plan.ID() != "valid" {
		t.Errorf("ID() = %q, want %q", plan.ID(), "valid")
	}
	if plan.Name() != "Strength and boxing" {
		t.Errorf("Name() = %q, want %q", plan.Name(), "Strength and boxing")
	}

	want := []string{"Gym A", "Boxing", "Zone 2"}
	got := plan.SessionNames()
	if len(got) != len(want) {
		t.Fatalf("SessionNames() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SessionNames()[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	gym, err := plan.Session("Gym A")
	if err != nil {
		t.Fatalf("Session() unexpected error: %v", err)
	}
	if gym.Kind != domain.KindStrength {
		t.Errorf("Kind = %q, want %q", gym.Kind, domain.KindStrength)
	}
	if len(gym.Exercises) != 2 {
		t.Fatalf("got %d exercises, want 2", len(gym.Exercises))
	}

	squat := gym.Exercises[0]
	if squat.Name != "Front squat" || squat.Sets != 3 || squat.RepsMin != 4 || squat.RepsMax != 6 {
		t.Errorf("front squat mapped as %+v", squat)
	}
	if squat.BodyWeight {
		t.Error("front squat should not be bodyweight")
	}

	// reps_max is omitted in the file and defaults to reps_min.
	jump := gym.Exercises[1]
	if jump.RepsMin != 3 || jump.RepsMax != 3 {
		t.Errorf("vertical jump reps = %d-%d, want 3-3", jump.RepsMin, jump.RepsMax)
	}
	if !jump.BodyWeight {
		t.Error("vertical jump should be bodyweight")
	}
}

func TestPlanReaderLoadPlanErrors(t *testing.T) {
	r := repository.NewPlanReader("testdata")

	tests := []struct {
		name string
		id   string
	}{
		{name: "unknown kind", id: "unknown-kind"},
		{name: "strength session without exercises", id: "strength-without-exercises"},
		{name: "no frontmatter", id: "no-frontmatter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := r.LoadPlan(tt.id); err == nil {
				t.Fatal("LoadPlan() expected an error, got nil")
			}
		})
	}
}

func TestPlanReaderLoadPlanNotFound(t *testing.T) {
	r := repository.NewPlanReader("testdata")

	_, err := r.LoadPlan("does-not-exist")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("LoadPlan() error = %v, want ErrNotFound", err)
	}
}

func TestPlanReaderListPlans(t *testing.T) {
	r := repository.NewPlanReader("testdata")

	ids, err := r.ListPlans()
	if err != nil {
		t.Fatalf("ListPlans() unexpected error: %v", err)
	}

	// Broken plans are listed too: ListPlans does not open the files.
	if len(ids) < 4 {
		t.Errorf("ListPlans() = %v, want every markdown file in testdata", ids)
	}

	var found bool
	for _, id := range ids {
		if id == "valid" {
			found = true
		}
	}
	if !found {
		t.Errorf("ListPlans() = %v, missing %q", ids, "valid")
	}
}
