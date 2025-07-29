package ui

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	flavor = catppuccin.Mocha

	// Styles using Catppuccin Mocha palette
	AccentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(flavor.Blue().Hex))

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

	GreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Green().Hex))

	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Yellow().Hex))

	RedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Red().Hex))
)
