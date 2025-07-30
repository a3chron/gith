package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/a3chron/gith/internal"
	ui "github.com/a3chron/gith/internal/ui"
)

var (
	Version = "dev (or version injection still broken)"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

func initialModel() ui.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.AccentStyle

	return ui.Model{
		CurrentStep: ui.StepLoad,
		Loading:     true,
		ActionModel: ui.ActionModel{
			Actions: []string{"Branch", "Status", "Commit", "Tag", "Remote", "Changes", "Options"},
		},
		BranchModel: ui.BranchModel{
			Actions: []string{"Switch Branch", "Create Branch", "List Branches", "Delete Branch"},
		},
		CommitModel: ui.CommitModel{
			Actions: []string{"Undo Last Commit"},
		},
		TagModel: ui.TagModel{
			Actions: []string{"Add Tag", "Remove Tag", "List Tags", "Push Tag"},
		},
		RemoteModel: ui.RemoteModel{
			Actions: []string{"List Remotes", "Add Remote", "Remove Remote"},
		},
		ConfigModel: ui.ConfigModel{
			Accents: []string{"Rosewater", "Flamingo", "Pink", "Mauve", "Red", "Maroon", "Peach", "Yellow", "Green", "Teal", "Blue", "Sapphire", "Sky", "Lavender", "Gray"},
			Flavors: []string{"Latte", "Frappe", "Macchiato", "Mocha"},
		},
		Selected: 0,
		Level:    0,
		Spinner:  s,
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) > 1 {
		return handleCliArgs()
	}

	isRepo, err := internal.IsGitRepository()
	if err != nil || !isRepo {
		return fmt.Errorf("not in a git repository")
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}

func handleCliArgs() error {
	switch os.Args[1] {
	case "version", "-v", "--version":
		checkForUpdate := len(os.Args) == 3 && os.Args[2] == "check"
		internal.PrintVersion(checkForUpdate, Version, Commit, Date, BuiltBy)
		return nil

	case "help", "-h", "--help":
		internal.PrintHelp()
		return nil

	default:
		return fmt.Errorf("unknown command: %s\nUse 'gith help' for usage information", os.Args[1])
	}
}
