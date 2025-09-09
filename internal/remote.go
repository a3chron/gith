package internal

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleRemoteActionSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.SelectedAction = m.RemoteModel.Actions[m.Selected]
	return m.HandleRemoteOperation()
}

func (m Model) HandleRemoteSelection() (tea.Model, tea.Cmd) {
	m.RemoteModel.SelectedOption = m.RemoteModel.Options[m.Selected]
	return m.ExecuteRemoteAction()
}

func (m *Model) HandleRemoteOperation() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "List Remotes":
		out, err := git.GetRemotes()
		if err != "" {
			m.OutputByLevel("\\crError:\n" + err)
			m.Err = out
		} else {
			m.OutputByLevel(git.FormatFancyRemotes(out))
			m.Success = "Listed all Remotes"
		}
		return m, tea.Quit

	case "Add Remote":
		m.CurrentStep = StepRemoteNameInput
		m.Level = 3
		return m, nil

	case "Remove Remote":
		return m.PrepareRemoteSelection()
	}
	return m, nil
}

func (m *Model) PrepareRemoteSelection() (*Model, tea.Cmd) {
	out, err := git.GetRemotesClean()

	if err != "" {
		m.OutputByLevel("\\crError:\n" + err)
		m.Err = out
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.OutputByLevel("You can add remotes with Remote -> Add Remote or `git remote add <remote_name> <remote_url>`")
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

func (m *Model) ExecuteRemoteAction() (*Model, tea.Cmd) {
	switch m.RemoteModel.SelectedAction {
	case "Remove Remote":
		out, err := git.RemoveRemote(m.RemoteModel.SelectedOption)
		if err != "" {
			m.OutputByLevel("\\crError:\n%s" + err)
			m.Err = out
		} else {
			m.OutputByLevel(out)
			m.Success = fmt.Sprintf("Removed Remote '%s'", m.RemoteModel.SelectedOption)
		}
	}
	return m, tea.Quit
}
