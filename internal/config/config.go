package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

const (
	defaultPlansDir      = "plans"
	defaultTrainingsDir  = "trainings"
	defaultDailyNotesDir = "Diario"
)

type Config struct {
	VaultPath     string `mapstructure:"vault_path"`
	PlansDir      string `mapstructure:"plans_dir"`
	TrainingsDir  string `mapstructure:"trainings_dir"`
	DailyNotesDir string `mapstructure:"daily_notes_dir"`
	ActivePlan    string `mapstructure:"active_plan"`
}

func (c *Config) PlansPath() string {
	return filepath.Join(c.VaultPath, c.PlansDir)
}

func (c *Config) TrainingsPath() string {
	return filepath.Join(c.VaultPath, c.TrainingsDir)
}

func (c *Config) DailyNotesPath() string {
	return filepath.Join(c.VaultPath, c.DailyNotesDir)
}

// Load reads the config file from disk. If it doesn't exist yet, it runs
// a short onboarding form to create it before returning.
func Load() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := createConfig(path); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("checking config file: %w", err)
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.VaultPath == "" {
		return nil, fmt.Errorf("vault_path is missing from %s", path)
	}

	return &cfg, nil
}

func SetActivePlan(id string) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	viper.Set("active_plan", id)
	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("plans_dir", defaultPlansDir)
	viper.SetDefault("trainings_dir", defaultTrainingsDir)
	viper.SetDefault("daily_notes_dir", defaultDailyNotesDir)
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, ".config", "lorren", "config.yaml"), nil
}

func createConfig(path string) error {
	var vaultPath string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Where is your vault?").
				Description("Full path, or ~/... from your home folder.").
				Value(&vaultPath),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	vaultPath, err := expandHome(vaultPath)
	if err != nil {
		return err
	}

	info, err := os.Stat(vaultPath)
	if err != nil {
		return fmt.Errorf("vault path %q does not exist. Create the folder first, then run lorren again", vaultPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("vault path %q is not a directory", vaultPath)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	viper.Set("vault_path", vaultPath)
	viper.Set("plans_dir", defaultPlansDir)
	viper.Set("trainings_dir", defaultTrainingsDir)
	viper.Set("daily_notes_dir", defaultDailyNotesDir)

	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

func expandHome(path string) (string, error) {
	if path == "" || path[0] != '~' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}

	if path == "~" {
		return home, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:]), nil
	}

	return path, nil
}
