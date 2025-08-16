package internal

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

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
		return m.SwitchBranch()
	case "Delete Branch":
		return m.DeleteBranch()
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

	return m.CreateBranch(strings.TrimSpace(m.BranchModel.Input))
}

func (m *Model) CreateBranch(branchName string) (*Model, tea.Cmd) {
	out, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()
	if err != nil {
		m.OutputByLevel(string(out))
		m.Err = fmt.Sprintf("Failed to create branch: %s", string(out))
	} else {
		m.Success = fmt.Sprintf("Successfully created branch '%s'", branchName)
	}
	return m, tea.Quit
}

func (m *Model) SwitchBranch() (*Model, tea.Cmd) {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+m.BranchModel.SelectedBranch)
	err := cmd.Run()

	var switchCmd *exec.Cmd
	if err != nil {
		// Branch doesn't exist locally, create and track remote
		switchCmd = exec.Command("git", "checkout", "-b", m.BranchModel.SelectedBranch, "origin/"+m.BranchModel.SelectedBranch)
	} else {
		// Branch exists locally, just switch
		switchCmd = exec.Command("git", "switch", m.BranchModel.SelectedBranch)
	}

	output, err := switchCmd.CombinedOutput()
	if err != nil {
		m.Err = fmt.Sprintf("Switch failed: %s", string(output))
	} else {
		m.Success = fmt.Sprintf("Switched to branch '%s'", m.BranchModel.SelectedBranch)
	}
	return m, tea.Quit
}

func (m *Model) DeleteBranch() (*Model, tea.Cmd) {
	// Try to delete local branch first
	cmd := exec.Command("git", "branch", "-d", m.BranchModel.SelectedBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.OutputByLevel(string(output) + "\n")
		// suggest force delete:
		m.OutputByLevel("\\cpTo force delete use: \ngit branch -D" + m.BranchModel.SelectedBranch)
		m.Err = "Delete failed"
	} else {
		m.OutputByLevel(string(output) + "\n")
		m.Success = fmt.Sprintf("Deleted branch '%s'", m.BranchModel.SelectedBranch)
	}
	return m, tea.Quit
}

func (m Model) HandleBranchActionSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedAction = m.BranchModel.Actions[m.Selected]
	return m.HandleBranchOperation()
}

func (m Model) HandleBranchSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedBranch = m.BranchModel.Branches[m.Selected]
	return m.ExecuteBranchAction()
}
