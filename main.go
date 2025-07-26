package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kurtschambach/gith/internal"
	ui "github.com/kurtschambach/gith/internal/ui"
)

var (
	Version = "dev"
)

func initialModel() ui.Model {
	return ui.Model{
		CurrentStep: ui.StepAction,
		ActionModel: ui.ActionModel{
			Actions: []string{"Branch", "Status", "Commit", "Tag", "Remote", "Changes", "Options"},
		},
		BranchModel: ui.BranchModel{
			Actions: []string{"Switch Branch", "Create Branch", "Delete Branch"},
		},
		CommitModel: ui.CommitModel{
			Actions: []string{"Undo Last Commit"},
		},
		TagModel: ui.TagModel{
			Actions: []string{"Add Tag", "Remove Tag", "List Tags", "Push Tag"},
		},
		RemoteModel: ui.RemoteModel{
			Actions: []string{"List Remotes", "Add Remote"},
		},
		ConfigModel: ui.ConfigModel{
			Accents: []string{"Rosewater", "Flamingo", "Pink", "Mauve", "Red", "Maroon", "Peach", "Yellow", "Green", "Teal", "Blue", "Sapphire", "Sky", "Lavender", "Gray"},
		},
		Selected: 0,
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v", "--version":
			if len(os.Args) == 3 && os.Args[2] == "check" {
				internal.PrintVersion(true, Version)
			} else {
				internal.PrintVersion(false, Version)
			}
			return
		case "help", "-h", "--help":
			internal.PrintHelp()
			return
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Use 'gith help' for usage information")
			os.Exit(1)
		}
	}

	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", string(out))
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
