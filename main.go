package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/a3chron/gith/internal"
	"github.com/a3chron/gith/internal/config"
	ui "github.com/a3chron/gith/internal/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func getVersion() string {
	// If version was set by GoReleaser, use it
	if version != "dev" {
		return version
	}

	// Fallback for go install: try to get version from build info
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" && info.Main.Version != "" {
			return info.Main.Version
		}
	}

	return "dev"
}

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
			Actions: []string{"Change Flavor", "Change Accent", "Reset to Defaults"},
			Flavors: config.GetAvailableFlavors(),
			Accents: config.GetAvailableAccents(),
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

	cfg, err := config.LoadConfig()
	if err != nil {
		// Fall back to defaults if config loading fails
		fmt.Fprintf(os.Stderr, "Warning: failed to load config, using defaults: %v\n", err)
		cfg = &config.Config{
			Accent: config.DefaultConfig.Accent,
			Flavor: config.DefaultConfig.Flavor,
		}
	}

	// Initialize styles with the loaded config
	ui.UpdateStylesByConfig(cfg)

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
		internal.PrintVersion(checkForUpdate, getVersion(), commit, date, builtBy)
		return nil

	case "help", "-h", "--help":
		internal.PrintHelp()
		return nil

	case "config":
		return handleConfigCommand()

	default:
		return fmt.Errorf("unknown command: %s\nUse 'gith help' for usage information", os.Args[1])
	}
}

func handleConfigCommand() error {
	if len(os.Args) < 3 {
		return printConfigUsage()
	}

	switch os.Args[2] {
	case "show":
		return showConfig()
	case "reset":
		return resetConfig()
	case "path":
		return showConfigPath()
	case "update":
		return updateConfig()
	default:
		return printConfigUsage()
	}
}

func printConfigUsage() error {
	fmt.Println("Config commands:")
	fmt.Println("  gith config show   	- Show current configuration")
	fmt.Println("  gith config reset  	- Reset configuration to defaults")
	fmt.Println("  gith config path   	- Show configuration file path")
	fmt.Println("  gith config update [--flavor=<flavor>] [--accent=<accent>]   - Update configuration")
	fmt.Println("    <flavor> can be any of the catppuccin flavors (case insensitive)")
	fmt.Println("    <accent> can be any of the catppuccin accents (case insensitive)")
	return nil
}

func showConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Current configuration:\n")
	fmt.Printf("  Flavor: %s\n", cfg.Flavor)
	fmt.Printf("  Accent: %s\n", cfg.Accent)
	return nil
}

func updateConfig() error {
	if len(os.Args) <= 3 {
		return printConfigUsage()
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// TODO: when adding more to configure, make sure to handle this automatically

	first := strings.ToLower(os.Args[3])
	second := ""

	if len(os.Args) >= 5 {
		second = strings.ToLower(os.Args[4])
	}

	if flavor, found := strings.CutPrefix(first, "--flavor="); found {
		if config.IsValidFlavor(flavor) {
			cfg.Flavor = flavor
		} else {
			return fmt.Errorf("not a valid flavor: %s\nlist of valid flavors: %s", flavor, config.GetAvailableFlavors())
		}
	} else if accent, found := strings.CutPrefix(first, "--accent="); found {
		if config.IsValidAccent(accent) {
			cfg.Accent = accent
		} else {
			return fmt.Errorf("not a valid accent: %s\nlist of valid accent: %s", accent, config.GetAvailableAccents())
		}
	} else {
		return fmt.Errorf("'%s' is not a valid flag\nRun 'gith config help' to see valid flags", strings.TrimPrefix(strings.Split(first, "=")[0], "--"))
	}

	if second != "" {
		if flavor, found := strings.CutPrefix(second, "--flavor="); found {
			if config.IsValidFlavor(flavor) {
				cfg.Flavor = flavor
			} else {
				return fmt.Errorf("not a valid flavor: %s\nlist of valid flavors: %s", flavor, config.GetAvailableFlavors())
			}
		} else if accent, found := strings.CutPrefix(second, "--accent="); found {
			if config.IsValidAccent(accent) {
				cfg.Accent = accent
			} else {
				return fmt.Errorf("not a valid accent: %s\nlist of valid accent: %s", accent, config.GetAvailableAccents())
			}
		} else {
			return fmt.Errorf("'%s' is not a valid flag\nRun 'gith config help' to see valid flags", strings.TrimPrefix(strings.Split(first, "=")[0], "--"))
		}
	}

	config.SaveConfig(cfg)

	return nil
}

func resetConfig() error {
	if err := config.SaveConfig(&config.DefaultConfig); err != nil {
		return fmt.Errorf("failed to reset config: %w", err)
	}
	fmt.Println("Configuration reset to defaults")
	return nil
}

func showConfigPath() error {
	path, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}
	fmt.Println(path)
	return nil
}
