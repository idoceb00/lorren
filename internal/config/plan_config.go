package config

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/idoceb00/lorren/internal/domain"
)

// ResolveActivePlan loads the active plan, asking the user to choose one the
// first time and persisting that choice. A configured plan that no longer
// exists is treated as no choice at all, so renaming a plan file does not
// leave the tool stuck.
func (c *Config) ResolveActivePlan(reader domain.PlanReader) (domain.Plan, error) {
	if c.ActivePlan != "" {
		plan, err := reader.LoadPlan(c.ActivePlan)
		if err == nil {
			return plan, nil
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return domain.Plan{}, err
		}
		fmt.Printf("Plan %q is no longer in %s, pick another one.\n\n", c.ActivePlan, c.PlansPath())
	}

	id, err := c.choosePlan(reader)
	if err != nil {
		return domain.Plan{}, err
	}

	plan, err := reader.LoadPlan(id)
	if err != nil {
		return domain.Plan{}, err
	}

	if err := SetActivePlan(id); err != nil {
		return domain.Plan{}, err
	}
	c.ActivePlan = id

	return plan, nil
}

func (c *Config) choosePlan(reader domain.PlanReader) (string, error) {
	ids, err := reader.ListPlans()
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", fmt.Errorf("no plans directory at %s. Create it and add a plan template", c.PlansPath())
		}
		return "", err
	}

	switch len(ids) {
	case 0:
		return "", fmt.Errorf("no plan templates found in %s", c.PlansPath())
	case 1:
		return ids[0], nil
	}

	var id string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which training plan are you following?").
				Description("Saved to your config. Edit active_plan to change it.").
				Options(huh.NewOptions(ids...)...).
				Value(&id),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}

	return id, nil
}
