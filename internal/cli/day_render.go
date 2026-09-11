package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/idoceb00/lorren/internal/domain"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Width(12)

	bodyStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Width(72)

	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	missingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	emptyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

const emptyValue = "—"

func renderDailyLog(w io.Writer, log *domain.DailyLog) {
	fmt.Fprintln(w, titleStyle.Render(log.Date.Format("Monday, 2 January 2006")))

	fmt.Fprintln(w, sectionStyle.Render("Habits"))
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

	fmt.Fprintln(w, sectionStyle.Render("Meals"))
	for _, meal := range []struct{ label, text string }{
		{"Breakfast", log.Breakfast},
		{"Lunch", log.Lunch},
		{"Dinner", log.Dinner},
		{"Snacks", log.Snacks},
	} {
		fmt.Fprintf(w, "  %s%s\n", labelStyle.Render(meal.label+":"), orEmpty(meal.text))
	}

	fmt.Fprintf(w, "\n%s\n  %s Well spent\n", sectionStyle.Render("The day"), mark(log.DayWellSpent))

	for _, section := range []struct{ label, text string }{
		{"What I did", log.WhatIDidToday},
		{"What went well", log.WhatWentWell},
		{"What to improve", log.WhatToImprove},
		{"Quick notes", log.QuickNotes},
	} {
		fmt.Fprintf(w, "\n  %s\n%s\n", labelStyle.Render(section.label), bodyStyle.Render(orEmpty(section.text)))
	}
}

func mark(done bool) string {
	if done {
		return doneStyle.Render("✓")
	}

	return missingStyle.Render("✗")
}

func orEmpty(text string) string {
	if strings.TrimSpace(text) == "" {
		return emptyStyle.Render(emptyValue)
	}

	return text
}
