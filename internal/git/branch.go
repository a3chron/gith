package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetBranches() ([]string, error) {
	out, err := exec.Command("git", "branch", "-a").Output()
	if err != nil {
		return []string{"Failed to get branches"}, fmt.Errorf("%w", err)
	}

	lines := strings.Split(string(out), "\n")
	var branches []string
	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		currentBranch, err := GetCurrentBranch()
		if err != nil {
			return []string{currentBranch}, fmt.Errorf("%w", err)
		}

		// Skip current branch
		if strings.HasPrefix(line, "* ") || line == currentBranch {
			continue
		}

		// Handle remote branches
		if after, hasPrefix := strings.CutPrefix(line, "remotes/origin/"); hasPrefix {
			branch := after
			// Skip HEAD pointer
			if strings.Contains(branch, "HEAD ->") {
				continue
			}
			if branch != "" && !seen[branch] {
				branches = append(branches, branch)
				seen[branch] = true
			}
		} else {
			// Local branches
			if line != "" && !seen[line] {
				branches = append(branches, line)
				seen[line] = true
			}
		}
	}

	return branches, nil
}

func GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "Failed to get current branch", fmt.Errorf("%w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func CreateBranch(branchName string) (string, error) {
	out, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()
	if err != nil {
		return string(out), err
	} else {
		return string(out), nil
	}
}

func SwitchBranch(selectedBranch string) (string, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+selectedBranch)
	err := cmd.Run()

	var switchCmd *exec.Cmd
	if err != nil {
		// Branch doesn't exist locally, create and track remote
		switchCmd = exec.Command("git", "checkout", "-b", selectedBranch, "origin/"+selectedBranch)
	} else {
		// Branch exists locally, just switch
		switchCmd = exec.Command("git", "switch", selectedBranch)
	}

	out, err := switchCmd.CombinedOutput()
	if err != nil {
		return string(out), err
	} else {
		return string(out), nil
	}
}

func DeleteBranch(selectedBranch string) (string, error) {
	// Try to delete local branch first
	cmd := exec.Command("git", "branch", "-d", selectedBranch)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out) + "\n" + "\\cpTo force delete use: \ngit branch -D" + selectedBranch, err
	} else {
		return string(out), nil
	}
}
