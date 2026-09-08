package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/idoceb00/lorren/internal/domain"
	"go.yaml.in/yaml/v3"
)

type ObsidianWriter struct {
	dailyNotesDir string
}

// frontmatter mirrors the scalar fields written to the YAML block.
// Field order doesn't matter for yaml.Unmarshal, only the tags do.
type frontmatter struct {
	Training     bool    `yaml:"training"`
	Reading      bool    `yaml:"reading"`
	Coding       bool    `yaml:"coding"`
	Meditation   bool    `yaml:"meditation"`
	NoSmoking    bool    `yaml:"no_smoking"`
	Stretching   bool    `yaml:"stretching"`
	SleepHours   float64 `yaml:"sleep_hours"`
	DayWellSpent bool    `yaml:"day_well_spent"`
}

var fieldPatterns = map[string]*regexp.Regexp{
	"breakfast":        compileFieldPattern("breakfast"),
	"lunch":            compileFieldPattern("lunch"),
	"dinner":           compileFieldPattern("dinner"),
	"snacks":           compileFieldPattern("snacks"),
	"what_i_did_today": compileFieldPattern("what_i_did_today"),
	"what_went_well":   compileFieldPattern("what_went_well"),
	"what_to_improve":  compileFieldPattern("what_to_improve"),
	"quick_notes":      compileFieldPattern("quick_notes"),
}

func NewObsidianWriter(dir string) *ObsidianWriter {
	return &ObsidianWriter{dailyNotesDir: dir}
}

func (w *ObsidianWriter) SaveDailyLog(log *domain.DailyLog) error {
	if err := os.MkdirAll(w.dailyNotesDir, 0o755); err != nil {
		return fmt.Errorf("creating daily notes dir: %w", err)
	}

	filename := log.Date.Format("2006-01-02") + ".md"
	path := filepath.Join(w.dailyNotesDir, filename)

	content := buildMarkdown(log)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing daily log file: %w", err)
	}

	return nil
}

func (w *ObsidianWriter) FindByDate(date time.Time) (*domain.DailyLog, error) {
	filename := date.Format("2006-01-02") + ".md"
	path := filepath.Join(w.dailyNotesDir, filename)

	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("reading daily log file: %w", err)
	}

	fm, body, err := splitFrontmatter(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parsing daily log fil: %w", err)
	}

	var meta frontmatter
	if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	return domain.NewDailyLog(domain.NewDailyLogInput{
		Date:          date,
		Training:      meta.Training,
		Reading:       meta.Reading,
		Coding:        meta.Coding,
		Meditation:    meta.Meditation,
		NoSmoking:     meta.NoSmoking,
		Stretching:    meta.Stretching,
		SleepHours:    meta.SleepHours,
		DayWellSpent:  meta.DayWellSpent,
		Breakfast:     extractField(body, "breakfast"),
		Lunch:         extractField(body, "lunch"),
		Dinner:        extractField(body, "dinner"),
		Snacks:        extractField(body, "snacks"),
		WhatIDidToday: extractField(body, "what_i_did_today"),
		WhatWentWell:  extractField(body, "what_went_well"),
		WhatToImprove: extractField(body, "what_to_improve"),
		QuickNotes:    extractField(body, "quick_notes"),
	})
}

// splitFrontmatter separates the leading YAML block (between --- markers)
// from the rest of the markdown body.
func splitFrontmatter(content string) (fm, body string, err error) {
	const delim = "---\n"
	if !strings.HasPrefix(content, delim) {
		return "", "", fmt.Errorf("missing frontmatter opening delimiter")
	}
	rest := content[len(delim):]
	idx := strings.Index(rest, delim)
	if idx == -1 {
		return "", "", fmt.Errorf("missing frontmatter closing delimiter")
	}
	return rest[:idx], rest[idx+len(delim):], nil
}

func extractField(body, key string) string {
	re, ok := fieldPatterns[key]
	if !ok {
		return ""
	}
	match := re.FindStringSubmatch(body)
	if match == nil {
		return ""
	}
	return match[1]
}

func compileFieldPattern(key string) *regexp.Regexp {
	pattern := fmt.Sprintf(`(?s)%s\n(.*?)\n%s`, regexp.QuoteMeta(markerStart(key)), regexp.QuoteMeta(markerEnd(key)))
	return regexp.MustCompile(pattern)
}

func buildMarkdown(log *domain.DailyLog) string {
	var b strings.Builder

	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "date: %s\n", log.Date.Format("2006-01-02"))
	fmt.Fprintf(&b, "training: %t\n", log.Training)
	fmt.Fprintf(&b, "reading: %t\n", log.Reading)
	fmt.Fprintf(&b, "coding: %t\n", log.Coding)
	fmt.Fprintf(&b, "meditation: %t\n", log.Meditation)
	fmt.Fprintf(&b, "no_smoking: %t\n", log.NoSmoking)
	fmt.Fprintf(&b, "stretching: %t\n", log.Stretching)
	fmt.Fprintf(&b, "sleep_hours: %.2f\n", log.SleepHours)
	fmt.Fprintf(&b, "day_well_spent: %t\n", log.DayWellSpent)
	fmt.Fprintf(&b, "---\n\n")

	fmt.Fprintf(&b, "## 🍽️ Meals\n\n")
	writeField(&b, "Breakfast", "breakfast", log.Breakfast)
	writeField(&b, "Lunch", "lunch", log.Lunch)
	writeField(&b, "Dinner", "dinner", log.Dinner)
	writeField(&b, "Snacks", "snacks", log.Snacks)

	fmt.Fprintf(&b, "## 📆 Day\n\n")
	writeField(&b, "What I did today", "what_i_did_today", log.WhatIDidToday)
	writeField(&b, "What went well", "what_went_well", log.WhatWentWell)
	writeField(&b, "What to improve", "what_to_improve", log.WhatToImprove)

	fmt.Fprintf(&b, "## 🎯 Quick notes\n\n")
	writeField(&b, "", "quick_notes", log.QuickNotes)

	return b.String()
}

func writeField(b *strings.Builder, label, key, value string) {
	if label != "" {
		fmt.Fprintf(b, "**%s:**\n", label)
	}
	fmt.Fprintf(b, "%s\n%s\n%s\n\n", markerStart(key), value, markerEnd(key))
}

func markerStart(key string) string {
	return fmt.Sprintf("<!-- lorren:%s:start -->", key)
}

func markerEnd(key string) string {
	return fmt.Sprintf("<!-- lorren:%s:end -->", key)
}
