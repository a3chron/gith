package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
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

// Box drawing characters
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

func PrintVersion(checkForUpdate bool, version string, commit string, date string, builtBy string) {
	// Create the version box
	printVersionBox(version, commit, date, builtBy)

	if checkForUpdate {
		if comparison, err := checkVersionComparison(version); err != nil {
			printUpdateBox("Error checking for updates: "+err.Error(), Gray, false)
		} else {
			printUpdateBox(comparison.Message, comparison.Color, comparison.ShowUpgrade)
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

func printUpdateBox(message, color string, showUpgrade bool) {
	lines := []string{}
	for line := range strings.SplitSeq(message, "\n") {
		lines = append(lines, line)
	}

	if showUpgrade {
		lines = append(lines, "To install the latest version run:")
		lines = append(lines, "go install github.com/a3chron/gith@latest")
	}

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 4

	fmt.Println() // Empty line before update box

	// Top border
	fmt.Printf("%s%s%s\n", TopLeft, strings.Repeat(Horizontal, maxWidth-2), TopRight)

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	currentParts, err1 := parseVersion(current)
	latestParts, err2 := parseVersion(latest)

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
	if currentParts[0] < latestParts[0] {
		// Major version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "major",
			Color:       Red,
			Message:     fmt.Sprintf("Major update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentParts[0] > latestParts[0] {
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

	if currentParts[1] < latestParts[1] {
		// Minor version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "minor",
			Color:       Orange,
			Message:     fmt.Sprintf("Minor update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentParts[1] > latestParts[1] {
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

	if currentParts[2] < latestParts[2] {
		// Patch version difference
		return &VersionComparison{
			Current:     current,
			Latest:      latest,
			Difference:  "patch",
			Color:       Green,
			Message:     fmt.Sprintf("Patch update available: %s → %s", current, latest),
			ShowUpgrade: true,
		}
	} else if currentParts[2] > latestParts[2] {
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

func parseVersion(version string) ([3]int, error) {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	// Use regex to match semantic version pattern
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-.*)?(?:\+.*)?$`)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return [3]int{}, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return [3]int{}, err
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return [3]int{}, err
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return [3]int{}, err
	}

	return [3]int{major, minor, patch}, nil
}

func PrintHelp() {
	fmt.Print(`gith - TUI Git Helper

Usage:
  gith                 Start interactive mode
  gith version         Show version information  
  gith version check   Show version & check for updates
  gith help            Show this help message

Interactive Commands:
  ↑↓ or j/k            Navigate menu items
  Enter                Select item
  Ctrl+H               Go back to previous step
  Q/Esc                Quit application

Features:
  • Branch management
  • Commit operations
  • Tag management
  • Remote operations
  • Git status checking
  • Configuration options

`)
}

func IsGitRepository() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil, nil
}
