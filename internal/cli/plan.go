package cli

import (
	"fmt"
	"strings"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/repository"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show the active training plan",
	Long:  "Plan loads the active training plan template and prints it as lorren understands it. Use it to check a plan after editing ti by hand.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var planReader domain.PlanReader = repository.NewPlanReader(appConfig.PlansPath())

		plan, err := appConfig.ResolveActivePlan(planReader)
		if err != nil {
			return err
		}

		printPlan(plan)

		return nil
	},
}

func printPlan(plan domain.Plan) {
	fmt.Printf("%s (%s)\n", plan.Name(), plan.ID())

	for _, s := range plan.Sessions() {
		fmt.Printf("\n%s [%s", s.Name, s.Kind)
		if s.Modality != "" {
			fmt.Printf(", %s", s.Modality)
		}
		fmt.Printf("]\n")

		if len(s.Exercises) == 0 {
			continue
		}

		block := ""
		for _, e := range s.Exercises {
			if e.Block != block {
				block = e.Block
				fmt.Printf(" %s\n", strings.ToUpper(block))
			}
			fmt.Printf("    %-45s %d x %s%s\n", e.Name, e.Sets, reps(e), weightNote(e))
		}
	}
}

func reps(e domain.Exercise) string {
	if e.RepsMin == e.RepsMax {
		return fmt.Sprintf("%d", e.RepsMin)
	}
	return fmt.Sprintf("%d-%d", e.RepsMin, e.RepsMax)
}

func weightNote(e domain.Exercise) string {
	if e.BodyWeight {
		return "  (bodyweight)"
	}
	return ""
}

func init() {
	rootCmd.AddCommand(planCmd)
}
