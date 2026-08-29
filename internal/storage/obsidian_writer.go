package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/idoceb00/lorren/internal/domain"
)

type ObsidianWriter struct {
	dailyNotesDir string
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
	fmt.Fprintf(&b, "**Breakfast:** %s\n\n", log.Breakfast)
	fmt.Fprintf(&b, "**Lunch:** %s\n\n", log.Lunch)
	fmt.Fprintf(&b, "**Dinner:** %s\n\n", log.Dinner)
	fmt.Fprintf(&b, "**Snacks:** %s\n\n", log.Snacks)

	fmt.Fprintf(&b, "## 📆 Day\n\n")
	fmt.Fprintf(&b, "**What I did today:** %s\n\n", log.WhatIDidToday)
	fmt.Fprintf(&b, "**What went well:** %s\n\n", log.WhatWentWell)
	fmt.Fprintf(&b, "**What to improve:** %s\n\n", log.WhatToImprove)

	fmt.Fprintf(&b, "## 🎯 Quick notes\n\n%s\n", log.QuickNotes)

	return b.String()
}
