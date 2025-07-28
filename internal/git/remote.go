package git

import (
	"fmt"
	"os/exec"
)

func ListRemotes() (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return "", fmt.Errorf("failed to list remotes: %w", err)
	}
	return string(out), nil
}
