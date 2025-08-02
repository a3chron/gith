package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	if !IsValidFlavor(config.Flavor) {
		config.Flavor = DefaultConfig.Flavor
	}
	if !IsValidAccent(config.Accent) {
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
	switch strings.ToLower(flavorName) {
	case "latte":
		return catppuccin.Latte
	case "frappe":
		return catppuccin.Frappe
	case "macchiato":
		return catppuccin.Macchiato
	case "mocha":
		return catppuccin.Mocha
	default:
		return catppuccin.Mocha
	}
}

// GetAccentColor returns the accent color based on config
func GetAccentColor(flavor catppuccin.Flavor, accentName string) catppuccin.Color {
	switch strings.ToLower(accentName) {
	case "rosewater":
		return flavor.Rosewater()
	case "flamingo":
		return flavor.Flamingo()
	case "pink":
		return flavor.Pink()
	case "mauve":
		return flavor.Mauve()
	case "red":
		return flavor.Red()
	case "maroon":
		return flavor.Maroon()
	case "peach":
		return flavor.Peach()
	case "yellow":
		return flavor.Yellow()
	case "green":
		return flavor.Green()
	case "teal":
		return flavor.Teal()
	case "blue":
		return flavor.Blue()
	case "sapphire":
		return flavor.Sapphire()
	case "sky":
		return flavor.Sky()
	case "lavender":
		return flavor.Lavender()
	case "gray":
		return flavor.Overlay2()
	default:
		return flavor.Blue()
	}
}

// Validation helpers
func IsValidFlavor(flavor string) bool {
	lowerFlavor := strings.ToLower(flavor)
	validFlavors := []string{"latte", "frappe", "macchiato", "mocha"}
	return slices.Contains(validFlavors, lowerFlavor)
}

func IsValidAccent(accent string) bool {
	lowerAccent := strings.ToLower(accent)
	validAccents := []string{
		"rosewater", "flamingo", "pink", "mauve", "red", "maroon",
		"peach", "yellow", "green", "teal", "blue", "sapphire",
		"sky", "lavender", "gray",
	}
	return slices.Contains(validAccents, lowerAccent)
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
