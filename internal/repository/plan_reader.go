package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/frontmatter"
	"go.yaml.in/yaml/v3"
)

const planFileExt = ".md"

// PlanReader loads training plan templates from a directory of markdown files.
type PlanReader struct {
	dir string
}

func NewPlanReader(dir string) *PlanReader {
	return &PlanReader{dir: dir}
}

// planFile mirrors the YAML frontmatter of a plan file.
type planFile struct {
	Plan     string        `yaml:"plan"`
	Sessions []sessionFile `yaml:"sessions"`
}

type sessionFile struct {
	Name      string         `yaml:"name"`
	Kind      string         `yaml:"kind"`
	Modality  string         `yaml:"modality"`
	Exercises []exerciseFile `yaml:"exercises"`
}

type exerciseFile struct {
	Name       string `yaml:"name"`
	Block      string `yaml:"block"`
	Sets       int    `yaml:"sets"`
	RepsMin    int    `yaml:"reps_min"`
	RepsMax    int    `yaml:"reps_max"`
	BodyWeight bool   `yaml:"body_weight"`
}

func (r *PlanReader) LoadPlan(id string) (domain.Plan, error) {
	path := filepath.Join(r.dir, id+planFileExt)

	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.Plan{}, fmt.Errorf("plan %q: %w", id, domain.ErrNotFound)
		}
		return domain.Plan{}, fmt.Errorf("reading plan %q: %w", id, err)
	}

	fm, _, err := frontmatter.Split(raw)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("plan %q: %w", id, err)
	}

	var file planFile
	if err := yaml.Unmarshal(fm, &file); err != nil {
		return domain.Plan{}, fmt.Errorf("parsing plan %q: %w", id, err)
	}

	sessions, err := toSessions(file.Sessions)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("plan %q: %w", id, err)
	}

	plan, err := domain.NewPlan(id, file.Plan, sessions)
	if err != nil {
		return domain.Plan{}, err
	}

	return plan, nil
}

// ListPlans returns the ids of every plan file in the directory, sorted. The files are not opened or validated.
func (r *PlanReader) ListPlans() ([]string, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("plans directory %q: %w", r.dir, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("reading plans directory: %w", err)
	}

	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if !strings.EqualFold(ext, planFileExt) {
			continue
		}
		ids = append(ids, strings.TrimSuffix(e.Name(), ext))
	}
	sort.Strings(ids)

	return ids, nil
}

func toSessions(in []sessionFile) ([]domain.Session, error) {
	out := make([]domain.Session, 0, len(in))
	for _, s := range in {
		session, err := domain.NewSession(
			strings.TrimSpace(s.Name),
			domain.Kind(strings.ToLower(strings.TrimSpace(s.Kind))),
			strings.TrimSpace(s.Modality),
			toExercises(s.Exercises),
		)
		if err != nil {
			return nil, err
		}
		out = append(out, session)
	}
	return out, nil
}

func toExercises(in []exerciseFile) []domain.Exercise {
	out := make([]domain.Exercise, 0, len(in))
	for _, e := range in {
		repsMax := e.RepsMax
		if repsMax == 0 {
			// A single prescribed rep count may be written as reps_min only.
			repsMax = e.RepsMin
		}
		out = append(out, domain.Exercise{
			Name:       strings.TrimSpace(e.Name),
			Block:      strings.ToLower(strings.TrimSpace(e.Block)),
			Sets:       e.Sets,
			RepsMin:    e.RepsMin,
			RepsMax:    repsMax,
			BodyWeight: e.BodyWeight,
		})
	}
	return out
}
