package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/idoceb00/lorren/internal/domain"
)

const emptyValue = "—"

func renderDailyLog(w io.Writer, log *domain.DailyLog) {
	fmt.Fprintf(w, "%s\n\n", log.Date.Format("Monday, 2 January 2006"))

	fmt.Fprintln(w, "Habits")
	for _, habit := range []struct {
		label string
		done  bool
	}{
		{"Training", log.Training},
		{"Reading", log.Reading},
		{"Coding", log.Coding},
		{"Meditation", log.Meditation},
		{"No smoking", log.NoSmoking},
		{"Stretching", log.Stretching},
	} {
		fmt.Fprintf(w, "  %s %s\n", mark(habit.done), habit.label)
	}

	fmt.Fprintf(w, "  %.1f hours of sleep\n\n", log.SleepHours)

	fmt.Fprintln(w, "Meals")
	for _, meal := range []struct{ label, text string }{
		{"Breakfast", log.Breakfast},
		{"Lunch", log.Lunch},
		{"Dinner", log.Dinner},
		{"Snacks", log.Snacks},
	} {
		fmt.Fprintf(w, "  %-11s %s\n", meal.label+":", orEmpty(meal.text))
	}

	fmt.Fprintf(w, "\nThe day\n  %s Well spent\n", mark(log.DayWellSpent))

	for _, section := range []struct{ label, text string }{
		{"What I did", log.WhatIDidToday},
		{"What went well", log.WhatWentWell},
		{"What to improve", log.WhatToImprove},
		{"Quick notes", log.QuickNotes},
	} {
		fmt.Fprintf(w, "\n  %s\n  %s\n", section.label, orEmpty(section.text))
	}
}

func mark(done bool) string {
	if done {
		return "✓"
	}

	return "✗"
}

func orEmpty(text string) string {
	if strings.TrimSpace(text) == "" {
		return emptyValue
	}

	return text
}
