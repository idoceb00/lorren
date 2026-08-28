package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "lorren",
	Short: "Log your training sessions and daily habits to your Obsidian vault",
	Long:  `Lorren is a CLI wizard that interviews you about your training sessions and daily habits, then writes the results as structured markdown files into your Obsidian vault.`,
}

// Execute runs the root command. It's the single entry point called from main.go - main stays a thin wrapper around this.
func Execute() error {
	return rootCmd.Execute()
}
