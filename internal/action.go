package internal

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/git"
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
		status, length, err := git.GetStatusInfo()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to get status: %v", err)
		} else {
			if length == 0 {
				m.Success = "Working tree clean"
			} else {
				if len(status["modified"]) > 0 {
					m.OutputByLevel("\\cyModified:\n " + strings.Join(status["modified"], "\n ") + "\n╌\n")
				}
				if len(status["added"]) > 0 {
					m.OutputByLevel("\\cgAdded:\n " + strings.Join(status["added"], "\n ") + "\n╌\n")
				}
				if len(status["deleted"]) > 0 {
					m.OutputByLevel("\\crDeleted:\n " + strings.Join(status["deleted"], "\n ") + "\n╌\n")
				}
				if len(status["untracked"]) > 0 {
					m.OutputByLevel("\\ctUntracked:\n " + strings.Join(status["untracked"], "\n ") + "\n╌\n")
				}
				m.Success = fmt.Sprintf("Changes detected: %d %s", length, func() string {
					if length == 1 {
						return "file"
					}
					return "files"
				}())
			}
		}
		return m, tea.Quit
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
