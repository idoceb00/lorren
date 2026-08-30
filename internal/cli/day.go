package cli

import (
	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/interviewer"
	"github.com/idoceb00/lorren/internal/storage"
	"github.com/spf13/cobra"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Log today's habits",
	Long:  `Day starts an interactive wizard tha asks about your daily habits and writes the result as a markdown file into your Obsidian vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var interviewerPort domain.Interviewer = interviewer.NewHuhInterviewer()
		var repositoryPort domain.Repository = storage.NewObsidianWriter(appConfig.VaultPath)

		log, err := interviewerPort.AskDailyLog()
		if err != nil {
			return err
		}

		return repositoryPort.SaveDailyLog(log)
	},
}

// Special function executed automatically when the package is loaded
func init() {
	rootCmd.AddCommand(dayCmd)
}
