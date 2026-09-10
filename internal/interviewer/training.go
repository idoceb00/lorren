package interviewer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/idoceb00/lorren/internal/domain"
)

// exerciseInput holds the raw strings huh collects for one exercise. Weight is left empty when the exercise was not performed.
type exerciseInput struct {
	template domain.Exercise
	weight   string
	reps     string
	done     bool
}

func (h *HuhInterviewer) AskTrainingLog(plan domain.Plan) (*domain.TrainingLog, error) {
	name, err := askSessionName(plan)
	if err != nil {
		return nil, err
	}

	template, err := plan.Session(name)
	if err != nil {
		return nil, err
	}

	detail, err := askDetail(template)
	if err != nil {
		return nil, err
	}

	duration, notes, err := askCommon(template)
	if err != nil {
		return nil, err
	}

	return domain.NewTrainingSession(domain.NewTrainingLogInput{
		Date:        time.Now(),
		SessionName: template.Name,
		Modality:    template.Modality,
		Duration:    duration,
		Notes:       notes,
		Detail:      detail,
	})
}

func askSessionName(plan domain.Plan) (string, error) {
	var name string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What did you train today?").
				Options(huh.NewOptions(plan.SessionNames()...)...).
				Value(&name),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}

	return name, nil
}

// askDetail dispatches on kind: each kind collects its own detail.
func askDetail(t domain.Session) (domain.TrainingDetail, error) {
	switch t.Kind {
	case domain.KindStrength:
		return askStrengthDetail(t)
	case domain.KindCardio:
		return askCardioDetail()
	case domain.KindSport:
		return domain.SportDetail{}, nil
	default:
		return nil, fmt.Errorf("session %q: unknown kind %q", t.Name, t.Kind)
	}
}

func askStrengthDetail(t domain.Session) (domain.TrainingDetail, error) {
	inputs := make([]*exerciseInput, 0, len(t.Exercises))
	for _, e := range t.Exercises {
		inputs = append(inputs, &exerciseInput{template: e, done: true})
	}

	groups := buildExerciseGroups(inputs)

	if err := huh.NewForm(groups...).Run(); err != nil {
		return nil, err
	}

	exercises := make([]domain.PerformedExercise, 0, len(inputs))
	for _, in := range inputs {
		e, err := in.toPerformed()
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}

	return domain.StrengthDetail{Exercises: exercises}, nil
}

// buildExerciseGroups puts one huh group per block, keeping template order.
func buildExerciseGroups(inputs []*exerciseInput) []*huh.Group {
	var (
		groups []*huh.Group
		fields []huh.Field
		block  string
	)

	flush := func() {
		if len(fields) > 0 {
			groups = append(groups, huh.NewGroup(fields...).Title(blockTitle(block)))
			fields = nil
		}
	}

	for _, in := range inputs {
		if in.template.Block != block {
			flush()
			block = in.template.Block
		}
		fields = append(fields, in.fields()...)
	}
	flush()

	return groups
}

func blockTitle(block string) string {
	if block == "" {
		return "Exercises"
	}
	return strings.ToUpper(block[:1]) + block[1:]
}

// fields builds the huh widgets for one exercise: a confirm for bodyweight work, a weight input otherwise.
func (in *exerciseInput) fields() []huh.Field {
	label := fmt.Sprintf("%s  %d x %s", in.template.Name, in.template.Sets, prescribedReps(in.template))

	if in.template.BodyWeight {
		return []huh.Field{
			huh.NewConfirm().
				Title(label).
				Value(&in.done),
			huh.NewInput().
				Title("  reps").
				Value(&in.reps).
				Validate(func(s string) error {
					return validateReps(s, in.done)
				}),
		}
	}

	return []huh.Field{
		huh.NewInput().
			Title(label).
			Description("kg, leave empty if you skipped it").
			Value(&in.weight).
			Validate(validateOptionalFloat),
		huh.NewInput().
			Title("  reps").
			Value(&in.reps).
			Validate(func(s string) error {
				return validateReps(s, strings.TrimSpace(in.weight) != "")
			}),
	}
}

func validateReps(s string, performed bool) error {
	s = strings.TrimSpace(s)

	if !performed {
		if s != "" {
			return fmt.Errorf("you skipped this one, leave reps empty")
		}
		return nil
	}

	if s == "" {
		return fmt.Errorf("how many reps?")
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a whole number")
	}
	if n < 1 {
		return fmt.Errorf("must be at least 1")
	}

	return nil
}

func (in *exerciseInput) toPerformed() (domain.PerformedExercise, error) {
	e := domain.PerformedExercise{
		Name:       in.template.Name,
		Block:      in.template.Block,
		Sets:       in.template.Sets,
		RepsMin:    in.template.RepsMin,
		RepsMax:    in.template.RepsMax,
		BodyWeight: in.template.BodyWeight,
	}

	if in.template.BodyWeight {
		e.Done = in.done
	} else {
		e.Done = strings.TrimSpace(in.weight) != ""
	}

	if !e.Done {
		return e, nil
	}

	reps, err := strconv.Atoi(strings.TrimSpace(in.reps))
	if err != nil {
		return domain.PerformedExercise{}, fmt.Errorf("exercise %q: reading reps: %w", e.Name, err)
	}
	e.Reps = reps

	if !e.BodyWeight {
		weight, err := strconv.ParseFloat(strings.TrimSpace(in.weight), 64)
		if err != nil {
			return domain.PerformedExercise{}, fmt.Errorf("exercise %q: reading weight: %w", e.Name, err)
		}
		e.Weight = &weight
	}

	return e, nil
}

func askCardioDetail() (domain.TrainingDetail, error) {
	var activity string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What did you do?").
				Description("Running, bike, rower...").
				Value(&activity),
		),
	)
	if err := form.Run(); err != nil {
		return nil, err
	}

	return domain.CardioDetail{Activity: strings.TrimSpace(activity)}, nil
}

func askCommon(t domain.Session) (time.Duration, string, error) {
	var (
		minutes string
		notes   string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("How long did it take?").
				Description("Minutes.").
				Value(&minutes).
				Validate(validateOptionalInt),
			huh.NewText().
				Title("Notes").
				Description("How it felt, anything worth remembering.").
				Value(&notes),
		).Title(t.Name),
	)
	if err := form.Run(); err != nil {
		return 0, "", err
	}

	trimmed := strings.TrimSpace(minutes)
	if trimmed == "" {
		return 0, notes, nil
	}

	m, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, "", fmt.Errorf("reading duration: %w", err)
	}

	return time.Duration(m) * time.Minute, notes, nil
}

func prescribedReps(e domain.Exercise) string {
	if e.RepsMin == e.RepsMax {
		return strconv.Itoa(e.RepsMin)
	}
	return fmt.Sprintf("%d-%d", e.RepsMin, e.RepsMax)
}

func validateOptionalInt(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("must be a whole number")
	}
	return nil
}

func validateOptionalFloat(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return fmt.Errorf("must be a number")
	}
	return nil
}
