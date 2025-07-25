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
		// Step 1: Select action (always visible)
		if m.SelectedAction == "" {
			var bullet string
			bullet = BulletStyle.Render("◆")
			if m.Err != "" {
				bullet = ErrorStyle.Render("■")
			}

			content.WriteString(LineStyle.Render("╭─╌") + " " + AccentStyle.Render("gith") + "\n" + line + "\n" + bullet + " " + TextStyle.Render("Select action") + "\n")

			if m.Err == "" {
				for i, option := range m.Options {
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
			content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.SelectedAction) + "\n" + line + "\n")
		}

		// Step 2: Select branch (visible after action is selected)
		if m.SelectedAction != "" && (m.SelectedAction == "Switch Branch" || m.SelectedAction == "Delete Branch") {
			var bullet = BulletStyle.Render("◆")
			if m.Err != "" {
				bullet = ErrorStyle.Render("■")
			}

			if m.SelectedBranch == "" && len(m.Branches) > 0 {
				content.WriteString(bullet + " " + TextStyle.Render("Select branch") + "\n")
				if m.Err == "" {
					for i, branch := range m.Branches {
						var bullet string

						if i == m.Selected && m.CurrentStep == StepBranch {
							bullet = BulletStyle.Render("●")
							content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, SelectedStyle.Render(branch)))
						} else {
							bullet = DimStyle.Render("○")
							content.WriteString(fmt.Sprintf("%s %s %s\n", line, bullet, NormalStyle.Render(branch)))
						}
					}
				}
			} else if m.SelectedBranch != "" {
				// Show completed branch selection
				content.WriteString(bullet + " " + TextStyle.Render("Select branch") + "\n")
				content.WriteString(LineStyle.Render("├╌") + " " + CompletedStyle.Render(m.SelectedBranch) + "\n")
			}
		}

		if m.CurrentStep == StepOptions {
			content.WriteString(BulletStyle.Render("◆") + " Options will come here soon" + "\n")
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
		if m.CurrentStep == StepComplete {
			if m.Err != "" {
				content.WriteString(ErrorStyle.Render(m.Err) + "\n")
			} else if m.Success != "" {
				content.WriteString(SuccessStyle.Render(m.Success) + "\n")
			}
			content.WriteString("\n\n" + DimStyle.Render("Press enter to continue, esc to restart"))
		} else {
			if m.Err != "" {
				content.WriteString(line + "\n" + LineStyle.Render("╰─╌") + " " + ErrorStyle.Render(m.Err) + "\n")
			} else if m.Success != "" {
				content.WriteString(line + "\n" + LineStyle.Render("╰─╌") + " " + SuccessStyle.Render(m.Success) + "\n")
			}
			// Show navigation hints
			content.WriteString("\n\n" + DimStyle.Render("Use ↑↓ to navigate, enter to select, ctrl+h to go back, q / esc to quit"))
		}
	}

	return ContainerStyle.Render(content.String())
}
