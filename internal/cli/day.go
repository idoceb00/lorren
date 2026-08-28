package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Log today's habits",
	Long:  `Day starts an interactive wizard tha asks about your daily habits and writes the result as a markdown file into your Obsidian vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("not implemented yet")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dayCmd)
}
