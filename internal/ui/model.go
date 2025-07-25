package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kurtschambach/gith/internal"
	git "github.com/kurtschambach/gith/internal/git"
)

type Step int

type Model struct {
	CurrentStep    Step
	SelectedAction string
	SelectedBranch string
	Options        []string
	Branches       []string
	Accent         string
	Selected       int
	Spinner        spinner.Model
	Loading        bool
	Output         string
	Err            string
	Success        string
}

type repoUpdatedMsg struct{}

const (
	StepLoad Step = iota
	StepAction
	StepBranch
	StepOptions
	StepComplete
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
			case StepBranch, StepComplete:
				m.CurrentStep = StepAction
				m.SelectedAction = ""
				m.SelectedBranch = ""
				m.Selected = 0
				m.Err = ""
				m.Output = ""
				m.Success = ""
				m.Branches = nil
			default:
				return m, tea.Quit
			}
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			} else {
				switch m.CurrentStep {
				case StepAction:
					m.Selected = len(m.Options) - 1
				case StepBranch:
					m.Selected = len(m.Branches) - 1
				}
			}
		case "down", "j":
			switch m.CurrentStep {
			case StepAction:
				if m.Selected < len(m.Options)-1 {
					m.Selected++
				} else {
					m.Selected = 0
				}
			case StepBranch:
				if m.Selected < len(m.Branches)-1 {
					m.Selected++
				} else {
					m.Selected = 0
				}
			}
		case "enter":
			switch m.CurrentStep {
			case StepAction:
				m.Output = ""
				m.SelectedAction = m.Options[m.Selected]

				switch m.SelectedAction {
				case "Switch Branch", "Delete Branch":
					branches, err := git.GetBranches()
					if err != nil {
						m.Err = fmt.Sprintf("Failed to fetch branches: %v", err)
						m.CurrentStep = StepComplete
						return m, nil
					}
					if len(branches) == 0 {
						m.Err = "No branches available"
						m.CurrentStep = StepComplete
						return m, nil
					}
					m.Branches = branches
					m.Selected = 0
					m.CurrentStep = StepBranch
					m.Err = ""
				case "Create Branch":
					m.Err = "Create branch functionality not implemented yet"
					return m, tea.Quit
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
				case "Options":
					m.CurrentStep = StepOptions
					m.Success = "End"
					return m, tea.Quit
				}
			case StepBranch:
				m.Output = ""

				if len(m.Branches) > 0 && m.Selected < len(m.Branches) {
					m.SelectedBranch = m.Branches[m.Selected]

					switch m.SelectedAction {
					case "Switch Branch":
						// First check if branch exists locally
						cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+m.SelectedBranch)
						err := cmd.Run()

						var switchCmd *exec.Cmd
						if err != nil {
							// Branch doesn't exist locally, create and track remote
							switchCmd = exec.Command("git", "checkout", "-b", m.SelectedBranch, "origin/"+m.SelectedBranch)
						} else {
							// Branch exists locally, just switch
							switchCmd = exec.Command("git", "checkout", m.SelectedBranch)
						}

						output, err := switchCmd.CombinedOutput()
						if err != nil {
							m.Err = fmt.Sprintf("Switch failed: %s", string(output))
						} else {
							m.Success = fmt.Sprintf("Switched to branch '%s'", m.SelectedBranch)
						}
					case "Delete Branch":
						// Try to delete local branch first
						cmd := exec.Command("git", "branch", "-d", m.SelectedBranch)
						output, err := cmd.CombinedOutput()
						if err != nil {
							m.Err += string(output) + "\n" // TODO: m.ErrOut
							// If that fails, try force delete
							cmd = exec.Command("git", "branch", "-D", m.SelectedBranch)
							output, err = cmd.CombinedOutput()
							if err != nil {
								m.Err = fmt.Sprintf("Delete failed: %s", string(output))
							} else {
								m.Success = fmt.Sprintf("Force deleted branch '%s'", m.SelectedBranch)
							}
						} else {
							m.Output += string(output) + "\n"
							m.Success = fmt.Sprintf("Deleted branch '%s'", m.SelectedBranch)
						}
					}
					return m, tea.Quit
				}
			case StepComplete:
				// Reset to beginning
				m.CurrentStep = StepAction
				m.SelectedAction = ""
				m.SelectedBranch = ""
				m.Selected = 0
				m.Err = ""
				m.Output = ""
				m.Success = ""
				m.Branches = nil
			}
		}
	}
	return m, nil
}
