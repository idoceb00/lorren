package cli

import (
	"time"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/interviewer"
	"github.com/idoceb00/lorren/internal/repository"
	"github.com/idoceb00/lorren/internal/service"
	"github.com/spf13/cobra"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Log today's habits",
	Long:  `Day starts an interactive wizard that asks about your daily habits and writes the result as a markdown file with YAML frontmatter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var interviewerPort domain.Interviewer = interviewer.NewHuhInterviewer()
		var repositoryPort domain.Repository = repository.NewMarkdownRepository(appConfig.DailyNotesPath())

		return service.RecordDay(interviewerPort, repositoryPort, time.Now())
	},
}

// Special function executed automatically when the package is loaded
func init() {
	rootCmd.AddCommand(dayCmd)
}
