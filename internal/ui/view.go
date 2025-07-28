package ui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	if m.CurrentStep == StepLoad {
		return m.renderLoadingView()
	}
	return m.renderMainView()
}

func (m Model) renderLoadingView() string {
	return `
		  ▘▗ ▌ 
		▛▌▌▜▘▛▌
		▙▌▌▐▖▌▌
		▄▌ - by a3chron
	` + "\n\n" + m.Spinner.View() + " Running git fetch"
}

func (m Model) renderMainView() string {
	var content strings.Builder
	line := LineStyle.Render("│")

	content.WriteString(m.renderHeader(line))

	content.WriteString(m.renderActionSelection(line))

	if m.ActionModel.SelectedAction != "" {
		content.WriteString(m.renderSubActions(line))
	}

	if m.Output != "" {
		content.WriteString(m.renderOutput(line))
	}

	content.WriteString(m.renderResult(line))

	content.WriteString(m.renderNavigationHints())

	return ContainerStyle.Render(content.String())
}

func (m Model) renderHeader(line string) string {
	return LineStyle.Render("╭─╌") + " " + AccentStyle.Render("gith") + "\n" + line + "\n"
}

func (m Model) renderActionSelection(line string) string {
	var content strings.Builder
	bullet := m.getBullet()

	content.WriteString(bullet + " " + TextStyle.Render("Select action") + "\n")

	if m.ActionModel.SelectedAction == "" {
		// Show action options
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.ActionModel.Actions, line, m.CurrentStep == StepAction))
		}
	} else {
		// Show completed action
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ActionModel.SelectedAction) + "\n")
	}

	return content.String()
}

func (m Model) renderSubActions(line string) string {
	switch m.ActionModel.SelectedAction {
	case "Branch":
		return m.renderBranchActions(line)
	case "Commit":
		return m.renderCommitActions(line)
	case "Tag":
		return m.renderTagActions(line)
	case "Remote":
		return m.renderRemoteActions(line)
	}
	return ""
}

func (m Model) renderBranchActions(line string) string {
	var content strings.Builder
	bullet := m.getBullet()

	content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select branch action") + "\n")

	if m.BranchModel.SelectedAction == "" && len(m.BranchModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.BranchModel.Actions, line, m.CurrentStep == StepBranchAction))
		}
	} else if m.BranchModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedAction) + "\n")

		if (m.BranchModel.SelectedAction == "Switch Branch" || m.BranchModel.SelectedAction == "Delete Branch") && len(m.BranchModel.Branches) > 0 {
			action := strings.ToLower(strings.Split(m.BranchModel.SelectedAction, " ")[0])
			content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Branch to "+action) + "\n")

			if m.BranchModel.SelectedBranch == "" {
				if m.Err == "" {
					content.WriteString(m.renderOptions(m.BranchModel.Branches, line, m.CurrentStep == StepBranchSelect))
				}
			} else {
				content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedBranch) + "\n")
			}
		}
	}

	return content.String()
}

func (m Model) renderCommitActions(line string) string {
	var content strings.Builder
	bullet := m.getBullet()

	content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select commit action") + "\n")

	if m.CommitModel.SelectedAction == "" && len(m.CommitModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.CommitModel.Actions, line, m.CurrentStep == StepCommit))
		}
	} else if m.CommitModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.CommitModel.SelectedAction) + "\n")
	}

	return content.String()
}

func (m Model) renderTagActions(line string) string {
	var content strings.Builder
	bullet := m.getBullet()

	content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select tag action") + "\n")

	if m.TagModel.SelectedAction == "" && len(m.TagModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.TagModel.Actions, line, m.CurrentStep == StepTag))
		}
	} else if m.TagModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.TagModel.SelectedAction) + "\n")

		if (m.TagModel.SelectedAction == "Remove Tag" || m.TagModel.SelectedAction == "Push Tag") && len(m.TagModel.Options) > 0 {
			actionText := strings.ToLower(m.TagModel.SelectedAction)
			content.WriteString(line + "\n" + bullet + " " + TextStyle.Render(actionText) + "\n")

			if m.TagModel.Selected == "" {
				if m.Err == "" {
					content.WriteString(m.renderOptions(m.TagModel.Options, line, m.CurrentStep == StepTagSelect))
				}
			} else {
				content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.TagModel.Selected) + "\n")
			}
		}
	}

	return content.String()
}

func (m Model) renderRemoteActions(line string) string {
	var content strings.Builder
	bullet := m.getBullet()

	content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select remote action") + "\n")

	if m.RemoteModel.SelectedAction == "" && len(m.RemoteModel.Actions) > 0 {
		if m.Err == "" {
			content.WriteString(m.renderOptions(m.RemoteModel.Actions, line, m.CurrentStep == StepRemote))
		}
	} else if m.RemoteModel.SelectedAction != "" {
		content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.RemoteModel.SelectedAction) + "\n")
	}

	return content.String()
}

func (m Model) renderOptions(options []string, line string, isCurrentStep bool) string {
	var content strings.Builder

	for i, option := range options {
		var bullet string
		if i == m.Selected && isCurrentStep {
			bullet = BulletStyle.Render("●")
			content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, SelectedStyle.Render(option)))
		} else {
			bullet = DimStyle.Render("○")
			content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, NormalStyle.Render(option)))
		}
	}

	return content.String()
}

func (m Model) renderOutput(line string) string {
	var content strings.Builder
	outputLines := strings.Split(m.Output, "\n")

	content.WriteString(line + "\n")
	for _, outputLine := range outputLines {
		trimmed := strings.TrimSpace(outputLine)
		if trimmed != "" {
			content.WriteString(line + " " + DimStyle.Render(outputLine) + "\n")
		}
	}

	return content.String()
}

func (m Model) renderResult(line string) string {
	if m.Err != "" {
		return line + "\n" + LineStyle.Render("╰─╌") + " " + ErrorStyle.Render(m.Err) + "\n"
	} else if m.Success != "" {
		return line + "\n" + LineStyle.Render("╰─╌") + " " + SuccessStyle.Render(m.Success) + "\n"
	}
	return ""
}

func (m Model) renderNavigationHints() string {
	return "\n\n" + DimStyle.Render("Use ↑↓ to navigate, enter to select, ctrl+h to go back, q / esc to quit")
}

func (m Model) getBullet() string {
	if m.Err != "" {
		return ErrorStyle.Render("■")
	}
	return BulletStyle.Render("◆")
}
