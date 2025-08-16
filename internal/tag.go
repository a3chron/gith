package internal

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/a3chron/gith/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleTagActionSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedAction = m.TagModel.Actions[m.Selected]
	return m.HandleTagOperation()
}

func (m Model) HandleTagSelection() (tea.Model, tea.Cmd) {
	m.TagModel.SelectedOption = m.TagModel.Options[m.Selected]
	return m.ExecuteTagAction()
}

func (m Model) HandleTagAddSelection() (tea.Model, tea.Cmd) {
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
			return m.CreateTag(newTag)
		} else {
			m.Err = "Failed to parse version from selection"
			return m, tea.Quit
		}
	}
}

func (m *Model) HandleTagOperation() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "List Tags":
		out, err := git.GetNLatestTags(10)
		m.OutputByLevel("\\ct10 latest tags:\n" + out)
		if err != nil {
			m.Err = fmt.Sprintf("%v", err)
		} else {
			m.Success = "Listed Tags"
		}
		return m, tea.Quit

	case "Add Tag":
		return m.PrepareTagAddition()

	case "Remove Tag", "Push Tag":
		return m.PrepareTagSelection()
	}
	return m, nil
}

// function to handle manual tag input submission
func (m Model) HandleTagInputSubmit() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.TagModel.ManualInput) == "" {
		m.Err = "Tag name cannot be empty"
		return m, tea.Quit
	}

	return m.CreateTag(strings.TrimSpace(m.TagModel.ManualInput))
}

func (m *Model) CreateTag(tagName string) (*Model, tea.Cmd) {
	out, err := exec.Command("git", "tag", tagName).CombinedOutput()
	if err != nil {
		m.OutputByLevel(string(out))
		m.Err = fmt.Sprintf("Failed to create tag: %s", string(out))
	} else {
		m.Success = fmt.Sprintf("Successfully created tag '%s'", tagName)
	}
	return m, tea.Quit
}

func (m *Model) PrepareTagSelection() (*Model, tea.Cmd) {
	var out string
	var err error

	if m.TagModel.SelectedAction == "Remove Tag" {
		out, err = git.GetNLatestTags(10)
	} else {
		out, err = git.GetNLatestTags(1)
	}

	if err != nil {
		m.OutputByLevel("\nError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get tags"
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.OutputByLevel("You can add tags with Tags -> Add Tag or `git tag <tag_name>`")
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

func (m *Model) PrepareTagAddition() (*Model, tea.Cmd) {
	// Get the latest tag
	out, err := git.GetNLatestTags(1)
	if err != nil {
		m.OutputByLevel("\\crError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get latest tag"
		return m, tea.Quit
	}

	latestTag := strings.TrimSpace(out)
	m.TagModel.CurrentTag = latestTag

	// Prepare add options
	if latestTag == "" {
		// No existing tags, only allow manual input
		m.TagModel.AddOptions = []string{"Manual Input"}
		m.OutputByLevel("\\ctNo existing tags found. You can create the first tag manually.")
	} else {
		// Check if latest tag is semantic version
		if _, _, _, _, err := ParseSemanticVersion(latestTag); err == nil {
			// Valid semantic version, offer all options
			patchVersion, _ := GenerateNextSemanticVersion(latestTag, "Patch")
			minorVersion, _ := GenerateNextSemanticVersion(latestTag, "Minor")
			majorVersion, _ := GenerateNextSemanticVersion(latestTag, "Major")

			m.TagModel.AddOptions = []string{
				fmt.Sprintf("Patch (%s)", patchVersion),
				fmt.Sprintf("Minor (%s)", minorVersion),
				fmt.Sprintf("Major (%s)", majorVersion),
				"Manual Input",
			}
			m.OutputByLevel(fmt.Sprintf("\\ctCurrent latest tag: %s", latestTag))
		} else {
			// Not semantic version, only allow manual input
			m.TagModel.AddOptions = []string{"Manual Input"}
			m.OutputByLevel("\\ctNote:\n You will get automatic completions for adding tags when using semantic versioning")
		}
	}

	m.Selected = 0
	m.CurrentStep = StepTagAdd
	m.Level = 3
	return m, nil
}

func (m *Model) LoadAllTags() (*Model, tea.Cmd) {
	out, err := git.GetAllTags()

	if err != nil {
		m.OutputByLevel("\nError:\n" + fmt.Sprintf("%v", err))
		m.Err = "Failed to get tags"
		return m, tea.Quit
	}

	if strings.TrimSpace(out) == "" {
		m.OutputByLevel("You can add tags with Tags -> Add Tag or `git tag <tag_name>`")
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

func (m *Model) ExecuteTagAction() (*Model, tea.Cmd) {
	switch m.TagModel.SelectedAction {
	case "Remove Tag":
		switch m.TagModel.SelectedOption {
		case "Load all tags":
			m.LoadAllTags()
			return m, nil
		case "Only load 10 latest tags":
			m.PrepareTagSelection()
			return m, nil
		default:
			out, err := exec.Command("git", "tag", "-d", m.TagModel.SelectedOption).CombinedOutput()
			m.OutputByLevel(string(out))
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
				m.OutputByLevel("\\crError:\n" + strings.Join(outResult, "\n"))
				m.Err = "Failed to push"
			} else {
				m.Success = fmt.Sprintf("Pushed Tag '%s'", m.TagModel.SelectedOption)
			}
		}
	}
	return m, tea.Quit
}
