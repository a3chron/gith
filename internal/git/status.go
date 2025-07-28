package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsWorkingTreeClean() (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, fmt.Errorf("failed to check working tree status: %w", err)
	}
	return len(strings.TrimSpace(string(out))) == 0, nil
}

func GetStatusInfo() (map[string][]string, int, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get status: %w", err)
	}

	length := len(strings.Split(strings.TrimSpace(string(out)), "\n"))

	status := make(map[string][]string)
	status["modified"] = []string{}
	status["added"] = []string{}
	status["deleted"] = []string{}
	status["untracked"] = []string{}

	lines := strings.SplitSeq(strings.TrimSpace(string(out)), "\n")
	for line := range lines {
		if len(line) < 3 {
			continue
		}

		statusCode := strings.TrimSpace(line[:2])
		filename := strings.TrimSpace(line[2:])

		switch {
		case strings.Contains(statusCode, "M"):
			status["modified"] = append(status["modified"], filename)
		case strings.Contains(statusCode, "A"):
			status["added"] = append(status["added"], filename)
		case strings.Contains(statusCode, "D"):
			status["deleted"] = append(status["deleted"], filename)
		case strings.Contains(statusCode, "?"):
			status["untracked"] = append(status["untracked"], filename)
		}
	}

	return status, length, nil
}
