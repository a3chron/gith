package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kurtschambach/gith/internal"
	"github.com/kurtschambach/gith/internal/git"
)

type Step int

type ActionModel struct {
	Actions        []string
	SelectedAction string
}

type BranchModel struct {
	Actions        []string
	SelectedAction string
	Branches       []string
	SelectedBranch string
}

// type CommitModel struct
// type TagModel struct
// type RemoteModel struct
// type ChangesModel struct

type ConfigModel struct {
	Accent   string
	Flavor   string
	Flavores []string
	Accents  []string
}

type Model struct {
	CurrentStep Step
	Loading     bool
	Selected    int
	ActionModel ActionModel
	BranchModel BranchModel
	ConfigModel ConfigModel
	Spinner     spinner.Model
	Output      string
	Err         string
	Success     string
}

type repoUpdatedMsg struct{}

const (
	StepLoad   Step = iota
	StepAction      // select an action
	StepBranchAction
	StepBranchSelect
	StepCommit
	StepTag
	StepRemote
	StepChanges
	StepOptions // preferences
)

func UpdateOnInit() tea.Cmd {
	return func() tea.Msg {
		internal.UpdateRepo()
		return repoUpdatedMsg{}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		UpdateOnInit(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case repoUpdatedMsg:
		m.Loading = false
		m.CurrentStep = StepAction
		return m, nil

	case spinner.TickMsg:
		m.Spinner, _ = m.Spinner.Update(msg)
		return m, m.Spinner.Tick

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.Err = "User quit"
			return m, tea.Quit
		case "ctrl+h":
			// Go back to previous Step
			switch m.CurrentStep {
			case StepAction, StepLoad:
				return m, tea.Quit
			default:
				m.CurrentStep = StepAction
				m.ActionModel.SelectedAction = ""
				m.BranchModel.SelectedBranch = ""
				m.Selected = 0
				m.Err = ""
				m.Output = ""
				m.Success = ""
			}
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			} else {
				switch m.CurrentStep {
				case StepAction:
					m.Selected = len(m.ActionModel.Actions) - 1
				case StepBranchAction:
					m.Selected = len(m.BranchModel.Actions) - 1
				case StepBranchSelect:
					m.Selected = len(m.BranchModel.Branches) - 1
				}
				// TODO: add other steps
			}
		case "down", "j":
			switch m.CurrentStep {
			case StepAction:
				if m.Selected < len(m.ActionModel.Actions)-1 {
					m.Selected++
				} else {
					m.Selected = 0
				}
			case StepBranchAction:
				if m.Selected < len(m.BranchModel.Actions)-1 {
					m.Selected++
				} else {
					m.Selected = 0
				}
			case StepBranchSelect:
				if m.Selected < len(m.BranchModel.Branches)-1 {
					m.Selected++
				} else {
					m.Selected = 0
				}
			}
			// TODO: add other steps
		case "enter":
			switch m.CurrentStep {
			case StepAction:
				m.Output = ""
				m.ActionModel.SelectedAction = m.ActionModel.Actions[m.Selected]

				//"Branch", "Status", "Commit", "Tag", "Remote", "Changes", "Options"
				switch m.ActionModel.SelectedAction {
				case "Branch":
					m.Selected = 0
					m.CurrentStep = StepBranchAction
				case "Status":
					cmd := exec.Command("git", "status", "--porcelain")
					output, err := cmd.Output()
					if err != nil {
						m.Err = fmt.Sprintf("Failed to get status: %v", err)
					} else {
						if len(output) == 0 {
							m.Success = "Working tree clean"
						} else {
							lines := strings.Split(strings.TrimSpace(string(output)), "\n")
							m.Output += string(output) + "\n"
							m.Success = fmt.Sprintf("Changes detected: %d files", len(lines))
						}
					}
					return m, tea.Quit
				// TODO: add other actions
				case "Options":
					m.CurrentStep = StepOptions
					m.Err = "Options will come here soon"
					return m, tea.Quit
				}
			case StepBranchAction:
				m.BranchModel.SelectedAction = m.BranchModel.Actions[m.Selected]

				switch m.BranchModel.SelectedAction {
				case "Create branch":
					m.Err = "Create branch functionality not implemented yet"
					return m, tea.Quit

				case "Switch branch", "Delete branch":
					branches, err := git.GetBranches()
					if err != nil {
						m.Err = fmt.Sprintf("Failed to fetch branches: %v", err)
						return m, nil
					}
					if len(branches) == 0 {
						m.Err = "No branches available"
						return m, tea.Quit //TODO: nil or Quit?
					}

					m.BranchModel.Branches = branches

					m.Selected = 0
					m.CurrentStep = StepBranchSelect
				}

			case StepBranchSelect:
				m.BranchModel.SelectedBranch = m.BranchModel.Branches[m.Selected]

				switch m.BranchModel.SelectedAction {
				case "Switch branch":
					// First check if branch exists locally
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
				case "Delete branch":
					// Try to delete local branch first
					cmd := exec.Command("git", "branch", "-d", m.BranchModel.SelectedBranch)
					output, err := cmd.CombinedOutput()
					if err != nil {
						m.Err += string(output) + "\n" // TODO: m.ErrOut
						// If that fails, try force delete
						cmd = exec.Command("git", "branch", "-D", m.BranchModel.SelectedBranch)
						output, err = cmd.CombinedOutput()
						if err != nil {
							m.Err = fmt.Sprintf("Delete failed: %s", string(output))
						} else {
							m.Success = fmt.Sprintf("Force deleted branch '%s'", m.BranchModel.SelectedBranch)
						}
					} else {
						m.Output += string(output) + "\n"
						m.Success = fmt.Sprintf("Deleted branch '%s'", m.BranchModel.SelectedBranch)
					}
				}
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
