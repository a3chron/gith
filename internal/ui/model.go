package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/a3chron/gith/internal"
	"github.com/a3chron/gith/internal/git"
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
	StepRemoteSelect
	StepChanges
	StepOptions
)

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
	Options        []string
	Selected       string
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
	Level       int
	Output      []string
	Err         string
	Success     string
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
	case StepCommit:
		return m.CommitModel.Actions
	case StepTag:
		return m.TagModel.Actions
	case StepTagSelect:
		return m.TagModel.Options
	case StepRemote:
		return m.RemoteModel.Actions
	case StepRemoteSelect:
		return m.RemoteModel.Options
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
	m.Level = 1
	m.Err = ""
	m.Output = []string{} // TODO add m.Output[0] if exisiting
	m.Success = ""
}

func (m *Model) outputByLevel(out string) {
	for len(m.Output) <= m.Level {
		m.Output = append(m.Output, "")
	}
	m.Output[m.Level] += out
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
		m.outputByLevel(string(output) + "\n")
		// If that fails, try force delete
		cmd = exec.Command("git", "branch", "-D", m.BranchModel.SelectedBranch)
		output, err = cmd.CombinedOutput()
		if err != nil {
			m.Err = fmt.Sprintf("Delete failed: %s", string(output))
		} else {
			m.Success = fmt.Sprintf("Force deleted branch '%s'", m.BranchModel.SelectedBranch)
		}
	} else {
		m.outputByLevel(string(output) + "\n")
		m.Success = fmt.Sprintf("Deleted branch '%s'", m.BranchModel.SelectedBranch)
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
	m.TagModel.Selected = ""
	m.CurrentStep = StepTagSelect
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
	m.TagModel.Selected = ""
	m.CurrentStep = StepTagSelect
	m.Level = 3
	return m, nil
}

func (m *Model) executeTagAction() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "Remove Tag":
		switch m.TagModel.Selected {
		case "Load all tags":
			m.loadAllTags()
			return m, nil
		case "Only load 10 latest tags":
			m.prepareTagSelection()
			return m, nil
		default:
			out, err := exec.Command("git", "tag", "-d", m.TagModel.Selected).CombinedOutput()
			m.outputByLevel(string(out))
			if err != nil {
				m.Err = fmt.Sprintf("Failed to remove: %s", string(out))
			} else {
				m.Success = fmt.Sprintf("Removed Tag '%s'", m.TagModel.Selected)
			}
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
		m.Err = "Add Remote functionality not implemented yet"
		return m, tea.Quit

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
	m.RemoteModel.Selected = ""
	m.CurrentStep = StepRemoteSelect
	m.Level = 3
	return m, nil
}

func (m *Model) executeRemoteAction() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "Remove Remote":
		out, err := git.RemoveRemote(m.RemoteModel.Selected)
		if err != "" {
			m.outputByLevel("\\crError:\n%s" + err)
			m.Err = out
		} else {
			m.outputByLevel(out)
			m.Success = fmt.Sprintf("Removed Remote '%s'", m.RemoteModel.Selected)
		}
	}
	return m, tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case repoUpdatedMsg:
		m.Loading = false
		m.CurrentStep = StepAction
		m.Level = 1
		return m, nil

	case repoUpdateErrorMsg:
		m.Loading = false
		m.CurrentStep = StepAction
		m.outputByLevel("Failed to fetch from remote")
		m.Level = 1
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
	case StepRemoteSelect:
		return m.handleRemoteSelection()
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
				m.Success = fmt.Sprintf("Changes detected: %d files", length)
			}
		}
		return m, tea.Quit
	case "Commit":
		m.Selected = 0
		m.CurrentStep = StepCommit
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
		m.CurrentStep = StepChanges
		m.Level = 2
		m.Err = "Changes will come here soon"
		return m, tea.Quit
	case "Options":
		m.CurrentStep = StepOptions
		m.Level = 2
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
		m.outputByLevel(out)
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

func (m Model) handleRemoteSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.Selected = m.RemoteModel.Options[m.Selected]
	return m.executeRemoteAction()
}
