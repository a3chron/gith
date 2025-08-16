package internal

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleActionSelection() (tea.Model, tea.Cmd) {
	m.ActionModel.SelectedAction = m.ActionModel.Actions[m.Selected]

	switch m.ActionModel.SelectedAction {
	case "Branch":
		m.Selected = 0
		m.CurrentStep = StepBranchAction
		m.Level = 2
	case "Status":
		return m.ExecuteStatus()
	case "Commit":
		m.Selected = 0
		m.CurrentStep = StepCommitAction
		m.Level = 2
	case "Tag":
		m.Selected = 0
		m.CurrentStep = StepTag
		m.Level = 2
	case "Remote":
		m.Selected = 0
		m.CurrentStep = StepRemote
		m.Level = 2
	case "Changes":
		m.Selected = 0
		m.CurrentStep = StepChanges
		m.Level = 2
		m.Err = "Changes will come here soon"
		return m, tea.Quit
	case "Options":
		m.Selected = 0
		m.CurrentStep = StepOptions
		m.Level = 2
	}
	return m, nil
}
