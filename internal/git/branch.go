package git

import (
	"os/exec"
	"strings"
)

func GetBranches() ([]string, error) {
	out, err := exec.Command("git", "branch", "-a").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var branches []string
	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip current branch (has * prefix)
		if strings.HasPrefix(line, "* ") {
			continue
		}

		// Handle remote branches
		if strings.HasPrefix(line, "remotes/origin/") {
			branch := strings.TrimPrefix(line, "remotes/origin/")
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
