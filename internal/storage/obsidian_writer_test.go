package storage_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/storage"
)

// frontmatter mirrors the scalar fields ObsidianWriter writes to YAML,
// used to parse and assert on the generated frontmatter block.
type frontmatter struct {
	Date         string  `yaml:"date"`
	Training     bool    `yaml:"training"`
	Reading      bool    `yaml:"reading"`
	Coding       bool    `yaml:"coding"`
	Meditation   bool    `yaml:"meditation"`
	NoSmoking    bool    `yaml:"no_smoking"`
	Stretching   bool    `yaml:"stretching"`
	SleepHours   float64 `yaml:"sleep_hours"`
	DayWellSpent bool    `yaml:"day_well_spent"`
}

func markedField(key, value string) string {
	return fmt.Sprintf("<!-- lorren:%s:start -->\n%s\n<!-- lorren:%s:end -->", key, value, key)
}

func labeledField(label, key, value string) string {
	return fmt.Sprintf("**%s:**\n%s", label, markedField(key, value))
}

func TestObsidianWriter_SaveDailyLog(t *testing.T) {
	fixedDate := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)

	fullLog := &domain.DailyLog{
		Date:          fixedDate,
		Training:      true,
		Reading:       false,
		Coding:        true,
		Meditation:    false,
		NoSmoking:     true,
		Stretching:    false,
		SleepHours:    7.5,
		Breakfast:     "oats and coffee",
		Lunch:         "chicken and rice",
		Dinner:        "salad",
		Snacks:        "almonds",
		DayWellSpent:  true,
		WhatIDidToday: "gym and work",
		WhatWentWell:  "good focus",
		WhatToImprove: "sleep earlier",
		QuickNotes:    "felt strong today",
	}

	emptyLog := &domain.DailyLog{
		Date: fixedDate,
		// all bools false, all strings "", SleepHours 0 — zero value struct
	}

	tests := []struct {
		name string
		// subdir is joined onto t.TempDir() as the writer's target directory.
		// Leaving it empty uses the temp dir itself (already exists).
		// A non-empty value that doesn't exist yet exercises os.MkdirAll.
		subdir            string
		log               *domain.DailyLog
		wantFrontmatter   frontmatter
		wantBodyFragments []string
		wantErr           bool
	}{
		{
			name:   "writes full log to existing directory",
			subdir: "",
			log:    fullLog,
			wantFrontmatter: frontmatter{
				Date:         "2026-09-05",
				Training:     true,
				Reading:      false,
				Coding:       true,
				Meditation:   false,
				NoSmoking:    true,
				Stretching:   false,
				SleepHours:   7.5,
				DayWellSpent: true,
			},
			wantBodyFragments: []string{
				labeledField("Breakfast", "breakfast", "oats and coffee"),
				labeledField("Lunch", "lunch", "chicken and rice"),
				labeledField("Dinner", "dinner", "salad"),
				labeledField("Snacks", "snacks", "almonds"),
				labeledField("What I did today", "what_i_did_today", "gym and work"),
				labeledField("What went well", "what_went_well", "good focus"),
				labeledField("What to improve", "what_to_improve", "sleep earlier"),
				markedField("quick_notes", "felt strong today"),
			},
			wantErr: false,
		},
		{
			name:   "creates nested directory when it does not exist yet",
			subdir: filepath.Join("vault", "daily"),
			log:    fullLog,
			wantFrontmatter: frontmatter{
				Date:         "2026-09-05",
				Training:     true,
				Reading:      false,
				Coding:       true,
				Meditation:   false,
				NoSmoking:    true,
				Stretching:   false,
				SleepHours:   7.5,
				DayWellSpent: true,
			},
			wantBodyFragments: []string{
				labeledField("Breakfast", "breakfast", "oats and coffee"),
			},
			wantErr: false,
		},
		{
			name:   "writes log with all zero-value fields",
			subdir: "",
			log:    emptyLog,
			wantFrontmatter: frontmatter{
				Date:         "2026-09-05",
				Training:     false,
				Reading:      false,
				Coding:       false,
				Meditation:   false,
				NoSmoking:    false,
				Stretching:   false,
				SleepHours:   0,
				DayWellSpent: false,
			},
			wantBodyFragments: []string{
				markedField("breakfast", ""),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.subdir != "" {
				dir = filepath.Join(dir, tt.subdir)
			}
			writer := storage.NewObsidianWriter(dir)

			gotErr := writer.SaveDailyLog(tt.log)
			if gotErr != nil {
				if !tt.wantErr {
					t.Fatalf("SaveDailyLog() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SaveDailyLog() succeeded unexpectedly")
			}

			wantPath := filepath.Join(dir, tt.log.Date.Format("2006-01-02")+".md")
			raw, err := os.ReadFile(wantPath)
			if err != nil {
				t.Fatalf("expected file at %s, read failed: %v", wantPath, err)
			}
			content := string(raw)

			// --- frontmatter: parse and compare as structured YAML ---
			parts := strings.SplitN(content, "---", 3)
			if len(parts) < 3 {
				t.Fatalf("content does not have a valid frontmatter block:\n%s", content)
			}

			var got frontmatter
			if err := yaml.Unmarshal([]byte(parts[1]), &got); err != nil {
				t.Fatalf("frontmatter is not valid YAML: %v", err)
			}

			if diff := cmp.Diff(tt.wantFrontmatter, got); diff != "" {
				t.Errorf("frontmatter mismatch (-want +got):\n%s", diff)
			}

			// --- body: check key fragments are present ---
			body := parts[2]
			for _, fragment := range tt.wantBodyFragments {
				if !strings.Contains(body, fragment) {
					t.Errorf("body missing expected fragment: %q\nfull body:\n%s", fragment, body)
				}
			}
		})
	}
}
