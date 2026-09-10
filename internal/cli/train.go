package cli

import (
	"fmt"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/interviewer"
	"github.com/idoceb00/lorren/internal/repository"
	"github.com/idoceb00/lorren/internal/service"
	"github.com/spf13/cobra"
)

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Log a training session",
	Long:  `Train starts an interactive wizard over your active plan and writes the session as a markdown file with YAML frontmatter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var planReader domain.PlanReader = repository.NewPlanReader(appConfig.PlansPath())

		plan, err := appConfig.ResolveActivePlan(planReader)
		if err != nil {
			return err
		}

		interviewerPort := interviewer.NewHuhInterviewer()
		repo := repository.NewTrainingRepository(appConfig.TrainingsPath())

		path, err := service.RecordTraining(interviewerPort, repo, plan)
		if err != nil {
			return err
		}

		fmt.Printf("Saved to %s\n", path)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(trainCmd)
}
