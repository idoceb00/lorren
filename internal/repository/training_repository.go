package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/idoceb00/lorren/internal/domain"
)

// TrainingRepository writes training session logs as markdown files with YAML frontmatter, one file per session.
type TrainingRepository struct {
	dir string
}

func NewTrainingRepository(dir string) *TrainingRepository {
	return &TrainingRepository{dir: dir}
}

func (r *TrainingRepository) SaveTrainingSession(s *domain.TrainingSession) (string, error) {
	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return "", fmt.Errorf("creating trainings dir: %w", err)
	}

	path := filepath.Join(r.dir, trainingFileName(s))

	if err := os.WriteFile(path, []byte(buildTrainingMarkdown(s)), 0o644); err != nil {
		return "", fmt.Errorf("writing training log file: %w", err)
	}

	return path, nil
}

func trainingFileName(s *domain.TrainingSession) string {
	return fmt.Sprintf("%s %s%s", s.Date.Format("2006-01-02"), sanitize(s.SessionName), planFileExt)
}

// sanitize removes characters that do not belong in a file name.
func sanitize(name string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-")
	return strings.TrimSpace(replacer.Replace(name))
}

func buildTrainingMarkdown(s *domain.TrainingSession) string {
	var b strings.Builder

	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "date: %s\n", s.Date.Format("2006-01-02"))
	fmt.Fprintf(&b, "session: %s\n", s.SessionName)
	fmt.Fprintf(&b, "kind: %s\n", s.Kind())
	if s.Modality != "" {
		fmt.Fprintf(&b, "modality: %s\n", s.Modality)
	}
	fmt.Fprintf(&b, "duration_min: %d\n", int(s.Duration.Minutes()))

	writeDetailFrontmatter(&b, s.Detail)

	fmt.Fprintf(&b, "---\n\n")

	fmt.Fprintf(&b, "# %s\n\n", s.SessionName)

	writeDetailBody(&b, s.Detail)

	fmt.Fprintf(&b, "## Notes\n\n")
	fmt.Fprintf(&b, "%s\n%s\n%s\n", markerStart("training_notes"), s.Notes, markerEnd("training_notes"))

	return b.String()
}

// writeDetailFrontmatter adds the queryable scalars each kind contributes.
func writeDetailFrontmatter(b *strings.Builder, detail domain.SessionDetail) {
	switch d := detail.(type) {
	case domain.StrengthDetail:
		fmt.Fprintf(b, "total_volume_kg: %.1f\n", totalVolume(d))
	case domain.CardioDetail:
		fmt.Fprintf(b, "activity: %s\n", d.Activity)
	}
}

// writeDetailBody writes the part that is read, not queried.
func writeDetailBody(b *strings.Builder, detail domain.SessionDetail) {
	d, ok := detail.(domain.StrengthDetail)
	if !ok {
		return
	}

	block := ""
	for _, e := range d.Exercises {
		if e.Block != block {
			if block != "" {
				fmt.Fprintf(b, "\n")
			}
			block = e.Block
			fmt.Fprintf(b, "## %s\n\n", blockHeading(block))
			fmt.Fprintf(b, "| Exercise | Prescribed | Weight | Reps |\n")
			fmt.Fprintf(b, "|---|---|---|---|\n")
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s |\n", e.Name, prescribed(e), weightCell(e), repsCell(e))
	}
	fmt.Fprintf(b, "\n")
}

func blockHeading(block string) string {
	if block == "" {
		return "Exercises"
	}
	return strings.ToUpper(block[:1]) + block[1:]
}

func prescribed(e domain.PerformedExercise) string {
	if e.RepsMin == e.RepsMax {
		return fmt.Sprintf("%d x %d", e.Sets, e.RepsMin)
	}
	return fmt.Sprintf("%d x %d-%d", e.Sets, e.RepsMin, e.RepsMax)
}

func weightCell(e domain.PerformedExercise) string {
	switch {
	case !e.Done:
		return "—"
	case e.Weight == nil:
		return "bw"
	default:
		return fmt.Sprintf("%g kg", *e.Weight)
	}
}

func repsCell(e domain.PerformedExercise) string {
	if !e.Done {
		return "—"
	}
	return fmt.Sprintf("%d", e.Reps)
}

// totalVolume is weight times reps summed over the loaded exercises actually performed.
func totalVolume(d domain.StrengthDetail) float64 {
	var total float64
	for _, e := range d.Exercises {
		if !e.Done || e.Weight == nil {
			continue
		}
		total += *e.Weight * float64(e.Reps) * float64(e.Sets)
	}
	return total
}
