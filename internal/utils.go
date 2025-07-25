package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func UpdateRepo() {
	exec.Command("git", "fetch", "-p").Run()
	exec.Command("git", "fetch", "-a").Run()
}

func PrintVersion(checkForUpdate bool, Version string) {
	fmt.Printf("gith version %s\n", Version)
	if checkForUpdate {
		resp, err := http.Get("https://api.github.com/repos/kurtschambach/gith/releases/latest")
		if err != nil {
			fmt.Println("error fetching latest release:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("unexpected status:", resp.Status)
			return
		}

		var latestRelease githubRelease
		if err := json.NewDecoder(resp.Body).Decode(&latestRelease); err != nil {
			fmt.Println("error decoding JSON:", err)
			return
		}

		if latestRelease.TagName != Version {
			fmt.Printf("new version available: %s (current: %s)\n", latestRelease.TagName, Version)
			fmt.Print("To install the latest version run: go install github.com/kurtschambach/gith@latest")
		} else {
			fmt.Println("you are using the latest version.")
		}
	}
}

func PrintHelp() {
	fmt.Printf(`gith - TUI Git Helper

Usage:
  gith              	Start interactive mode
  gith version      	Show version information
  gith version check	Show version & check for updates
  gith help         	Show this help message

Interactive Commands:
  ↑↓ or j/k        Navigate menu items
  Enter            Select item
  Ctrl+H           Go back to previous step
  Q/Esc            Quit application

`)
}
