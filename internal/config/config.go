package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	catppuccin "github.com/catppuccin/go"
)

type Config struct {
	Accent string `json:"accent"`
	Flavor string `json:"flavor"`
}

var DefaultConfig = Config{
	Accent: "Blue",
	Flavor: "Mocha",
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "gith", "config.json"), nil
	}

	// Fall back to ~/.config/gith/config.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "gith", "config.json"), nil
}

// LoadConfig loads configuration from file, creating default if it doesn't exist
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(&DefaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return &DefaultConfig, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config values
	if !isValidFlavor(config.Flavor) {
		config.Flavor = DefaultConfig.Flavor
	}
	if !isValidAccent(config.Accent) {
		config.Accent = DefaultConfig.Accent
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCatppuccinFlavor returns the catppuccin flavor based on config
func GetCatppuccinFlavor(flavorName string) catppuccin.Flavor {
	switch flavorName {
	case "Latte":
		return catppuccin.Latte
	case "Frappe":
		return catppuccin.Frappe
	case "Macchiato":
		return catppuccin.Macchiato
	case "Mocha":
		return catppuccin.Mocha
	default:
		return catppuccin.Mocha
	}
}

// GetAccentColor returns the accent color based on config
func GetAccentColor(flavor catppuccin.Flavor, accentName string) catppuccin.Color {
	switch accentName {
	case "Rosewater":
		return flavor.Rosewater()
	case "Flamingo":
		return flavor.Flamingo()
	case "Pink":
		return flavor.Pink()
	case "Mauve":
		return flavor.Mauve()
	case "Red":
		return flavor.Red()
	case "Maroon":
		return flavor.Maroon()
	case "Peach":
		return flavor.Peach()
	case "Yellow":
		return flavor.Yellow()
	case "Green":
		return flavor.Green()
	case "Teal":
		return flavor.Teal()
	case "Blue":
		return flavor.Blue()
	case "Sapphire":
		return flavor.Sapphire()
	case "Sky":
		return flavor.Sky()
	case "Lavender":
		return flavor.Lavender()
	case "Gray":
		return flavor.Overlay2()
	default:
		return flavor.Blue()
	}
}

// Validation helpers
func isValidFlavor(flavor string) bool {
	validFlavors := []string{"Latte", "Frappe", "Macchiato", "Mocha"}
	return slices.Contains(validFlavors, flavor)
}

func isValidAccent(accent string) bool {
	validAccents := []string{
		"Rosewater", "Flamingo", "Pink", "Mauve", "Red", "Maroon",
		"Peach", "Yellow", "Green", "Teal", "Blue", "Sapphire",
		"Sky", "Lavender", "Gray",
	}
	return slices.Contains(validAccents, accent)
}

// GetAvailableFlavors returns list of available flavors
func GetAvailableFlavors() []string {
	return []string{"Latte", "Frappe", "Macchiato", "Mocha"}
}

// GetAvailableAccents returns list of available accents
func GetAvailableAccents() []string {
	return []string{
		"Rosewater", "Flamingo", "Pink", "Mauve", "Red", "Maroon",
		"Peach", "Yellow", "Green", "Teal", "Blue", "Sapphire",
		"Sky", "Lavender", "Gray",
	}
}
