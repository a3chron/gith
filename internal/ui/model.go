package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/a3chron/gith/internal"
	"github.com/a3chron/gith/internal/config"
	"github.com/a3chron/gith/internal/git"
)

type Step int

const (
	StepLoad Step = iota
	StepAction

	StepBranchAction
	StepBranchSelect
	StepBranchCreate
	StepBranchInput

	StepCommitAction
	StepCommitSelectPrefix
	StepCommitInput

	StepTag
	StepTagSelect
	StepTagAdd
	StepTagInput

	StepRemote
	StepRemoteSelect
	StepRemoteNameInput
	StepRemoteUrlInput

	StepChanges

	StepOptions
	StepOptionsFlavorSelect
	StepOptionsAccentSelect
)

type ActionModel struct {
	Actions        []string
	SelectedAction string
}

type BranchModel struct {
	Actions        []string
	SelectedAction string
	Branches       []string
	SelectedBranch string
	Options        []string
	SelectedOption string
	Input          string
}

type CommitModel struct {
	Actions        []string
	SelectedAction string
	CommitPrefixes []string
	SelectedPrefix string
	CommitMessage  string
}

type TagModel struct {
	Actions        []string
	SelectedAction string
	Options        []string
	SelectedOption string
	AddOptions     []string
	SelectedAddTag string
	CurrentTag     string
	ManualInput    string
}

type RemoteModel struct {
	Actions        []string
	SelectedAction string
	Options        []string
	SelectedOption string
	NameInput      string
	UrlInput       string
}

type ConfigModel struct {
	Actions        []string
	SelectedAction string
	Flavors        []string
	SelectedFlavor string
	Accents        []string
	SelectedAccent string
}

type Model struct {
	CurrentStep   Step
	Loading       bool
	Selected      int
	ActionModel   ActionModel
	BranchModel   BranchModel
	CommitModel   CommitModel
	RemoteModel   RemoteModel
	TagModel      TagModel
	ConfigModel   ConfigModel
	CurrentConfig *config.Config
	Spinner       spinner.Model
	Level         int
	Output        []string
	Err           string
	Success       string
	StartAt       string // allows jumping directly to a specific flow when the UI loads
	StartAtLevel  int
}

type repoUpdatedMsg struct{}
type repoUpdateErrorMsg struct{}

func UpdateOnInit() tea.Cmd {
	return func() tea.Msg {
		if err := internal.UpdateRepo(); err != nil {
			return repoUpdateErrorMsg{}
		}
		return repoUpdatedMsg{}
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

// isInputStep returns true if the current step expects free-text input
func isInputStep(step Step) bool {
	switch step {
	case StepTagInput, StepBranchInput, StepRemoteNameInput, StepRemoteUrlInput, StepCommitInput:
		return true
	default:
		return false
	}
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
		UpdateStyles(flavor, config.GetAccentColor(flavor, m.ConfigModel.Accents[m.Selected]))
	}

	if m.CurrentStep == StepOptionsFlavorSelect {
		flavor := config.GetCatppuccinFlavor(m.ConfigModel.Flavors[m.Selected])
		UpdateStyles(flavor, config.GetAccentColor(flavor, m.CurrentConfig.Accent))
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

func (m *Model) outputByLevel(out string) {
	// Fill output array up to current level
	for len(m.Output) <= m.Level {
		m.Output = append(m.Output, "")
	}
	m.Output[m.Level] += out
}

func (m *Model) handleBranchOperation() (*Model, tea.Cmd) {
	switch m.BranchModel.SelectedAction {
	case "Create Branch":
		return m.prepareBranchAddition()

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

		m.outputByLevel(strings.Join(branches, "\n"))
		m.Success = "Listed branches"
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) prepareBranchAddition() (*Model, tea.Cmd) {
	m.Selected = 0
	m.CurrentStep = StepBranchCreate
	m.Level = 3
	return m, nil
}

func (m *Model) executeBranchAction() (*Model, tea.Cmd) {
	switch m.BranchModel.SelectedAction {
	case "Switch Branch":
		return m.switchBranch()
	case "Delete Branch":
		return m.deleteBranch()
	}
	return m, tea.Quit
}

func (m Model) handleBranchCreateSelection() (tea.Model, tea.Cmd) {
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
func (m *Model) handleBranchInputSubmit() (*Model, tea.Cmd) {
	if strings.TrimSpace(m.BranchModel.Input) == "" {
		m.Err = "Branch name cannot be empty"
		return m, tea.Quit
	}

	return m.createBranch(strings.TrimSpace(m.BranchModel.Input))
}

func (m *Model) createBranch(branchName string) (*Model, tea.Cmd) {
	out, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()
	if err != nil {
		m.outputByLevel(string(out))
		m.Err = fmt.Sprintf("Failed to create branch: %s", string(out))
	} else {
		m.Success = fmt.Sprintf("Successfully created branch '%s'", branchName)
	}
	return m, tea.Quit
}

func (m *Model) switchBranch() (*Model, tea.Cmd) {
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

func (m *Model) deleteBranch() (*Model, tea.Cmd) {
	// Try to delete local branch first
	cmd := exec.Command("git", "branch", "-d", m.BranchModel.SelectedBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.outputByLevel(string(output) + "\n")
		// suggest force delete:
		m.outputByLevel("\\cpTo force delete use: \ngit branch -D" + m.BranchModel.SelectedBranch)
		m.Err = "Delete failed"
	} else {
		m.outputByLevel(string(output) + "\n")
		m.Success = fmt.Sprintf("Deleted branch '%s'", m.BranchModel.SelectedBranch)
	}
	return m, tea.Quit
}

// function to handle commit message submission
func (m *Model) handleCommitMessageSubmit() (*Model, tea.Cmd) {
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
		m.outputByLevel("\\ceError:\n" + out)
		m.Err = "Failed to Commit"
	} else {
		m.Success = "Commited Changes"
	}

	return m, tea.Quit
}

func (m *Model) handleTagOperation() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "List Tags":
		out, err := git.GetNLatestTags(10)
		m.outputByLevel("\\ct10 latest tags:\n" + out)
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Listed Tags"
		}
		return m, tea.Quit

	case "Add Tag":
		return m.prepareTagAddition()

	case "Remove Tag", "Push Tag":
		return m.prepareTagSelection()
	}
	return m, nil
}

func (m Model) handleTagAddSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedAddTag = m.TagModel.AddOptions[m.Selected]

	if m.TagModel.SelectedAddTag == "Manual Input" {
		// Switch to input mode
		m.TagModel.ManualInput = ""
		m.CurrentStep = StepTagInput
		return m, nil
	} else {
		// Extract the version from the option (e.g., "Patch (v1.2.4)" -> "v1.2.4")
		start := strings.Index(m.TagModel.SelectedAddTag, "(")
		end := strings.Index(m.TagModel.SelectedAddTag, ")")
		if start != -1 && end != -1 && end > start {
			newTag := m.TagModel.SelectedAddTag[start+1 : end]
			return m.createTag(newTag)
		} else {
			m.Err = "Failed to parse version from selection"
			return m, tea.Quit
		}
	}
}

// function to handle manual tag input submission
func (m Model) handleTagInputSubmit() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.TagModel.ManualInput) == "" {
		m.Err = "Tag name cannot be empty"
		return m, tea.Quit
	}

	return m.createTag(strings.TrimSpace(m.TagModel.ManualInput))
}

func (m *Model) createTag(tagName string) (*Model, tea.Cmd) {
	out, err := exec.Command("git", "tag", tagName).CombinedOutput()
	if err != nil {
		m.outputByLevel(string(out))
		m.Err = fmt.Sprintf("Failed to create tag: %s", string(out))
	} else {
		m.Success = fmt.Sprintf("Successfully created tag '%s'", tagName)
	}
	return m, tea.Quit
}

func (m *Model) prepareTagSelection() (*Model, tea.Cmd) {
	var out string
	var err error

	if m.TagModel.SelectedAction == "Remove Tag" {
		out, err = git.GetNLatestTags(10)
	} else {
		out, err = git.GetNLatestTags(1)
	}

	if err != nil {
		m.outputByLevel("\nError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get tags"
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.outputByLevel("You can add tags with Tags -> Add Tag or `git tag <tag_name>`")
		m.Err = "No tags found"
		return m, tea.Quit
	}

	m.Selected = 0
	m.TagModel.Options = strings.Split(strings.TrimSpace(out), "\n")
	m.TagModel.Options = append(m.TagModel.Options, "Load all tags")
	m.TagModel.SelectedOption = ""
	m.CurrentStep = StepTagSelect
	m.Level = 3
	return m, nil
}

func (m *Model) prepareTagAddition() (*Model, tea.Cmd) {
	// Get the latest tag
	out, err := git.GetNLatestTags(1)
	if err != nil {
		m.outputByLevel("\\crError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get latest tag"
		return m, tea.Quit
	}

	latestTag := strings.TrimSpace(out)
	m.TagModel.CurrentTag = latestTag

	// Prepare add options
	if latestTag == "" {
		// No existing tags, only allow manual input
		m.TagModel.AddOptions = []string{"Manual Input"}
		m.outputByLevel("\\ctNo existing tags found. You can create the first tag manually.")
	} else {
		// Check if latest tag is semantic version
		if _, _, _, _, err := internal.ParseSemanticVersion(latestTag); err == nil {
			// Valid semantic version, offer all options
			patchVersion, _ := internal.GenerateNextSemanticVersion(latestTag, "Patch")
			minorVersion, _ := internal.GenerateNextSemanticVersion(latestTag, "Minor")
			majorVersion, _ := internal.GenerateNextSemanticVersion(latestTag, "Major")

			m.TagModel.AddOptions = []string{
				fmt.Sprintf("Patch (%s)", patchVersion),
				fmt.Sprintf("Minor (%s)", minorVersion),
				fmt.Sprintf("Major (%s)", majorVersion),
				"Manual Input",
			}
			m.outputByLevel(fmt.Sprintf("\\ctCurrent latest tag: %s", latestTag))
		} else {
			// Not semantic version, only allow manual input
			m.TagModel.AddOptions = []string{"Manual Input"}
			m.outputByLevel("\\ctNote:\n You will get automatic completions for adding tags when using semantic versioning")
		}
	}

	m.Selected = 0
	m.CurrentStep = StepTagAdd
	m.Level = 3
	return m, nil
}

func (m *Model) loadAllTags() (*Model, tea.Cmd) {
	out, err := git.GetAllTags()

	if err != nil {
		m.outputByLevel("\nError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get tags"
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.outputByLevel("You can add tags with Tags -> Add Tag or `git tag <tag_name>`")
		m.Err = "No tags found"
		return m, tea.Quit
	}

	m.Selected = 0
	m.TagModel.Options = append([]string{"Only load 10 latest tags"}, strings.Split(strings.TrimSpace(out), "\n")...)
	m.TagModel.SelectedOption = ""
	m.CurrentStep = StepTagSelect
	m.Level = 3
	return m, nil
}

func (m *Model) executeTagAction() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "Remove Tag":
		switch m.TagModel.SelectedOption {
		case "Load all tags":
			m.loadAllTags()
			return m, nil
		case "Only load 10 latest tags":
			m.prepareTagSelection()
			return m, nil
		default:
			out, err := exec.Command("git", "tag", "-d", m.TagModel.SelectedOption).CombinedOutput()
			m.outputByLevel(string(out))
			if err != nil {
				m.Err = fmt.Sprintf("Failed to remove: %s", string(out))
			} else {
				m.Success = fmt.Sprintf("Removed Tag '%s'", m.TagModel.SelectedOption)
			}
		}
	case "Push Tag":
		switch m.TagModel.SelectedOption {
		case "Load all tags":
			//m.loadAllTags()
			m.Err = "Not implemented yet"
		case "Only load 10 latest tags": //TODO: add to options from start
			//m.prepareTagSelection()
			m.Err = "Not implemented yet"
		default:
			out, err := exec.Command("git", "push", "origin", m.TagModel.SelectedOption).CombinedOutput()
			if err != nil {
				//FIXME: next line is a hotfix for some very weird bug with the output when usng raw out
				outResult := []string{}
				for line := range strings.SplitSeq(string(out), "\n") {
					if strings.TrimSpace(line) != "" {
						outResult = append(outResult, strings.TrimSpace(line))
					}
				}
				m.outputByLevel("\\crError:\n" + strings.Join(outResult, "\n"))
				m.Err = "Failed to push"
			} else {
				m.Success = fmt.Sprintf("Pushed Tag '%s'", m.TagModel.SelectedOption)
			}
		}
	}
	return m, tea.Quit
}

func (m *Model) handleRemoteOperation() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "List Remotes":
		out, err := git.GetRemotes()
		if err != "" {
			m.outputByLevel("\\crError:\n" + err)
			m.Err = out
		} else {
			m.outputByLevel(out)
			m.Success = "Listed all Remotes"
		}
		return m, tea.Quit

	case "Add Remote":
		m.CurrentStep = StepRemoteNameInput
		m.Level = 3
		return m, nil

	case "Remove Remote":
		return m.prepareRemoteSelection()
	}
	return m, nil
}

func (m *Model) prepareRemoteSelection() (*Model, tea.Cmd) {
	out, err := git.GetRemotesClean()

	if err != "" {
		m.outputByLevel("\\crError:\n" + err)
		m.Err = out
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.outputByLevel("You can add remotes with Remote -> Add Remote or `git remote add <remote_name> <remote_url>`")
		m.Err = "No remotes found"
		return m, tea.Quit
	}

	m.Selected = 0
	m.RemoteModel.Options = strings.Split(strings.TrimSpace(out), "\n")
	m.RemoteModel.SelectedOption = ""
	m.CurrentStep = StepRemoteSelect
	m.Level = 3
	return m, nil
}

func (m *Model) executeRemoteAction() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "Remove Remote":
		out, err := git.RemoveRemote(m.RemoteModel.SelectedOption)
		if err != "" {
			m.outputByLevel("\\crError:\n%s" + err)
			m.Err = out
		} else {
			m.outputByLevel(out)
			m.Success = fmt.Sprintf("Removed Remote '%s'", m.RemoteModel.SelectedOption)
		}
	}
	return m, tea.Quit
}

func (m Model) handleOptionsActionSelection() (tea.Model, tea.Cmd) {
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

	case "Reset to Defaults":
		m.CurrentConfig = &config.Config{
			Accent: config.DefaultConfig.Accent,
			Flavor: config.DefaultConfig.Flavor,
		}
		if err := config.SaveConfig(m.CurrentConfig); err != nil {
			m.Err = fmt.Sprintf("Failed to save config: %v", err)
		} else {
			UpdateStylesByConfig(m.CurrentConfig)
			m.Success = "Configuration reset to defaults"
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleOptionsFlavorSelection() (tea.Model, tea.Cmd) {
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
		UpdateStylesByConfig(m.CurrentConfig)
		m.Success = fmt.Sprintf("Flavor changed to %s", selectedFlavor)
	}
	return m, tea.Quit
}

func (m Model) handleOptionsAccentSelection() (tea.Model, tea.Cmd) {
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
		UpdateStylesByConfig(m.CurrentConfig)
		m.Success = fmt.Sprintf("Accent changed to %s", selectedAccent)
	}
	return m, tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case repoUpdatedMsg, repoUpdateErrorMsg:
		m.Loading = false

		if _, isErr := msg.(repoUpdateErrorMsg); isErr {
			m.outputByLevel("Failed to fetch from remote")
		}

		// quick-start flows triggered via StartAt
		switch m.StartAt {
		case "add-tag":
			m.prepareTagAddition()
			m.ActionModel.SelectedAction = "Tag"
			m.TagModel.SelectedAction = "Add Tag"

		case "push-tag":
			m.prepareTagSelection()
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
					return m.handleTagInputSubmit()
				case StepBranchInput:
					return m.handleBranchInputSubmit()
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
						m.outputByLevel(string(out))
						m.Err = "Failed to add remote"
					} else {
						m.Success = fmt.Sprintf("Added Remote '%s' -> %s", m.RemoteModel.NameInput, m.RemoteModel.UrlInput)
					}
					return m, tea.Quit
				case StepCommitInput:
					return m.handleCommitMessageSubmit()
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
		return m.handleActionSelection()

	case StepBranchAction:
		return m.handleBranchActionSelection()
	case StepBranchSelect:
		return m.handleBranchSelection()
	case StepBranchCreate:
		return m.handleBranchCreateSelection()

	case StepCommitAction:
		return m.handleCommitSelection()
	case StepCommitSelectPrefix:
		return m.handleCommitPrefixSelection()

	case StepTag:
		return m.handleTagActionSelection()
	case StepTagSelect:
		return m.handleTagSelection()
	case StepTagAdd:
		return m.handleTagAddSelection()

	case StepRemote:
		return m.handleRemoteActionSelection()
	case StepRemoteSelect:
		return m.handleRemoteSelection()

	case StepOptions:
		return m.handleOptionsActionSelection()
	case StepOptionsFlavorSelect:
		return m.handleOptionsFlavorSelection()
	case StepOptionsAccentSelect:
		return m.handleOptionsAccentSelection()
	}
	return m, nil
}

func (m Model) handleActionSelection() (tea.Model, tea.Cmd) {
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
					m.outputByLevel("\\cyModified:\n " + strings.Join(status["modified"], "\n ") + "\n╌\n")
				}
				if len(status["added"]) > 0 {
					m.outputByLevel("\\cgAdded:\n " + strings.Join(status["added"], "\n ") + "\n╌\n")
				}
				if len(status["deleted"]) > 0 {
					m.outputByLevel("\\crDeleted:\n " + strings.Join(status["deleted"], "\n ") + "\n╌\n")
				}
				if len(status["untracked"]) > 0 {
					m.outputByLevel("\\ctUntracked:\n " + strings.Join(status["untracked"], "\n ") + "\n╌\n")
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

func (m Model) handleBranchActionSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedAction = m.BranchModel.Actions[m.Selected]
	return m.handleBranchOperation()
}

func (m Model) handleBranchSelection() (tea.Model, tea.Cmd) {
	m.BranchModel.SelectedBranch = m.BranchModel.Branches[m.Selected]
	return m.executeBranchAction()
}

func (m Model) handleCommitSelection() (tea.Model, tea.Cmd) {
	m.CommitModel.SelectedAction = m.CommitModel.Actions[m.Selected]

	if m.CommitModel.SelectedAction == "Undo Last Commit" {
		out, err := git.UndoLastCommit()
		m.outputByLevel(out)
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

func (m Model) handleCommitPrefixSelection() (tea.Model, tea.Cmd) {
	m.CommitModel.SelectedPrefix = m.CommitModel.CommitPrefixes[m.Selected]

	if m.CommitModel.SelectedPrefix != "Custom Prefix" {
		m.CommitModel.CommitMessage = m.CommitModel.SelectedPrefix + ": "
	}

	m.Selected = 0
	m.CurrentStep = StepCommitInput
	return m, nil
}

func (m Model) handleTagActionSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedAction = m.TagModel.Actions[m.Selected]
	return m.handleTagOperation()
}

func (m Model) handleTagSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedOption = m.TagModel.Options[m.Selected]
	return m.executeTagAction()
}

func (m Model) handleRemoteActionSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.SelectedAction = m.RemoteModel.Actions[m.Selected]
	return m.handleRemoteOperation()
}

func (m Model) handleRemoteSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.SelectedOption = m.RemoteModel.Options[m.Selected]
	return m.executeRemoteAction()
}
