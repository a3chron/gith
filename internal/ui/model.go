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

const (
	StepLoad Step = iota
	StepAction
	StepBranchAction
	StepBranchSelect
	StepCommit
	StepTag
	StepTagSelect
	StepRemote
	StepChanges
	StepOptions
)

type SelectionModel interface {
	GetActions() []string
	GetSelectedAction() string
	SetSelectedAction(string)
}

type ActionModel struct {
	Actions        []string
	SelectedAction string
}

func (a *ActionModel) GetActions() []string       { return a.Actions }
func (a *ActionModel) GetSelectedAction() string  { return a.SelectedAction }
func (a *ActionModel) SetSelectedAction(s string) { a.SelectedAction = s }

type BranchModel struct {
	Actions        []string
	SelectedAction string
	Branches       []string
	SelectedBranch string
}

func (b *BranchModel) GetActions() []string       { return b.Actions }
func (b *BranchModel) GetSelectedAction() string  { return b.SelectedAction }
func (b *BranchModel) SetSelectedAction(s string) { b.SelectedAction = s }

type CommitModel struct {
	Actions        []string
	SelectedAction string
}

func (c *CommitModel) GetActions() []string       { return c.Actions }
func (c *CommitModel) GetSelectedAction() string  { return c.SelectedAction }
func (c *CommitModel) SetSelectedAction(s string) { c.SelectedAction = s }

type TagModel struct {
	Actions        []string
	SelectedAction string
	Options        []string
	Selected       string
}

func (t *TagModel) GetActions() []string       { return t.Actions }
func (t *TagModel) GetSelectedAction() string  { return t.SelectedAction }
func (t *TagModel) SetSelectedAction(s string) { t.SelectedAction = s }

type RemoteModel struct {
	Actions        []string
	SelectedAction string
}

func (r *RemoteModel) GetActions() []string       { return r.Actions }
func (r *RemoteModel) GetSelectedAction() string  { return r.SelectedAction }
func (r *RemoteModel) SetSelectedAction(s string) { r.SelectedAction = s }

type ConfigModel struct {
	Accent  string
	Flavor  string
	Flavors []string
	Accents []string
}

type Model struct {
	CurrentStep Step
	Loading     bool
	Selected    int
	ActionModel ActionModel
	BranchModel BranchModel
	CommitModel CommitModel
	RemoteModel RemoteModel
	TagModel    TagModel
	ConfigModel ConfigModel
	Spinner     spinner.Model
	Output      string
	Err         string
	Success     string
}

type repoUpdatedMsg struct{}

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

func (m Model) getCurrentOptions() []string {
	switch m.CurrentStep {
	case StepAction:
		return m.ActionModel.Actions
	case StepBranchAction:
		return m.BranchModel.Actions
	case StepBranchSelect:
		return m.BranchModel.Branches
	case StepCommit:
		return m.CommitModel.Actions
	case StepTag:
		return m.TagModel.Actions
	case StepTagSelect:
		return m.TagModel.Options
	case StepRemote:
		return m.RemoteModel.Actions
	default:
		return []string{}
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
}

func (m *Model) resetState() {
	m.CurrentStep = StepAction
	m.ActionModel.SelectedAction = ""
	m.BranchModel.SelectedBranch = ""
	m.BranchModel.SelectedAction = ""
	m.CommitModel.SelectedAction = ""
	m.TagModel.SelectedAction = ""
	m.TagModel.Selected = ""
	m.RemoteModel.SelectedAction = ""
	m.Selected = 0
	m.Err = ""
	m.Output = ""
	m.Success = ""
}

func (m *Model) executeGitStatus() (*Model, tea.Cmd) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		m.Err = fmt.Sprintf("Failed to get status: %v", err)
	} else {
		if len(output) == 0 {
			m.Success = "Working tree clean"
		} else {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			m.Output = string(output)
			m.Success = fmt.Sprintf("Changes detected: %d files", len(lines))
		}
	}
	return m, tea.Quit
}

func (m *Model) handleBranchOperation() (*Model, tea.Cmd) {
	switch m.BranchModel.SelectedAction {
	case "Create Branch":
		m.Err = "Create branch functionality not implemented yet"
		return m, tea.Quit

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
	}
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
		m.Output += string(output) + "\n"
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
	return m, tea.Quit
}

func (m *Model) handleTagOperation() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "List Tags":
		out, err := git.GetNLatestTags(10)
		m.Output = "10 latest tags:\n---\n" + out
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Listed Tags"
		}
		return m, tea.Quit

	case "Add Tag":
		m.Err = "Tag functionality not implemented yet"
		return m, tea.Quit

	case "Remove Tag", "Push Tag":
		return m.prepareTagSelection()
	}
	return m, nil
}

func (m *Model) prepareTagSelection() (*Model, tea.Cmd) {
	var out string
	var err error

	if m.TagModel.SelectedAction == "Remove Tag" {
		out, err = git.GetAllTags()
	} else {
		out, err = git.GetNLatestTags(1)
	}

	if err != nil {
		m.Output += "\nError:\n" + fmt.Sprintf("%v", err)
		m.Err = "Failed to get tags"
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.Output += "You can add tags with Tags -> Add Tag or `git tag <tag_name>`"
		m.Err = "No tags found"
		return m, tea.Quit
	}

	m.Selected = 0
	m.TagModel.Options = strings.Split(strings.TrimSpace(out), "\n")
	m.CurrentStep = StepTagSelect
	return m, nil
}

func (m *Model) executeTagAction() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "Remove Tag":
		out, err := exec.Command("git", "tag", "-d", m.TagModel.Selected).CombinedOutput()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to remove: %s", string(out))
		} else {
			m.Success = fmt.Sprintf("Removed Tag '%s'", m.TagModel.Selected)
		}
	case "Push Tag":
		out, err := exec.Command("git", "push", "origin", m.TagModel.Selected).CombinedOutput()
		if err != nil {
			m.Err = fmt.Sprintf("Failed to push: %s", string(out))
		} else {
			m.Success = fmt.Sprintf("Pushed Tag '%s'", m.TagModel.Selected)
		}
	}
	return m, tea.Quit
}

func (m *Model) handleRemoteOperation() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "List Remotes":
		out, err := git.ListRemotes()
		m.Output = out
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Listed all Remotes"
		}
		return m, tea.Quit

	case "Add Remote":
		m.Err = "Add Remote functionality not implemented yet"
		return m, tea.Quit
	}
	return m, nil
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
	case StepCommit:
		return m.handleCommitSelection()
	case StepTag:
		return m.handleTagActionSelection()
	case StepTagSelect:
		return m.handleTagSelection()
	case StepRemote:
		return m.handleRemoteActionSelection()
	}
	return m, nil
}

func (m Model) handleActionSelection() (tea.Model, tea.Cmd) {
	m.Output = ""
	m.ActionModel.SelectedAction = m.ActionModel.Actions[m.Selected]

	switch m.ActionModel.SelectedAction {
	case "Branch":
		m.Selected = 0
		m.CurrentStep = StepBranchAction
	case "Status":
		return m.executeGitStatus()
	case "Commit":
		m.Selected = 0
		m.CurrentStep = StepCommit
	case "Tag":
		m.Selected = 0
		m.CurrentStep = StepTag
	case "Remote":
		m.Selected = 0
		m.CurrentStep = StepRemote
	case "Changes":
		m.CurrentStep = StepChanges
		m.Err = "Changes will come here soon"
		return m, tea.Quit
	case "Options":
		m.CurrentStep = StepOptions
		m.Err = "Options will come here soon"
		return m, tea.Quit
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

	switch m.CommitModel.SelectedAction {
	case "Undo Last Commit":
		out, err := git.UndoLastCommit()
		m.Output = out
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Undo Commit Successful"
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleTagActionSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedAction = m.TagModel.Actions[m.Selected]
	return m.handleTagOperation()
}

func (m Model) handleTagSelection() (tea.Model, tea.Cmd) {
	m.TagModel.Selected = m.TagModel.Options[m.Selected]
	return m.executeTagAction()
}

func (m Model) handleRemoteActionSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.SelectedAction = m.RemoteModel.Actions[m.Selected]
	return m.handleRemoteOperation()
}
