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

func initialModel(cfg *config.Config) internal.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.AccentStyle

	return internal.Model{
		CurrentStep: internal.StepLoad,
		Loading:     true,
		ActionModel: internal.ActionModel{
			Actions: []string{"Branch", "Status", "Commit", "Tag", "Remote", "Changes", "Options"},
		},
		BranchModel: internal.BranchModel{
			Actions: []string{"Switch Branch", "Create Branch", "List Branches", "Delete Branch"},
			Options: []string{"feat/", "fix/", "refactor/", "docs/", "Manual Input"},
		},
		CommitModel: internal.CommitModel{
			Actions:        []string{"Commit Staged", "Commit All", "Undo Last Commit"},
			CommitPrefixes: []string{"feat", "fix", "chore", "build", "ci", "test", "perf", "refactor", "revert", "style", "docs", "Custom Prefix"},
		},
		TagModel: internal.TagModel{
			Actions: []string{"Add Tag", "Remove Tag", "List Tags", "Push Tag"},
		},
		RemoteModel: internal.RemoteModel{
			Actions: []string{"List Remotes", "Add Remote", "Remove Remote"},
		},
		ConfigModel: internal.ConfigModel{
			Actions:        []string{"Change Flavor", "Change Accent", "Fetch at Init Behaviour", "Reset to Defaults"},
			InitBehaviours: []string{"Always fetch on Init", "Do not fetch for Quick Selects", "Never fetch"},
		},
		CurrentConfig: cfg,
		Selected:      0,
		Level:         0,
		StartAt:       "",
		StartAtLevel:  0,
		Spinner:       s,
	}
}

// initialModelWithStart returns the base model but allows setting a quick-start flow.
func initialModelWithStart(startAt string, startAtLevel int, cfg *config.Config) internal.Model {
	m := initialModel(cfg)
	m.StartAt = startAt
	m.StartAtLevel = startAtLevel
	return m
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
			Accent:        config.DefaultConfig.Accent,
			Flavor:        config.DefaultConfig.Flavor,
			InitBehaviour: config.DefaultConfig.InitBehaviour,
		}
	}

	// Initialize styles with the loaded config
	ui.UpdateStylesByConfig(cfg)

	isRepo, err := internal.IsGitRepository()
	if err != nil || !isRepo {
		return fmt.Errorf("not in a git repository")
	}

	p := tea.NewProgram(initialModel(cfg))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}

// runQuick starts the UI directly at a specific flow (e.g., add-tag).
func runQuick(startAt string, startAtLevel int) error {
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

	p := tea.NewProgram(initialModelWithStart(startAt, startAtLevel, cfg))
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

	case "update":
		if version == "dev" {
			// installed via go install (or actual dev build lol)
			err := internal.Update()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			return nil
		}

		fmt.Fprintf(os.Stderr, "To update via 'gith update' install gith via go\nInstuctions are linked in the repo's readme\n")
		os.Exit(1)

	case "config":
		return handleConfigCommand()

	case "completion":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: gith completion <shell>\nSupported shells: bash, zsh, fish\n")
			os.Exit(1)
		}
		if err := internal.GenerateCompletions(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return nil

	case "add":
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Usage: gith add remote\n")
			os.Exit(1)
		}
		switch os.Args[2] {
		case "remote":
			return runQuick("add-remote", 3)
		default:
			os.Exit(1)
			return fmt.Errorf("unknown command: %s\nUse 'gith help' for usage information", os.Args[1]+" "+os.Args[2])
		}

	case "tag":
		if len(os.Args) != 2 {
			fmt.Fprintf(os.Stderr, "Usage: gith tag\n")
			os.Exit(1)
		}

		// Start UI directly at tag selection
		return runQuick("add-tag", 3)

	case "push":
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Usage: gith push tag\n")
			os.Exit(1)
		}
		switch os.Args[2] {
		case "tag":
			// Start UI directly at push tag selection
			return runQuick("push-tag", 3)
		default:
			os.Exit(1)
			return fmt.Errorf("unknown command: %s\nUse 'gith help' for usage information", os.Args[1]+" "+os.Args[2])
		}

	case "commit":
		if len(os.Args) >= 4 {
			fmt.Fprintf(os.Stderr, "Usage: gith commit [ all | staged ]\n")
			os.Exit(1)
		}
		if len(os.Args) == 2 {
			return runQuick("commit-staged", 3)
		}

		switch os.Args[2] {
		case "all":
			return runQuick("commit-all", 3)
		case "staged":
			return runQuick("commit-staged", 3)
		} // TODO: all error here if none matched

	case "undo":
		if len(os.Args) != 3 || os.Args[2] != "commit" {
			fmt.Fprintf(os.Stderr, "Usage: gith undo commit\n")
			os.Exit(1)
		}

		return runQuick("undo-commit", 2)

	case "status":
		if len(os.Args) != 2 {
			fmt.Fprintf(os.Stderr, "Usage: gith status\n")
			os.Exit(1)
		}

		return runQuick("status", 1)
	}

	fmt.Fprintf(os.Stderr, "unknown command: %s\nUse 'gith help' for usage information", strings.Join(os.Args[1:], " "))
	os.Exit(1)
	return fmt.Errorf("unknown command")
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
	helpText := `
Config commands:
  gith config show     - Show current configuration
  gith config reset    - Reset configuration to defaults
  gith config path     - Show configuration file path

  gith config update [--flavor=<flavor>] [--accent=<accent>] [--initFetch=<initFetch>]
    Update your configuration options. Flags are optional and can be combined.

    --flavor=<flavor>
        Set the Catppuccin flavor (case insensitive). Available options:
        Latte, Frappe, Macchiato, Mocha

    --accent=<accent>
        Set the Catppuccin accent color (case insensitive). Available options:
        Rosewater, Flamingo, Pink, Mauve, Red, Maroon, Peach, Yellow, Green,
        Teal, Blue, Sapphire, Sky, Lavender, Gray

    --initFetch=<initFetch>
        Set fetch behavior on init. Available options:
        always  - always fetch on init
        quick   - fetch only for full load, skip for quick selects
        never   - never fetch on init
`

	fmt.Println(helpText)

	return nil
}

func showConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Current configuration:\n")
	fmt.Printf("  Flavor:         %s\n", cfg.Flavor)
	fmt.Printf("  Accent:         %s\n", cfg.Accent)
	fmt.Printf("  Init Behaviour: %s\n", cfg.InitBehaviour)
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

	args := os.Args[3:]
	for _, arg := range args {
		arg = strings.ToLower(arg)

		switch {
		case strings.HasPrefix(arg, "--flavor="):
			val := strings.TrimPrefix(arg, "--flavor=")
			if !config.IsValidFlavor(val) {
				return fmt.Errorf("not a valid flavor: %s\nvalid flavors: %s", val, strings.Join(config.GetAvailableFlavors(), ", "))
			}
			cfg.Flavor = val

		case strings.HasPrefix(arg, "--accent="):
			val := strings.TrimPrefix(arg, "--accent=")
			if !config.IsValidAccent(val) {
				return fmt.Errorf("not a valid accent: %s\nvalid accents: %s", val, strings.Join(config.GetAvailableAccents(), ", "))
			}
			cfg.Accent = val

		case strings.HasPrefix(arg, "--initfetch="):
			val := strings.TrimPrefix(arg, "--initfetch=")
			valid := map[string]bool{"always": true, "quick": true, "never": true}
			if !valid[val] {
				return fmt.Errorf("not a valid initFetch: %s\nvalid options: always, quick, never", val)
			}
			cfg.InitBehaviour = val

		default:
			return fmt.Errorf("'%s' is not a valid flag\nRun 'gith config help' to see valid flags", strings.Split(arg, "=")[0])
		}
	}

	saveErr := config.SaveConfig(cfg)
	if saveErr != nil {
		return fmt.Errorf("failed to save config: %w", saveErr)
	}
	fmt.Println("Configuration updated")

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
