package internal

import (
	"fmt"

	"github.com/a3chron/gith/internal/config"
	"github.com/a3chron/gith/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleOptionsActionSelection() (tea.Model, tea.Cmd) {
	// Ensure CurrentConfig is not nil
	if m.CurrentConfig == nil {
		cfg, err := config.LoadConfig()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to load config: %v", err)
			return m, tea.Quit
		}
		m.CurrentConfig = cfg
	}

	m.ConfigModel.SelectedAction = m.ConfigModel.Actions[m.Selected]

	switch m.ConfigModel.SelectedAction {
	case "Change Flavor":
		m.ConfigModel.Flavors = config.GetAvailableFlavors()
		m.Selected = 0
		// Set current selection to match current config
		for i, flavor := range m.ConfigModel.Flavors {
			if flavor == m.CurrentConfig.Flavor {
				m.Selected = i
				break
			}
		}
		m.CurrentStep = StepOptionsFlavorSelect
		m.Level = 3

	case "Change Accent":
		m.ConfigModel.Accents = config.GetAvailableAccents()
		m.Selected = 0
		// Set current selection to match current config
		for i, accent := range m.ConfigModel.Accents {
			if accent == m.CurrentConfig.Accent {
				m.Selected = i
				break
			}
		}
		m.CurrentStep = StepOptionsAccentSelect
		m.Level = 3

	case "Fetch at Init Behaviour":
		m.Selected = 0
		// Set current selection to match current config
		for i, behaviour := range m.ConfigModel.InitBehaviours {
			if behaviour == m.CurrentConfig.InitBehaviour {
				m.Selected = i
				break
			}
		}
		m.CurrentStep = StepOptionsInitBehaviourSelect
		m.Level = 3

	case "Reset to Defaults":
		m.CurrentConfig = &config.Config{
			Accent:        config.DefaultConfig.Accent,
			Flavor:        config.DefaultConfig.Flavor,
			InitBehaviour: config.DefaultConfig.InitBehaviour,
		}
		if err := config.SaveConfig(m.CurrentConfig); err != nil {
			m.Err = fmt.Sprintf("Failed to save config: %v", err)
		} else {
			ui.UpdateStylesByConfig(m.CurrentConfig)
			m.Success = "Configuration reset to defaults"
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) HandleOptionsFlavorSelection() (tea.Model, tea.Cmd) {
	// Ensure CurrentConfig is not nil
	if m.CurrentConfig == nil {
		cfg, err := config.LoadConfig()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to load config: %v", err)
			return m, tea.Quit
		}
		m.CurrentConfig = cfg
	}

	selectedFlavor := m.ConfigModel.Flavors[m.Selected]
	m.ConfigModel.SelectedFlavor = selectedFlavor

	// Update config and save
	m.CurrentConfig.Flavor = selectedFlavor
	if err := config.SaveConfig(m.CurrentConfig); err != nil {
		m.Err = fmt.Sprintf("Failed to save config: %v", err)
	} else {
		ui.UpdateStylesByConfig(m.CurrentConfig)
		m.Success = fmt.Sprintf("Flavor changed to %s", selectedFlavor)
	}
	return m, tea.Quit
}

func (m Model) HandleOptionsAccentSelection() (tea.Model, tea.Cmd) {
	// Ensure CurrentConfig is not nil
	if m.CurrentConfig == nil {
		cfg, err := config.LoadConfig()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to load config: %v", err)
			return m, tea.Quit
		}
		m.CurrentConfig = cfg
	}

	selectedAccent := m.ConfigModel.Accents[m.Selected]
	m.ConfigModel.SelectedAccent = selectedAccent

	// Update config and save
	m.CurrentConfig.Accent = selectedAccent
	if err := config.SaveConfig(m.CurrentConfig); err != nil {
		m.Err = fmt.Sprintf("Failed to save config: %v", err)
	} else {
		ui.UpdateStylesByConfig(m.CurrentConfig)
		m.Success = fmt.Sprintf("Accent changed to %s", selectedAccent)
	}
	return m, tea.Quit
}

func (m Model) HandleOptionsInitBehaviourSelection() (tea.Model, tea.Cmd) {
	// Ensure CurrentConfig is not nil
	if m.CurrentConfig == nil {
		cfg, err := config.LoadConfig()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to load config: %v", err)
			return m, tea.Quit
		}
		m.CurrentConfig = cfg
	}

	m.ConfigModel.SelectedBehaviour = m.ConfigModel.InitBehaviours[m.Selected]

	// Update config and save
	m.CurrentConfig.InitBehaviour = m.ConfigModel.SelectedBehaviour
	if err := config.SaveConfig(m.CurrentConfig); err != nil {
		m.Err = fmt.Sprintf("Failed to save config: %v", err)
	} else {
		m.Success = fmt.Sprintf("Init Behaviour changed to %s", m.ConfigModel.SelectedBehaviour)
	}
	return m, tea.Quit
}
