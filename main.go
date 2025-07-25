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
		Options:     []string{"Switch Branch", "Create Branch", "Delete Branch", "Status", "Options"},
		Selected:    0,
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
