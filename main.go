package main

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateMainMenu state = iota
	stateBranchList
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	logoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Bold(true)
)

// Item for bubble list

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

// Model

type model struct {
	state       state
	mainMenu    list.Model
	branchList  list.Model
	branches    []string
	err         error
}

func initialModel() model {
	items := []list.Item{
		item("Switch Branch"),
		item("Commit"),
	}

	mainMenu := list.New(items, list.NewDefaultDelegate(), 0, 0)
	mainMenu.Title = "Choose action"

	branchList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	branchList.Title = "Select branch"

	return model{
		state:      stateMainMenu,
		mainMenu:   mainMenu,
		branchList: branchList,
	}
}

// Fetch branches
func fetchBranches() ([]string, error) {
	exec.Command("git", "fetch", "-a").Run() // don't care about error for now
	out, err := exec.Command("git", "branch").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var branches []string
	for _, line := range lines {
		b := strings.TrimSpace(strings.TrimPrefix(line, "origin/"))
		if b != "" && !strings.Contains(b, "*") {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.mainMenu.SetSize(msg.Width-4, msg.Height-6)
		m.branchList.SetSize(msg.Width-4, msg.Height-6)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case stateMainMenu:
				branches, err := fetchBranches()
				if err != nil {
					m.err = err
					return m, nil
				}
				m.branches = branches
				items := make([]list.Item, len(branches))
				for i, b := range branches {
					items[i] = item(b)
				}
				m.branchList.SetItems(items)
				m.state = stateBranchList
			case stateBranchList:
				selected := m.branchList.SelectedItem().(item)
				cmd := exec.Command("git", "switch", string(selected))
				output, err := cmd.CombinedOutput()
				if err != nil {
					m.err = fmt.Errorf("switch failed: %s", output)
				} else {
					m.err = fmt.Errorf("✔ Switched to %s", selected)
				}
				return m, tea.Quit
			}
		}
	}

	if m.state == stateMainMenu {
		var cmd tea.Cmd
		m.mainMenu, cmd = m.mainMenu.Update(msg)
		return m, cmd
	} else if m.state == stateBranchList {
		var cmd tea.Cmd
		m.branchList, cmd = m.branchList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	header := logoStyle.Render(`
  ▘▗ ▌ 
▛▌▌▜▘▛▌
▙▌▌▐▖▌▌
▄▌     
`)

	if m.err != nil {
		header += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.err.Error()) + "\n"
	}

	switch m.state {
	case stateMainMenu:
		return header + "\n" + m.mainMenu.View()
	case stateBranchList:
		return header + "\n" + m.branchList.View()
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
