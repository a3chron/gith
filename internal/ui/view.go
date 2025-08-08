package ui

import (
	"fmt"
	"strings"

	"github.com/a3chron/gith/internal/config"
	"github.com/charmbracelet/lipgloss"
)

// View renders the entire UI based on the current model state.
func (m Model) View() string {
	if m.CurrentStep == StepLoad {
		return m.renderLoadingView()
	}
	return m.renderMainView()
}

// renderLoadingView displays a spinner and a loading message.
func (m Model) renderLoadingView() string {
	bigLogo := `
		        ██┐         ██┐
		        └─┘   ██┐   ██│
		██████┐ ██┐ ██████┐ ██████┐
		██┌─██│ ██│ └─██┌─┘ ██┌─██│
		██│ ██│ ██│   ██│   ██│ ██│
		██████│ ██│   ████┐ ██│ ██│
		└───██│ └─┘   └───┘ └─┘ └─┘
		    ██│
		██████│
		└─────┘
	`
	return bigLogo + "\n\n" + m.Spinner.View() + " Running git fetch"
}

// renderMainView constructs the main application view from its components.
func (m Model) renderMainView() string {
	var content strings.Builder
	line := LineStyle.Render("│")

	content.WriteString(m.renderHeader())
	content.WriteString(m.renderOutput(line, 0)) // For initial/global output

	content.WriteString(m.renderActionSelection())
	content.WriteString(m.renderOutput(line, 1)) // Output for level 1

	if m.ActionModel.SelectedAction != "" {
		content.WriteString(m.renderSubActions())
		content.WriteString(m.renderOutput(line, 2)) // Output for level 2

		content.WriteString(m.renderSubActions2())
		content.WriteString(m.renderOutput(line, 3)) // Output for level 3

		if m.CurrentStep == StepTagInput {
			content.WriteString(m.renderOutput(line, 4)) // Output for level 4
		}
	}

	content.WriteString(m.renderResult())

	content.WriteString(m.renderNavigationHints())

	return ContainerStyle.Render(content.String())
}

// renderHeader returns the application title bar.
func (m Model) renderHeader() string {
	return LineStyle.Render("╭─╌") + " " + AccentStyle.Render("gith") + "\n"
}

// renderActionSelection renders the primary action choices (Branch, Commit, etc.).
func (m Model) renderActionSelection() string {
	var content strings.Builder
	bullet := m.getBullet(1)

	content.WriteString(bullet + " " + TextStyle.Render("Select action") + "\n")

	if m.ActionModel.SelectedAction == "" {
		// Show action options
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.ActionModel.Actions, m.CurrentStep == StepAction))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else {
		// Show completed action
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ActionModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderSubActions delegates rendering of the second-level action lists.
func (m Model) renderSubActions() string {
	switch m.ActionModel.SelectedAction {
	case "Branch":
		return m.renderBranchActions()
	case "Commit":
		return m.renderCommitActions()
	case "Tag":
		return m.renderTagActions()
	case "Remote":
		return m.renderRemoteActions()
	case "Options":
		return m.renderOptionsActions()
	}
	return ""
}

// renderSubActions2 delegates rendering of the third-level action lists and confirmations.
func (m Model) renderSubActions2() string {
	switch m.ActionModel.SelectedAction {
	case "Branch":
		return m.renderBranchSubActions2()
	case "Tag":
		return m.renderTagSubActions2()
	case "Remote":
		return m.renderRemoteSubActions2()
	case "Options":
		return m.renderOptionsSubActions2()
	}
	return ""
}

// renderBranchActions now only renders the list of available branch actions.
func (m Model) renderBranchActions() string {
	var content strings.Builder
	bullet := m.getBullet(2)

	content.WriteString(bullet + " " + TextStyle.Render("Select branch action") + "\n")

	if m.BranchModel.SelectedAction == "" && len(m.BranchModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.BranchModel.Actions, m.CurrentStep == StepBranchAction))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else if m.BranchModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderBranchSubActions2 handles selection of a specific branch.
func (m Model) renderBranchSubActions2() string {
	var content strings.Builder
	bullet := m.getBullet(3)

	switch m.BranchModel.SelectedAction {
	case "Switch Branch", "Delete Branch":
		content.WriteString(bullet + " " + TextStyle.Render(m.BranchModel.SelectedAction) + "\n")

		// If no branch is selected yet, show the list of branches to choose from.
		if m.BranchModel.SelectedBranch == "" {
			if len(m.BranchModel.Branches) > 0 && m.Err == "" {
				content.WriteString(m.renderOptions(m.BranchModel.Branches, m.CurrentStep == StepBranchSelect))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else { // A branch has been selected, show it as completed.
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedBranch) + "\n")
		}
	case "Create Branch":
		content.WriteString(bullet + " " + TextStyle.Render("Create Branch") + "\n")

		if m.BranchModel.SelectedOption == "" {
			// Show add options (feat, fix, refactor, ..., manual)
			if len(m.BranchModel.Options) > 0 && m.Err == "" {
				content.WriteString(m.renderOptions(m.BranchModel.Options, m.CurrentStep == StepBranchCreate))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else if m.CurrentStep == StepBranchInput {
			// Show input field
			if m.Err == "" {
				content.WriteString(m.renderBranchInput())
			}
		} else {
			// Show completed selection
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedOption) + "\n")
		}
	}

	return content.String()
}

// renderCommitActions renders the list of available commit actions.
func (m Model) renderCommitActions() string {
	var content strings.Builder
	bullet := m.getBullet(2)

	content.WriteString(bullet + " " + TextStyle.Render("Select commit action") + "\n")

	if m.CommitModel.SelectedAction == "" && len(m.CommitModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.CommitModel.Actions, m.CurrentStep == StepCommit))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else if m.CommitModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.CommitModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderTagActions now only renders the list of available tag actions.
func (m Model) renderTagActions() string {
	var content strings.Builder
	bullet := m.getBullet(2)

	content.WriteString(bullet + " " + TextStyle.Render("Select tag action") + "\n")

	if m.TagModel.SelectedAction == "" && len(m.TagModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.TagModel.Actions, m.CurrentStep == StepTag))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else if m.TagModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.TagModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderTagSubActions2 handles selection of a specific tag.
func (m Model) renderTagSubActions2() string {
	var content strings.Builder
	bullet := m.getBullet(3)

	switch m.TagModel.SelectedAction {
	case "Remove Tag", "Push Tag":
		content.WriteString(bullet + " " + TextStyle.Render(m.TagModel.SelectedAction) + "\n")

		// If no tag is selected yet, show the list of tags to choose from.
		if m.TagModel.SelectedOption == "" {
			if len(m.TagModel.Options) > 0 && m.Err == "" {
				content.WriteString(m.renderOptions(m.TagModel.Options, m.CurrentStep == StepTagSelect))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else { // A tag has been selected, show it as completed.
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.TagModel.SelectedOption) + "\n")
		}

	case "Add Tag":
		content.WriteString(bullet + " " + TextStyle.Render("Add Tag") + "\n")

		if m.TagModel.SelectedAddTag == "" {
			// Show add options (patch, minor, major, manual)
			if len(m.TagModel.AddOptions) > 0 && m.Err == "" {
				content.WriteString(AccentStyle.Render("├╌") + " " + DimStyle.Render("Latest tag:") + " " + NormalStyle.Render(m.TagModel.CurrentTag) + "\n")
				content.WriteString(m.renderOptions(m.TagModel.AddOptions, m.CurrentStep == StepTagAdd))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else if m.CurrentStep == StepTagInput {
			// Show input field
			if m.Err == "" {
				content.WriteString(m.renderTagInput())
			}
		} else {
			// Show completed selection
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.TagModel.SelectedAddTag) + "\n")
		}
	}

	return content.String()
}

func (m Model) renderRemoteSubActions2() string {
	var content strings.Builder
	bullet := m.getBullet(3)

	switch m.RemoteModel.SelectedAction {
	case "Remove Remote":
		content.WriteString(bullet + " " + TextStyle.Render(m.RemoteModel.SelectedAction) + "\n")

		// If no remote is selected yet, show the list of remotes to choose from.
		if m.RemoteModel.SelectedOption == "" {
			if len(m.RemoteModel.Options) > 0 && m.Err == "" {
				content.WriteString(m.renderOptions(m.RemoteModel.Options, m.CurrentStep == StepRemoteSelect))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else { // A remote has been selected, show it as completed.
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.SelectedOption) + "\n")
		}
	case "Add Remote":
		content.WriteString(bullet + " " + TextStyle.Render("Add Tag") + "\n")

		if m.CurrentStep == StepRemoteNameInput {
			// Show input field
			if m.Err == "" {
				content.WriteString(AccentStyle.Render("├╌") + " " + AccentStyle.Render("Remote Name:") + " " + DimStyle.Render("(e.g. origin)") + "\n")
				content.WriteString(m.renderRemoteInput())
			}
		} else {
			// Completed Remote Name Input
			// Show second input for url input
			if m.CurrentStep == StepRemoteUrlInput && m.Err == "" && m.Success == "" {
				// Show input of still active add remote input
				content.WriteString(AccentStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.NameInput) + "\n")

				// Show input field
				if m.Err == "" {
					content.WriteString(AccentStyle.Render("├╌") + " " + AccentStyle.Render("Remote Url:") + " " + DimStyle.Render("(e.g. git@github.com:a3chron/gith.git)") + "\n")
					content.WriteString(m.renderRemoteInput())
				}
			} else {
				// Show completed selection
				content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.NameInput) + "\n")
				// Show completed selection
				content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.UrlInput) + "\n")
			}
		}
	}

	return content.String()
}

// renderRemoteActions renders the list of available remote actions.
func (m Model) renderRemoteActions() string {
	var content strings.Builder
	bullet := m.getBullet(2)

	content.WriteString(bullet + " " + TextStyle.Render("Select remote action") + "\n")

	if m.RemoteModel.SelectedAction == "" && len(m.RemoteModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.RemoteModel.Actions, m.CurrentStep == StepRemote))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else if m.RemoteModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderOptionsActions renders the list of available options actions.
func (m Model) renderOptionsActions() string {
	var content strings.Builder
	bullet := m.getBullet(2)

	content.WriteString(bullet + " " + TextStyle.Render("Select options action") + "\n")

	if m.ConfigModel.SelectedAction == "" && len(m.ConfigModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.ConfigModel.Actions, m.CurrentStep == StepOptions))
			content.WriteString(AccentStyle.Render("╰─╌") + "\n")
		}
	} else if m.ConfigModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ConfigModel.SelectedAction) + "\n")
	}

	return content.String()
}

// renderOptionsSubActions2 handles selection of specific configuration options.
func (m Model) renderOptionsSubActions2() string {
	var content strings.Builder
	bullet := m.getBullet(3)

	// Safety check for CurrentConfig
	if m.CurrentConfig == nil {
		return ""
	}

	switch m.ConfigModel.SelectedAction {
	case "Change Flavor":
		content.WriteString(bullet + " " + TextStyle.Render("Select flavor") + "\n")

		if m.ConfigModel.SelectedFlavor == "" {
			if len(m.ConfigModel.Flavors) > 0 && m.Err == "" {
				content.WriteString(AccentStyle.Render("├╌") + " " + DimStyle.Render("Current: ") + AccentStyle.Render(m.CurrentConfig.Flavor) + "\n")
				content.WriteString(m.renderOptionsWithPreview(m.ConfigModel.Flavors, m.CurrentStep == StepOptionsFlavorSelect, "flavor"))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else {
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ConfigModel.SelectedFlavor) + "\n")
		}

	case "Change Accent":
		content.WriteString(bullet + " " + TextStyle.Render("Select accent") + "\n")

		if m.ConfigModel.SelectedAccent == "" {
			if len(m.ConfigModel.Accents) > 0 && m.Err == "" {
				content.WriteString(AccentStyle.Render("├╌") + " " + DimStyle.Render("Current: ") + AccentStyle.Render(m.CurrentConfig.Accent) + "\n")
				content.WriteString(m.renderOptionsWithPreview(m.ConfigModel.Accents, m.CurrentStep == StepOptionsAccentSelect, "accent"))
				content.WriteString(AccentStyle.Render("╰─╌") + "\n")
			}
		} else {
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ConfigModel.SelectedAccent) + "\n")
		}
	}

	return content.String()
}

// Update renderOptionsWithPreview with safety checks:
func (m Model) renderOptionsWithPreview(options []string, isCurrentStep bool, optionType string) string {
	var content strings.Builder
	accLine := AccentStyle.Render("│")

	// Safety check for CurrentConfig
	if m.CurrentConfig == nil {
		return content.String()
	}

	for i, option := range options {
		var bullet string
		var preview string

		// Create a preview of what the option would look like
		switch optionType {
		case "flavor":
			flavor := config.GetCatppuccinFlavor(option)
			accentColor := config.GetAccentColor(flavor, m.CurrentConfig.Accent)
			preview = lipgloss.NewStyle().Foreground(lipgloss.Color(accentColor.Hex)).Render("●")
		case "accent":
			flavor := config.GetCatppuccinFlavor(m.CurrentConfig.Flavor)
			accentColor := config.GetAccentColor(flavor, option)
			preview = lipgloss.NewStyle().Foreground(lipgloss.Color(accentColor.Hex)).Render("●")
		}

		if i == m.Selected && isCurrentStep {
			bullet = BulletStyle.Render("●")
			content.WriteString(fmt.Sprintf("%s %s %s %s\n", accLine, bullet, SelectedStyle.Render(option), preview))
		} else {
			bullet = DimStyle.Render("○")
			content.WriteString(fmt.Sprintf("%s %s %s %s\n", accLine, bullet, NormalStyle.Render(option), preview))
		}
	}

	return content.String()
}

// renderOptions formats a list of choices for selection.
func (m Model) renderOptions(options []string, isCurrentStep bool) string {
	var content strings.Builder
	accLine := AccentStyle.Render("│")

	for i, option := range options {
		var bullet string
		if i == m.Selected && isCurrentStep {
			bullet = BulletStyle.Render("●")
			content.WriteString(fmt.Sprintf("%s %s %s\n", accLine, bullet, SelectedStyle.Render(option)))
		} else {
			bullet = DimStyle.Render("○")
			content.WriteString(fmt.Sprintf("%s %s %s\n", accLine, bullet, NormalStyle.Render(option)))
		}
	}

	return content.String()
}

func (m Model) renderTagInput() string {
	var content strings.Builder
	cursor := "_"
	line := LineStyle.Render("│")
	inputText := m.TagModel.ManualInput

	if m.Success == "" {
		line = AccentStyle.Render("│")
		content.WriteString(line + " " + NormalStyle.Render("Enter tag name:") + "\n")
	}

	// Add cursor at the end
	displayText := inputText + cursor

	if m.Success == "" {
		content.WriteString(line + " " + AccentStyle.Render("> ") + displayText + "\n")
		content.WriteString(AccentStyle.Render("╰─╌") + "\n")
	} else {
		content.WriteString(line + " " + CompletedStyle.Render("> "+inputText) + "\n")
	}

	return content.String()
}

// TODO: maybe combine this and above into one, a lot common
func (m Model) renderBranchInput() string {
	var content strings.Builder
	cursor := "_"
	line := LineStyle.Render("│")
	inputText := m.BranchModel.Input

	if m.Success == "" {
		line = AccentStyle.Render("│")
		content.WriteString(line + " " + NormalStyle.Render("Enter branch name:") + "\n")
	}

	// Add cursor at the end
	displayText := inputText + cursor

	if m.Success == "" {
		content.WriteString(line + " " + AccentStyle.Render("> ") + displayText + "\n")
		content.WriteString(AccentStyle.Render("╰─╌") + "\n")
	} else {
		content.WriteString(line + " " + CompletedStyle.Render("> "+inputText) + "\n")
	}

	return content.String()
}

func (m Model) renderRemoteInput() string {
	var content strings.Builder
	var inputText string

	cursor := "_"
	line := LineStyle.Render("│")
	accentLine := AccentStyle.Render("│")

	if m.CurrentStep == StepRemoteNameInput {
		inputText = m.RemoteModel.NameInput
	} else {
		inputText = m.RemoteModel.UrlInput
	}

	// Add cursor at the end
	displayText := inputText + cursor

	if m.Success == "" {
		content.WriteString(accentLine + " " + AccentStyle.Render("> ") + displayText + "\n")
		content.WriteString(AccentStyle.Render("╰─╌") + "\n")
	} else {
		content.WriteString(line + " " + CompletedStyle.Render("> "+inputText) + "\n")
	}

	return content.String()
}

// renderOutput displays command output at a specific interaction level.
// It handles multi-line and color-coded messages.
func (m Model) renderOutput(line string, level int) string {
	// Don't render output for a level that hasn't been reached or has no message.
	if m.Level > level {
		return line + "\n"
	}

	if len(m.Output) <= level {
		return ""
	}

	var content strings.Builder
	outputLines := strings.Split(m.Output[level], "\n")

	var renderedLines []string
	for _, outputLine := range outputLines {
		trimmed := strings.TrimSpace(outputLine)
		if trimmed != "" {
			// Check for color prefix (e.g., "\cgSuccess!")
			if coloredText, found := strings.CutPrefix(trimmed, "\\c"); found && len(coloredText) > 1 {
				color := coloredText[:1]
				text := coloredText[1:]

				switch color {
				case "g":
					renderedLines = append(renderedLines, line+" "+GreenStyle.Render(text))
				case "y":
					renderedLines = append(renderedLines, line+" "+YellowStyle.Render(text))
				case "r":
					renderedLines = append(renderedLines, line+" "+RedStyle.Render(text))
				case "t":
					renderedLines = append(renderedLines, line+" "+TextStyle.Render(text))
				case "a":
					renderedLines = append(renderedLines, line+" "+AccentStyle.Render(text))
				default:
					// Fallback for unknown color codes
					// TODO: add all colors & maybe simplify
					renderedLines = append(renderedLines, line+" "+DimStyle.Render(outputLine))
				}
			} else {
				renderedLines = append(renderedLines, line+" "+DimStyle.Render(outputLine))
			}
		}
	}

	if len(renderedLines) == 0 {
		return line + "\n"
	}

	content.WriteString(line + "\n")
	for _, l := range renderedLines {
		content.WriteString(l + "\n")
	}
	content.WriteString(line + "\n")

	return content.String()
}

// renderResult displays the final success or error message.
func (m Model) renderResult() string {
	if m.Err != "" {
		return LineStyle.Render("│\n╰─╌") + " " + ErrorStyle.Render(m.Err) + "\n"
	} else if m.Success != "" {
		return LineStyle.Render("│\n╰─╌") + " " + SuccessStyle.Render(m.Success) + "\n"
	}
	return ""
}

// renderNavigationHints displays help text for the user.
func (m Model) renderNavigationHints() string {
	switch m.CurrentStep {
	case StepTagInput, StepBranchInput, StepRemoteNameInput, StepRemoteUrlInput:
		return "\n\n" + DimStyle.Render("Type name, enter to confirm, ctrl+h to go back, esc to quit")
	case StepOptionsAccentSelect:
		return "\n\n" + DimStyle.Render("Select Accent to preview, enter to confirm, ctrl+h to go back, esc to quit")
	default:
		return "\n\n" + DimStyle.Render("Use ↑↓ to navigate, enter to select, ctrl+h to go back, q / esc to quit")
	}
}

// getBullet returns the appropriate bullet character based on the current state.
func (m Model) getBullet(forLevel int) string {
	if forLevel == m.Level {
		if m.Err != "" {
			// bullet for failed (=last) selection
			return ErrorStyle.Render("□")
		}
		// bullet for current selection
		return BulletStyle.Render("◆")
	}
	// bullet for already selected
	return BulletStyle.Render("◇")
}
