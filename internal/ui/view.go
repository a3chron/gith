package ui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var content strings.Builder
	var line = LineStyle.Render("│")

	if m.CurrentStep == StepLoad {
		return `
			  ▘▗ ▌ 
			▛▌▌▜▘▛▌
			▙▌▌▐▖▌▌
			▄▌ - by a3chron
		` + "\n\n" + m.Spinner.View() + "Running git fetch"
	} else {
		// Step 1: Select action
		if m.ActionModel.SelectedAction == "" {
			var bullet string
			bullet = BulletStyle.Render("◆")
			if m.Err != "" {
				bullet = ErrorStyle.Render("■")
			}

			content.WriteString(LineStyle.Render("╭─╌") + " " + AccentStyle.Render("gith") + "\n" + line + "\n" + bullet + " " + TextStyle.Render("Select action") + "\n")

			if m.Err == "" {
				for i, option := range m.ActionModel.Actions {
					if i == m.Selected && m.CurrentStep == StepAction {
						bullet = BulletStyle.Render("●")
						content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, SelectedStyle.Render(option)))
					} else {
						bullet = DimStyle.Render("○")
						content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, NormalStyle.Render(option)))
					}
				}
			}
		} else {
			// Show completed action selection
			content.WriteString(LineStyle.Render("╭─╌") + " " + AccentStyle.Render("gith") + "\n" + line + "\n" + BulletStyle.Render("◆") + " " + TextStyle.Render("Select action") + "\n")
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.ActionModel.SelectedAction) + "\n")
		}

		// Next Steps
		if m.ActionModel.SelectedAction != "" {
			// Step 2
			switch m.ActionModel.SelectedAction {
			case "Branch":
				var bullet = BulletStyle.Render("◆")
				if m.Err != "" {
					bullet = ErrorStyle.Render("■")
				}

				if m.BranchModel.SelectedAction == "" && len(m.BranchModel.Actions) > 0 {
					content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select branch action") + "\n")
					if m.Err == "" {
						for i, action := range m.BranchModel.Actions {
							var bullet string

							if i == m.Selected && m.CurrentStep == StepBranchAction {
								bullet = BulletStyle.Render("●")
								content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, SelectedStyle.Render(action)))
							} else {
								bullet = DimStyle.Render("○")
								content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, NormalStyle.Render(action)))
							}
						}
					}
				} else if m.BranchModel.SelectedAction != "" {
					// Show completed branch action selection
					content.WriteString(bullet + " " + TextStyle.Render("Select branch action") + "\n")
					content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedAction) + "\n")
				}

				// Step 3
				if m.BranchModel.SelectedAction == "Switch branch" || m.BranchModel.SelectedAction == "Delete branch" {
					if m.BranchModel.SelectedBranch == "" && len(m.BranchModel.Branches) > 0 {
						content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Branch to "+strings.ToLower(strings.Split(m.BranchModel.SelectedAction, " ")[0])) + "\n")
						if m.Err == "" {
							for i, branch := range m.BranchModel.Branches {
								var bullet string

								if i == m.Selected && m.CurrentStep == StepBranchSelect {
									bullet = BulletStyle.Render("●")
									content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, SelectedStyle.Render(branch)))
								} else {
									bullet = DimStyle.Render("○")
									content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, NormalStyle.Render(branch)))
								}
							}
						}
					} else if m.BranchModel.SelectedBranch != "" {
						// Show completed branch selection
						content.WriteString(line + "\n" + bullet + " " + TextStyle.Render("Select branch") + "\n")
						content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.BranchModel.SelectedBranch) + "\n")
					}
				}
			case "Commit":
				//TODO: thsi and other actions
			}
		}

		if m.Output != "" {
			var outputArr = strings.Split(m.Output, "\n")

			for _, out := range outputArr {
				var outTrim = strings.TrimSpace(out)
				if outTrim != "" && outTrim != "\n" {
					content.WriteString(line + " " + DimStyle.Render(out) + "\n")
				}
			}
		}

		// Show result
		if m.Err != "" {
			content.WriteString(line + "\n" + LineStyle.Render("╰─╌") + " " + ErrorStyle.Render(m.Err) + "\n")
		} else if m.Success != "" {
			content.WriteString(line + "\n" + LineStyle.Render("╰─╌") + " " + SuccessStyle.Render(m.Success) + "\n")
		}
		// Show navigation hints
		content.WriteString("\n\n" + DimStyle.Render("Use ↑↓ to navigate, enter to select, ctrl+h to go back, q / esc to quit"))
	}

	return ContainerStyle.Render(content.String())
}
