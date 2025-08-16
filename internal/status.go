package internal

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) ExecuteStatus() (*Model, tea.Cmd) {
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
}
