package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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

func PrintVersion(checkForUpdate bool, version string) {
	fmt.Printf("gith version: %s\n", version)

	if checkForUpdate {
		if err := checkAndPrintUpdate(version); err != nil {
			fmt.Printf("Error checking for updates: %v\n", err)
		}
	}
}

func checkAndPrintUpdate(currentVersion string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/kurtschambach/gith/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to close response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var latestRelease githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&latestRelease); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	if latestRelease.TagName != currentVersion {
		fmt.Printf("New version available: %s (current: %s)\n", latestRelease.TagName, currentVersion)
		fmt.Println("To install the latest version run:\ngo install github.com/kurtschambach/gith@latest")
	} else {
		fmt.Println("You are using the latest version.")
	}

	return nil
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
