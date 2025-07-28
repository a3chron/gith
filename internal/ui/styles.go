package ui

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	mocha = catppuccin.Mocha

	// Styles using Catppuccin Mocha palette
	AccentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Blue().Hex))

	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Text().Hex))

	TitleStyle = AccentStyle.Bold(true)

	LogoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Text().Hex)).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Text().Hex)).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Subtext0().Hex))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Overlay0().Hex))

	CompletedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Subtext0().Hex))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Red().Hex))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Green().Hex))

	BulletStyle = AccentStyle.Bold(true)

	LineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Surface1().Hex))

	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	GreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(mocha.Green().Hex))

	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(mocha.Yellow().Hex))

	RedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(mocha.Red().Hex))
)
