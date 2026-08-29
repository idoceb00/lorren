package storage

import (
	"fmt"
	"os"
	"path/filepath"

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
	return fmt.Sprintf(`---
date: %s
reading: %t
coding: %t
meditation: %t
no_smoking: %t
stretching: %t
sleep_hours: %.2f
---
`,
		log.Date.Format("2006-01-02"),
		log.Reading,
		log.Coding,
		log.Meditation,
		log.NoSmoking,
		log.Stretching,
		log.SleepHours,
	)
}
