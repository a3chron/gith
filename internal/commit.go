package internal

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

// function to handle commit message submission
func (m *Model) HandleCommitMessageSubmit() (*Model, tea.Cmd) {
	if strings.TrimSpace(m.CommitModel.CommitMessage) == "" {
		m.Err = "Commit message cannot be empty"
		return m, tea.Quit
	}

	var out string
	var err error

	switch m.CommitModel.SelectedAction {
	case "Commit Staged":
		out, err = git.CommitStaged(m.CommitModel.CommitMessage)
	case "Commit All":
		out, err = git.CommitAll(m.CommitModel.CommitMessage)
	}

	if err != nil {
		m.OutputByLevel("\\ceError:\n" + out)
		m.Err = "Failed to Commit"
	} else {
		m.Success = "Commited Changes"
	}

	return m, tea.Quit
}

func (m Model) HandleCommitSelection() (tea.Model, tea.Cmd) {
	m.CommitModel.SelectedAction = m.CommitModel.Actions[m.Selected]

	if m.CommitModel.SelectedAction == "Undo Last Commit" {
		out, err := git.UndoLastCommit()
		m.OutputByLevel(out)
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Undo Commit Successful"
		}
		return m, tea.Quit
	}

	m.Level = 3
	m.Selected = 0
	m.CurrentStep = StepCommitSelectPrefix
	return m, nil
}

func (m Model) HandleCommitPrefixSelection() (tea.Model, tea.Cmd) {
	m.CommitModel.SelectedPrefix = m.CommitModel.CommitPrefixes[m.Selected]

	if m.CommitModel.SelectedPrefix != "Custom Prefix" {
		m.CommitModel.CommitMessage = m.CommitModel.SelectedPrefix + ": "
	}

	m.Selected = 0
	m.CurrentStep = StepCommitInput
	return m, nil
}
