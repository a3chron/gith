package main

import (
	"fmt"
	"os/exec"
	"strings"

	catppuccin "github.com/catppuccin/go"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepAction step = iota
	stepBranch
	stepComplete
)

var (
	mocha = catppuccin.Mocha

	// Styles using Catppuccin Mocha palette
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Blue().Hex)).
			Bold(true)

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Text().Hex)).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Base().Hex)).
			Background(lipgloss.Color(mocha.Blue().Hex)).
			Bold(true).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Text().Hex))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Overlay0().Hex))

	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Green().Hex))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Red().Hex))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Green().Hex))

	bulletStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Blue().Hex)).
			Bold(true)

	lineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mocha.Surface1().Hex))

	containerStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

type model struct {
	currentStep    step
	selectedAction string
	selectedBranch string
	options        []string
	branches       []string
	selected       int
	err            string
	success        string
}

func initialModel() model {
	return model{
		currentStep: stepAction,
		options:     []string{"Switch Branch", "Create Branch", "Delete Branch", "Status"},
		selected:    0,
	}
}

func fetchBranches() ([]string, error) {
	// Fetch remote branches first
	exec.Command("git", "fetch", "-a").Run()

	out, err := exec.Command("git", "branch", "-a").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var branches []string
	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip current branch (has * prefix)
		if strings.HasPrefix(line, "* ") {
			continue
		}

		// Handle remote branches
		if strings.HasPrefix(line, "remotes/origin/") {
			branch := strings.TrimPrefix(line, "remotes/origin/")
			// Skip HEAD pointer
			if strings.Contains(branch, "HEAD ->") {
				continue
			}
			if branch != "" && !seen[branch] {
				branches = append(branches, branch)
				seen[branch] = true
			}
		} else {
			// Local branches
			if line != "" && !seen[line] {
				branches = append(branches, line)
				seen[line] = true
			}
		}
	}

	return branches, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Go back to previous step
			switch m.currentStep {
			case stepBranch:
				m.currentStep = stepAction
				m.selected = 0
				m.err = ""
				m.branches = nil
			case stepComplete:
				m.currentStep = stepAction
				m.selectedAction = ""
				m.selectedBranch = ""
				m.selected = 0
				m.err = ""
				m.success = ""
				m.branches = nil
			default:
				return m, tea.Quit
			}
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			switch m.currentStep {
			case stepAction:
				if m.selected < len(m.options)-1 {
					m.selected++
				}
			case stepBranch:
				if m.selected < len(m.branches)-1 {
					m.selected++
				}
			}
		case "enter":
			switch m.currentStep {
			case stepAction:
				m.selectedAction = m.options[m.selected]
				switch m.selectedAction {
				case "Switch Branch", "Delete Branch":
					branches, err := fetchBranches()
					if err != nil {
						m.err = fmt.Sprintf("Failed to fetch branches: %v", err)
						m.currentStep = stepComplete
						return m, nil
					}
					if len(branches) == 0 {
						m.err = "No branches available"
						m.currentStep = stepComplete
						return m, nil
					}
					m.branches = branches
					m.selected = 0
					m.currentStep = stepBranch
					m.err = ""
				case "Create Branch":
					m.success = "Create branch functionality not implemented yet"
					m.currentStep = stepComplete
				case "Status":
					cmd := exec.Command("git", "status", "--porcelain")
					output, err := cmd.Output()
					if err != nil {
						m.err = fmt.Sprintf("Failed to get status: %v", err)
					} else {
						if len(output) == 0 {
							m.success = "Working tree clean"
						} else {
							lines := strings.Split(strings.TrimSpace(string(output)), "\n")
							m.success = fmt.Sprintf("Changes detected: %d files", len(lines))
						}
					}
					m.currentStep = stepComplete
				}
			case stepBranch:
				if len(m.branches) > 0 && m.selected < len(m.branches) {
					m.selectedBranch = m.branches[m.selected]

					switch m.selectedAction {
					case "Switch Branch":
						// First check if branch exists locally
						cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+m.selectedBranch)
						err := cmd.Run()

						var switchCmd *exec.Cmd
						if err != nil {
							// Branch doesn't exist locally, create and track remote
							switchCmd = exec.Command("git", "checkout", "-b", m.selectedBranch, "origin/"+m.selectedBranch)
						} else {
							// Branch exists locally, just switch
							switchCmd = exec.Command("git", "checkout", m.selectedBranch)
						}

						output, err := switchCmd.CombinedOutput()
						if err != nil {
							m.err = fmt.Sprintf("Switch failed: %s", string(output))
						} else {
							m.success = fmt.Sprintf("✓ Switched to branch '%s'", m.selectedBranch)
						}
					case "Delete Branch":
						// Try to delete local branch first
						cmd := exec.Command("git", "branch", "-d", m.selectedBranch)
						output, err := cmd.CombinedOutput()
						if err != nil {
							// If that fails, try force delete
							cmd = exec.Command("git", "branch", "-D", m.selectedBranch)
							output, err = cmd.CombinedOutput()
							if err != nil {
								m.err = fmt.Sprintf("Delete failed: %s", string(output))
							} else {
								m.success = fmt.Sprintf("✓ Force deleted branch '%s'", m.selectedBranch)
							}
						} else {
							m.success = fmt.Sprintf("✓ Deleted branch '%s'", m.selectedBranch)
						}
					}
					m.currentStep = stepComplete
				}
			case stepComplete:
				// Reset to beginning
				m.currentStep = stepAction
				m.selectedAction = ""
				m.selectedBranch = ""
				m.selected = 0
				m.err = ""
				m.success = ""
				m.branches = nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render(`
		  ▘▗ ▌ 
		▛▌▌▜▘▛▌
		▙▌▌▐▖▌▌
		▄▌ - by a3chron
	`) + "\n\n")

	// Step 1: Select action (always visible)
	if m.selectedAction == "" {
		content.WriteString(lineStyle.Render("┌─") + " " + bulletStyle.Render("◆") + " " + dimStyle.Render("Select action") + "\n")
		for i, option := range m.options {
			var bullet, line string
			if i == len(m.options)-1 {
				line = lineStyle.Render("└─")
			} else {
				line = lineStyle.Render("├─")
			}

			if i == m.selected && m.currentStep == stepAction {
				bullet = bulletStyle.Render("●")
				content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, selectedStyle.Render(option)))
			} else {
				bullet = dimStyle.Render("○")
				content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, normalStyle.Render(option)))
			}
		}
	} else {
		// Show completed action selection
		content.WriteString(lineStyle.Render("┌─") + " " + bulletStyle.Render("◆") + " " + dimStyle.Render("Select action") + "\n")
		content.WriteString(lineStyle.Render("├─") + " " + completedStyle.Render("●") + " " + completedStyle.Render(m.selectedAction) + "\n")
	}

	// Step 2: Select branch (visible after action is selected)
	if m.selectedAction != "" && (m.selectedAction == "Switch Branch" || m.selectedAction == "Delete Branch") {
		if m.selectedBranch == "" && len(m.branches) > 0 {
			content.WriteString(lineStyle.Render("├─") + " " + bulletStyle.Render("◆") + " " + dimStyle.Render("Select branch") + "\n")
			for i, branch := range m.branches {
				var bullet, line string
				if i == len(m.branches)-1 {
					line = lineStyle.Render("└─")
				} else {
					line = lineStyle.Render("├─")
				}

				if i == m.selected && m.currentStep == stepBranch {
					bullet = bulletStyle.Render("●")
					content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, selectedStyle.Render(branch)))
				} else {
					bullet = dimStyle.Render("○")
					content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, normalStyle.Render(branch)))
				}
			}
		} else if m.selectedBranch != "" {
			// Show completed branch selection
			content.WriteString(lineStyle.Render("├─") + " " + bulletStyle.Render("◆") + " " + dimStyle.Render("Select branch") + "\n")
			content.WriteString(lineStyle.Render("└─") + " " + completedStyle.Render("●") + " " + completedStyle.Render(m.selectedBranch) + "\n")
		}
	}

	// Show result
	if m.currentStep == stepComplete {
		content.WriteString("\n")
		if m.err != "" {
			content.WriteString(errorStyle.Render("✗ "+m.err) + "\n")
		} else if m.success != "" {
			content.WriteString(successStyle.Render(m.success) + "\n")
		}
		content.WriteString("\n" + dimStyle.Render("Press enter to continue, esc to restart"))
	} else {
		// Show navigation hints
		content.WriteString("\n" + dimStyle.Render("Use ↑↓ to navigate, enter to select, esc to go back, q to quit"))
	}

	return containerStyle.Render(content.String())
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
