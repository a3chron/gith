package internal

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/a3chron/gith/internal/config"
	"github.com/a3chron/gith/internal/ui"
)

func UpdateOnInit() tea.Cmd {
	return func() tea.Msg {
		if err := UpdateRepo(); err != nil {
			return RepoUpdateErrorMsg{}
		}
		return RepoUpdatedMsg{}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		UpdateOnInit(),
	)
}

func (m Model) getCurrentOptions() []string {
	switch m.CurrentStep {
	case StepAction:
		return m.ActionModel.Actions

	case StepBranchAction:
		return m.BranchModel.Actions
	case StepBranchSelect:
		return m.BranchModel.Branches
	case StepBranchCreate:
		return m.BranchModel.Options

	case StepCommitAction:
		return m.CommitModel.Actions
	case StepCommitSelectPrefix:
		return m.CommitModel.CommitPrefixes

	case StepTag:
		return m.TagModel.Actions
	case StepTagSelect:
		return m.TagModel.Options
	case StepTagAdd:
		return m.TagModel.AddOptions

	case StepRemote:
		return m.RemoteModel.Actions
	case StepRemoteSelect:
		return m.RemoteModel.Options

	case StepOptions:
		return m.ConfigModel.Actions
	case StepOptionsFlavorSelect:
		return m.ConfigModel.Flavors
	case StepOptionsAccentSelect:
		return m.ConfigModel.Accents
	default:
		return []string{}
	}
}

func (m *Model) resetState() {
	m.CurrentStep = StepAction
	m.ActionModel.SelectedAction = ""

	m.BranchModel.SelectedBranch = ""
	m.BranchModel.SelectedAction = ""
	m.BranchModel.SelectedOption = ""
	m.BranchModel.Input = ""

	m.CommitModel.SelectedAction = ""
	m.CommitModel.SelectedPrefix = ""
	m.CommitModel.CommitMessage = ""

	m.TagModel.SelectedAction = ""
	m.TagModel.SelectedOption = ""
	m.TagModel.SelectedAddTag = ""
	m.TagModel.CurrentTag = ""
	m.TagModel.ManualInput = ""

	m.RemoteModel.SelectedAction = ""
	m.RemoteModel.SelectedOption = ""
	m.RemoteModel.NameInput = ""
	m.RemoteModel.UrlInput = ""

	m.ConfigModel.SelectedAccent = ""
	m.ConfigModel.SelectedFlavor = ""
	m.ConfigModel.SelectedAction = ""
	m.Selected = 0
	m.Level = 1
	m.Output = []string{} // TODO add m.Output[0] if exisiting
	m.Err = ""
	m.Success = ""
	m.StartAt = ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RepoUpdatedMsg, RepoUpdateErrorMsg:
		m.Loading = false

		if _, isErr := msg.(RepoUpdateErrorMsg); isErr {
			m.OutputByLevel("Failed to fetch from remote")
		}

		// quick-start flows triggered via StartAt
		switch m.StartAt {
		case "add-tag":
			m.PrepareTagAddition()
			m.ActionModel.SelectedAction = "Tag"
			m.TagModel.SelectedAction = "Add Tag"

		case "push-tag":
			m.PrepareTagSelection()
			m.ActionModel.SelectedAction = "Tag"
			m.TagModel.SelectedAction = "Push Tag"

		case "commit-staged":
			m.Level = 3
			m.Selected = 0
			m.CurrentStep = StepCommitSelectPrefix
			m.ActionModel.SelectedAction = "Commit"
			m.CommitModel.SelectedAction = "Commit Staged"

		case "commit-all":
			m.Level = 3
			m.Selected = 0
			m.CurrentStep = StepCommitSelectPrefix
			m.ActionModel.SelectedAction = "Commit"
			m.CommitModel.SelectedAction = "Commit All"

		default:
			m.CurrentStep = StepAction
			m.Level = 1
		}

		return m, nil

	case spinner.TickMsg:
		m.Spinner, _ = m.Spinner.Update(msg)
		return m, m.Spinner.Tick

	case tea.KeyMsg:
		if isInputStep(m.CurrentStep) {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.Err = "User cancelled"
				return m, tea.Quit
			case "ctrl+h":
				m.resetState()
			case "enter":
				switch m.CurrentStep {
				case StepTagInput:
					return m.HandleTagInputSubmit()
				case StepBranchInput:
					return m.HandleBranchInputSubmit()
				case StepRemoteNameInput:
					// proceed to URL input if name provided
					if strings.TrimSpace(m.RemoteModel.NameInput) == "" {
						m.Err = "Remote name cannot be empty"
						return m, tea.Quit
					}
					m.CurrentStep = StepRemoteUrlInput
					return m, nil
				case StepRemoteUrlInput:
					if strings.TrimSpace(m.RemoteModel.UrlInput) == "" {
						m.Err = "Remote URL cannot be empty"
						return m, tea.Quit
					}
					out, err := exec.Command("git", "remote", "add", m.RemoteModel.NameInput, m.RemoteModel.UrlInput).CombinedOutput()
					if err != nil {
						m.OutputByLevel(string(out))
						m.Err = "Failed to add remote"
					} else {
						m.Success = fmt.Sprintf("Added Remote '%s' -> %s", m.RemoteModel.NameInput, m.RemoteModel.UrlInput)
					}
					return m, tea.Quit
				case StepCommitInput:
					return m.HandleCommitMessageSubmit()
				}
			case "backspace":
				switch m.CurrentStep {
				case StepTagInput:
					if len(m.TagModel.ManualInput) > 0 {
						m.TagModel.ManualInput = m.TagModel.ManualInput[:len(m.TagModel.ManualInput)-1]
					}
				case StepBranchInput:
					if len(m.BranchModel.Input) > 0 {
						m.BranchModel.Input = m.BranchModel.Input[:len(m.BranchModel.Input)-1]
					}
				case StepRemoteNameInput:
					if len(m.RemoteModel.NameInput) > 0 {
						m.RemoteModel.NameInput = m.RemoteModel.NameInput[:len(m.RemoteModel.NameInput)-1]
					}
				case StepRemoteUrlInput:
					if len(m.RemoteModel.UrlInput) > 0 {
						m.RemoteModel.UrlInput = m.RemoteModel.UrlInput[:len(m.RemoteModel.UrlInput)-1]
					}
				case StepCommitInput:
					if len(m.CommitModel.CommitMessage) > 0 {
						m.CommitModel.CommitMessage = m.CommitModel.CommitMessage[:len(m.CommitModel.CommitMessage)-1]
					}
				}
			default:
				// Add character to input
				if len(msg.String()) == 1 {
					switch m.CurrentStep {
					case StepTagInput:
						m.TagModel.ManualInput += msg.String()
					case StepBranchInput:
						m.BranchModel.Input += msg.String()
					case StepRemoteNameInput:
						m.RemoteModel.NameInput += msg.String()
					case StepRemoteUrlInput:
						m.RemoteModel.UrlInput += msg.String()
					case StepCommitInput:
						m.CommitModel.CommitMessage += msg.String()
					}
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.Err = "User quit"
			return m, tea.Quit

		case "ctrl+h":
			if m.CurrentStep == StepAction || m.CurrentStep == StepLoad {
				return m, tea.Quit
			}
			m.resetState()

		case "up", "k", "down", "j":
			m.handleNavigation(msg.String())

		case "enter":
			return m.handleEnterKey()
		}
	}
	return m, nil
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.CurrentStep {
	case StepAction:
		return m.HandleActionSelection()

	case StepBranchAction:
		return m.HandleBranchActionSelection()
	case StepBranchSelect:
		return m.HandleBranchSelection()
	case StepBranchCreate:
		return m.HandleBranchCreateSelection()

	case StepCommitAction:
		return m.HandleCommitSelection()
	case StepCommitSelectPrefix:
		return m.HandleCommitPrefixSelection()

	case StepTag:
		return m.HandleTagActionSelection()
	case StepTagSelect:
		return m.HandleTagSelection()
	case StepTagAdd:
		return m.HandleTagAddSelection()

	case StepRemote:
		return m.HandleRemoteActionSelection()
	case StepRemoteSelect:
		return m.HandleRemoteSelection()

	case StepOptions:
		return m.HandleOptionsActionSelection()
	case StepOptionsFlavorSelect:
		return m.HandleOptionsFlavorSelection()
	case StepOptionsAccentSelect:
		return m.HandleOptionsAccentSelection()
	}
	return m, nil
}

func (m *Model) handleNavigation(key string) {
	options := m.getCurrentOptions()
	if len(options) == 0 {
		return
	}

	switch key {
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		} else {
			m.Selected = len(options) - 1
		}
	case "down", "j":
		if m.Selected < len(options)-1 {
			m.Selected++
		} else {
			m.Selected = 0
		}
	}

	// handle preview when selecting a theme
	if m.CurrentStep == StepOptionsAccentSelect {
		flavor := config.GetCatppuccinFlavor(m.CurrentConfig.Flavor)
		ui.UpdateStyles(flavor, config.GetAccentColor(flavor, m.ConfigModel.Accents[m.Selected]))
	}

	if m.CurrentStep == StepOptionsFlavorSelect {
		flavor := config.GetCatppuccinFlavor(m.ConfigModel.Flavors[m.Selected])
		ui.UpdateStyles(flavor, config.GetAccentColor(flavor, m.CurrentConfig.Accent))
	}
}
