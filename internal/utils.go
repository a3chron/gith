package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func UpdateRepo() error {
	if err := exec.Command("git", "fetch", "--all", "--prune").Run(); err != nil {
		return fmt.Errorf("failed to fetch remotes: %w", err)
	}
	return nil
}

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Orange = "\033[33m"
	Red    = "\033[31m"
	Gray   = "\033[90m"
	Bold   = "\033[1m"
)

const (
	TopLeft     = "╭"
	TopRight    = "╮"
	BottomLeft  = "╰"
	BottomRight = "╯"
	Horizontal  = "─"
	Vertical    = "│"
)

type VersionComparison struct {
	Current     string
	Latest      string
	Difference  string // "same", "patch", "minor", "major", "dev", "invalid"
	Color       string
	Message     string
	ShowUpgrade bool
}

func Update() error {
	out, err := exec.Command("go", "install", "github.com/a3chron/gith@latest").CombinedOutput()

	fmt.Println(string(out))

	if err != nil {
		return fmt.Errorf("failed to update gith: %w", err)
	}

	return nil
}

func PrintVersion(checkForUpdate bool, version string, commit string, date string, builtBy string) {
	// Create the version box
	printVersionBox(version, commit, date, builtBy)

	if checkForUpdate {
		if comparison, err := checkVersionComparison(version); err != nil {
			printUpdateBox("Error checking for updates", Gray, false, version == "dev") // TODO: +err.Error() for message, but can be long, maybe add verbosity and just print
		} else {
			printUpdateBox(comparison.Message, comparison.Color, comparison.ShowUpgrade, version == "dev")
		}
	}
}

func maxInt(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func printVersionBox(version, commit, date, builtBy string) {
	// Calculate max width for the box
	title := fmt.Sprintf("gith version: %s", version)
	author := "by a3chron (https://a3chron.vercel.app/)"
	commitLine := fmt.Sprintf("commit: %s", commit)
	dateLine := fmt.Sprintf("date: %s", date)
	builtLine := fmt.Sprintf("built with %s", builtBy)

	maxWidth := maxInt(len(title), len(author), len(commitLine), len(dateLine), len(builtLine)) + 4

	// Top border
	fmt.Printf("%s%s%s\n", TopLeft, strings.Repeat(Horizontal, maxWidth-2), TopRight)

	// Title line
	fmt.Printf("%s %s%-*s %s\n", Vertical, Bold, maxWidth-4, title, Reset+Vertical)

	// Author line
	fmt.Printf("%s %s%-*s %s\n", Vertical, Gray, maxWidth-4, author, Reset+Vertical)

	// Separator line
	fmt.Printf("%s%s%s\n", Vertical, strings.Repeat(Horizontal, maxWidth-2), Vertical)

	// Version info lines
	fmt.Printf("%s %-*s %s\n", Vertical, maxWidth-4, commitLine, Vertical)
	fmt.Printf("%s %-*s %s\n", Vertical, maxWidth-4, dateLine, Vertical)
	fmt.Printf("%s %-*s %s\n", Vertical, maxWidth-4, builtLine, Vertical)

	// Bottom border
	fmt.Printf("%s%s%s\n", BottomLeft, strings.Repeat(Horizontal, maxWidth-2), BottomRight)
}

func printUpdateBox(message, color string, showUpgrade bool, goInstallation bool) {
	lines := []string{}
	for line := range strings.SplitSeq(message, "\n") {
		lines = append(lines, line)
	}

	if showUpgrade {
		lines = append(lines, "")
		if goInstallation {
			// installed via go install (or actual dev build)
			lines = append(lines, "To install the latest version run:")
			lines = append(lines, "gith update OR go install github.com/a3chron/gith@latest")
		} else {
			lines = append(lines, "To install the latest version you will have to")
			lines = append(lines, "download the latest binaries from github and install them")
		}
	}

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 4

	// Top border
	fmt.Printf("\n%s%s%s\n", TopLeft, strings.Repeat(Horizontal, maxWidth-2), TopRight)

	// Content lines
	for _, line := range lines {
		fmt.Printf("%s %s%-*s%s %s\n", Vertical, color, maxWidth-4, line, Reset, Vertical)
	}

	// Bottom border
	fmt.Printf("%s%s%s\n", BottomLeft, strings.Repeat(Horizontal, maxWidth-2), BottomRight)
}

func checkVersionComparison(currentVersion string) (*VersionComparison, error) {
	// Check if current version contains "dev"
	if strings.Contains(strings.ToLower(currentVersion), "dev") {
		latestOfficial, err := getLatestOfficialVersion()
		if err != nil {
			return &VersionComparison{
				Current:     currentVersion,
				Latest:      "unknown",
				Difference:  "dev",
				Color:       Green,
				Message:     fmt.Sprintf("Development version: %s", currentVersion),
				ShowUpgrade: false,
			}, nil
		}
		return &VersionComparison{
			Current:     currentVersion,
			Latest:      latestOfficial,
			Difference:  "dev",
			Color:       Green,
			Message:     fmt.Sprintf("Development version: %s\nLatest official: %s", currentVersion, latestOfficial),
			ShowUpgrade: false,
		}, nil
	}

	// Get latest version from GitHub
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/a3chron/gith/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var latestRelease githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&latestRelease); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
	currentClean := strings.TrimPrefix(currentVersion, "v")

	comparison := compareVersions(currentClean, latestVersion)
	return comparison, nil
}

func getLatestOfficialVersion() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/a3chron/gith/releases/latest")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var latestRelease githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&latestRelease); err != nil {
		return "", err
	}

	return strings.TrimPrefix(latestRelease.TagName, "v"), nil
}

func compareVersions(current, latest string) *VersionComparison {
	// Parse versions
	currentMajor, currentMinor, currentPatch, _, err1 := ParseSemanticVersion(current)
	latestMajor, latestMinor, latestPatch, _, err2 := ParseSemanticVersion(latest)

	if err1 != nil || err2 != nil {
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "invalid",
			Color:       Gray,
			Message:     "Unable to compare versions",
			ShowUpgrade: false,
		}
	}

	// Compare versions
	if currentMajor < latestMajor {
		// Major version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "major",
			Color:       Red,
			Message:     fmt.Sprintf("Major update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentMajor > latestMajor {
		// Newer major version
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "same",
			Color:       Green,
			Message:     fmt.Sprintf("You have a newer version: %s (latest: %s)", current, latest),
			ShowUpgrade: false,
		}
	}

	if currentMinor < latestMinor {
		// Minor version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "minor",
			Color:       Orange,
			Message:     fmt.Sprintf("Minor update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentMinor > latestMinor {
		// Newer minor version
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "same",
			Color:       Green,
			Message:     fmt.Sprintf("You have a newer version: %s (latest: %s)", current, latest),
			ShowUpgrade: false,
		}
	}

	if currentPatch < latestPatch {
		// Patch version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "patch",
			Color:       Green,
			Message:     fmt.Sprintf("Patch update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentPatch > latestPatch {
		// Newer patch version
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "same",
			Color:       Green,
			Message:     fmt.Sprintf("You have a newer version: %s (latest: %s)", current, latest),
			ShowUpgrade: false,
		}
	}

	// Same version
	return &VersionComparison{
		Current:     current,
		Latest:      latest,
		Difference:  "same",
		Color:       Green,
		Message:     "You are using the latest version",
		ShowUpgrade: false,
	}
}

func ParseSemanticVersion(tag string) (major, minor, patch int, prefix string, err error) {
	prefix = ""
	if strings.HasPrefix(tag, "v") {
		prefix = "v"
		tag = tag[1:]
	}

	parts := strings.Split(tag, ".")
	if len(parts) != 3 {
		return 0, 0, 0, prefix, fmt.Errorf("invalid semantic version format")
	}

	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, prefix, err
	}

	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, prefix, err
	}

	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, prefix, err
	}

	return major, minor, patch, prefix, nil
}

func GenerateNextSemanticVersion(currentTag, versionType string) (string, error) {
	major, minor, patch, prefix, err := ParseSemanticVersion(currentTag)
	if err != nil {
		return "", err
	}

	switch versionType {
	case "Patch":
		patch++
	case "Minor":
		minor++
		patch = 0
	case "Major":
		major++
		minor = 0
		patch = 0
	default:
		return "", fmt.Errorf("invalid version type")
	}

	return fmt.Sprintf("%s%d.%d.%d", prefix, major, minor, patch), nil
}

func PrintHelp() {
	fmt.Print(`gith - TUI Git Helper

Usage:
  gith                   Start interactive mode

  gith version           Show version information  
  gith version check     Show version & check for updates
  gith update			 Update to the latest version (for go users)

  // Quick Selects - Jump right to a specific selection
  gith tag             	 Add Tag
  gith push tag          Push Tag
  gith commit [staged]   Commit Staged
  gith commit all        Commit All Files
  gith undo commit       Undo Last Commit
  gith add remote        Add Remote
  gith status            Show Status

  gith config help       Show config related help message
  gith help              Show this help message

Interactive Commands:
  ↑↓ or j/k              Navigate menu items
  Enter                  Select item
  Ctrl+H, Ctrl+Y         Go back to previous step
  Q/Esc                  Quit application

Features:
  • Branch management
  • Commit operations
  • Tag management
  • Remote operations
  • Git status checking
  • Configuration options

Shell Completions:
  gith completion fish > ~/.config/fish/completions/gith.fish
  gith completion bash > ~/.local/share/bash-completion/completions/gith
  gith completion zsh > ~/.local/share/zsh/site-functions/_gith
`)
}

func GenerateCompletions(shell string) error {
	switch shell {
	case "fish":
		fmt.Print(`# Fish completion for gith
complete -c gith -f
complete -c gith -n "__fish_use_subcommand" -a "version update config help add push tag status undo commit" -d "Available commands"
complete -c gith -n "__fish_seen_subcommand_from version" -a "check" -d "Check for updates"
complete -c gith -n "__fish_seen_subcommand_from config" -a "show reset path update help" -d "Config commands"
complete -c gith -n "__fish_seen_subcommand_from config update" -l flavor -d "Catppuccin flavor" -a "latte frappe macchiato mocha"
complete -c gith -n "__fish_seen_subcommand_from config update" -l accent -d "Catppuccin accent" -a "rosewater flamingo pink mauve red maroon peach yellow green teal sky sapphire blue lavender gray"
complete -c gith -n "__fish_seen_subcommand_from config update" -l initFetch -d "Init fetch behaviour" -a "always quick never"
complete -c gith -n "__fish_seen_subcommand_from add" -a "remote" -d "Quick Select: Add Remote"
complete -c gith -n "__fish_seen_subcommand_from push" -a "tag" -d "Quick Select: Push Tag"
complete -c gith -n "__fish_seen_subcommand_from undo" -a "commit" -d "Quick Select: Status"
`)
	case "bash":
		fmt.Print(`# Bash completion for gith
_gith() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    case "${prev}" in
        gith)
            opts="version update config help add push tag status undo commit"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        version)
            opts="check"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        config)
            opts="show reset path update help"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
		add)
            opts="remote"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        push)
            opts="tag"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
		undo)
            opts="commit"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        --flavor)
            opts="latte frappe macchiato mocha"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        --accent)
            opts="rosewater flamingo pink mauve red maroon peach yellow green teal sky sapphire blue lavender"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
		--initFetch)
            opts="always quick never"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
    esac
}
complete -F _gith gith
`)
	case "zsh":
		fmt.Print(`#compdef gith
_gith() {
    local context state line
    _arguments \
        '1:command:(version update config help add push status undo commit)' \
        '*::arg:->args'
    
    case $state in
        args)
            case $words[1] in
                version)
                    _arguments '1:subcommand:(check)'
                    ;;
                config)
                    _arguments \
                        '1:subcommand:(show reset path update help tag)' \
                        '--flavor[Catppuccin flavor]:(latte frappe macchiato mocha)' \
                        '--accent[Catppuccin accent]:(rosewater flamingo pink mauve red maroon peach yellow green teal sky sapphire blue lavender)' \
						'--initFetch[Init fetch behaviour]:(always quick never)'
                    ;;
                add)
                    _arguments '1:subcommand:(remote)'
                    ;;
                push)
                    _arguments '1:subcommand:(tag)'
                    ;;
				undo)
                    _arguments '1:subcommand:(commit)'
                    ;;
            esac
            ;;
    esac
}
`)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	return nil
}

func IsGitRepository() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil, nil
}

// isInputStep returns true if the current step expects free-text input
func isInputStep(step Step) bool {
	switch step {
	case StepTagInput, StepBranchInput, StepRemoteNameInput, StepRemoteUrlInput, StepCommitInput:
		return true
	default:
		return false
	}
}

func (m *Model) OutputByLevel(out string) {
	// Fill output array up to current level
	for len(m.Output) <= m.Level {
		m.Output = append(m.Output, "")
	}
	m.Output[m.Level] += out
}
