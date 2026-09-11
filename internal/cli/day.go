package cli

import (
	"errors"
	"fmt"
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
		var repositoryPort domain.DailyRepository = repository.NewDailyRepository(appConfig.DailyNotesPath())

		show, err := cmd.Flags().GetBool("show")
		if err != nil {
			return err
		}

		if show {
			log, err := service.ShowDay(repositoryPort, time.Now())
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					fmt.Fprintln(cmd.OutOrStdout(), "No daily log for today yet.")
					return nil
				}

				return err
			}

			renderDailyLog(cmd.OutOrStdout(), log)

			return nil
		}

		var interviewerPort domain.DailyInterviewer = interviewer.NewHuhInterviewer()

		path, err := service.RecordDay(interviewerPort, repositoryPort, time.Now())
		if err != nil {
			return err
		}

		fmt.Printf("Saved to %s\n", path)

		return nil
	},
}

// Special function executed automatically when the package is loaded
func init() {
	dayCmd.Flags().Bool("show", false, "print today's log instead of starting the wizard")
	rootCmd.AddCommand(dayCmd)
}
