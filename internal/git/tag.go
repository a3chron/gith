package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// TODO: not yet used anywhere
func GetAllTagsWithCommit() (string, error) {
	out, err := exec.Command("git", "tag", "-n", "-l", "--sort=-v:refname").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tags: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func GetAllTags() (string, error) {
	out, err := exec.Command("git", "tag", "-l", "--sort=-v:refname").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tags: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func GetNLatestTags(n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid number of tags requested: %d", n)
	}

	out, err := GetAllTags()
	if err != nil {
		return "", err
	}

	if out == "" {
		return "", nil
	}

	outArr := strings.Split(out, "\n")
	if len(outArr) <= n {
		return out, nil
	}

	result := strings.Join(outArr[:n], "\n")
	return result, nil
}
