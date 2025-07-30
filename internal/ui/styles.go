package ui

import (
	"github.com/a3chron/gith/internal/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Global style variables - will be updated by UpdateStyles
	AccentStyle    lipgloss.Style
	TextStyle      lipgloss.Style
	TitleStyle     lipgloss.Style
	LogoStyle      lipgloss.Style
	SelectedStyle  lipgloss.Style
	NormalStyle    lipgloss.Style
	DimStyle       lipgloss.Style
	CompletedStyle lipgloss.Style
	ErrorStyle     lipgloss.Style
	SuccessStyle   lipgloss.Style
	BulletStyle    lipgloss.Style
	LineStyle      lipgloss.Style
	ContainerStyle lipgloss.Style
	GreenStyle     lipgloss.Style
	YellowStyle    lipgloss.Style
	RedStyle       lipgloss.Style
)

// UpdateStyles updates all styles based on the provided config
func UpdateStyles(cfg *config.Config) {
	flavor := config.GetCatppuccinFlavor(cfg.Flavor)
	accentColor := config.GetAccentColor(flavor, cfg.Accent)

	// Update all styles
	AccentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentColor.Hex))

	TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Text().Hex))

	TitleStyle = AccentStyle.Bold(true)

	LogoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Text().Hex)).
		Bold(true)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Text().Hex)).
		Bold(true)

	NormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Subtext0().Hex))

	DimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Overlay0().Hex))

	CompletedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Subtext0().Hex))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Red().Hex))

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Green().Hex))

	BulletStyle = AccentStyle.Bold(true)

	LineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Surface1().Hex))

	ContainerStyle = lipgloss.NewStyle().
		Padding(1, 2)

	GreenStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Green().Hex))

	YellowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Yellow().Hex))

	RedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(flavor.Red().Hex))
}

// InitializeStyles initializes styles with default or loaded config
func InitializeStyles() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fall back to default config if loading fails
		cfg = &config.DefaultConfig
	}
	UpdateStyles(cfg)
	return nil
}
