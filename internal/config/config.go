package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

type Config struct {
	VaultPath string `mapstructure:"vault_path"`
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

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &cfg, nil
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
				Title("Where is your Obsidian vault?").
				Description("Enter the fult path to an existing folder.").
				Value(&vaultPath),
		),
	)

	if err := form.Run(); err != nil {
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
	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
