package internal

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleBranchActionSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedAction = m.BranchModel.Actions[m.Selected]
	return m.HandleBranchOperation()
}

func (m Model) HandleBranchSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedBranch = m.BranchModel.Branches[m.Selected]
	return m.ExecuteBranchAction()
}

func (m *Model) HandleBranchOperation() (*Model, tea.Cmd) {
	switch m.BranchModel.SelectedAction {
	case "Create Branch":
		return m.PrepareBranchAddition()

	case "Switch Branch", "Delete Branch":
		branches, err := git.GetBranches()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to fetch branches: %v", err)
			return m, nil
		}
		if len(branches) == 0 {
			m.Err = "No branches available"
			return m, tea.Quit
		}

		m.BranchModel.Branches = branches
		m.Selected = 0
		m.CurrentStep = StepBranchSelect
		m.Level = 3

	case "List Branches":
		branches, err := git.GetBranches()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to fetch branches: %v", err)
			return m, nil
		}
		if len(branches) == 0 {
			m.Success = "No branches available"
			return m, tea.Quit
		}

		m.OutputByLevel(strings.Join(branches, "\n"))
		m.Success = "Listed branches"
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) PrepareBranchAddition() (*Model, tea.Cmd) {
	m.Selected = 0
	m.CurrentStep = StepBranchCreate
	m.Level = 3
	return m, nil
}

func (m *Model) ExecuteBranchAction() (*Model, tea.Cmd) {
	switch m.BranchModel.SelectedAction {
	case "Switch Branch":
		out, err := git.SwitchBranch(m.BranchModel.SelectedBranch)

		m.OutputByLevel(out)

		if err != nil {
			m.Err = "Failed to Switch Branch"
		} else {
			m.Success = "Switched Branch"
		}

		return m, tea.Quit

	case "Delete Branch":
		out, err := git.DeleteBranch(m.BranchModel.SelectedBranch)

		m.OutputByLevel(out)

		if err != nil {
			m.Err = "Failed to Delete Branch"
		} else {
			m.Success = "Deleted Branch"
		}

		return m, tea.Quit
	}
	return m, tea.Quit
}

func (m Model) HandleBranchCreateSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedOption = m.BranchModel.Options[m.Selected]

	if m.BranchModel.SelectedOption == "Manual Input" {
		// Switch to input mode
		m.BranchModel.Input = ""
	} else {
		// Add prefix to input & switch to input mode
		m.BranchModel.Input = m.BranchModel.SelectedOption
	}

	m.Selected = 0
	m.CurrentStep = StepBranchInput
	return m, nil
}

// function to handle manual branch input submission
func (m *Model) HandleBranchInputSubmit() (*Model, tea.Cmd) {
	if strings.TrimSpace(m.BranchModel.Input) == "" {
		m.Err = "Branch name cannot be empty"
		return m, tea.Quit
	}

	out, err := git.CreateBranch(strings.TrimSpace(m.BranchModel.Input))

	m.OutputByLevel(out)

	if err != nil {
		m.Err = "Failed to Create Branch"
	} else {
		m.Success = "Created Branch"
	}

	return m, tea.Quit
}
