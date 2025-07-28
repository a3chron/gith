package git

import (
	"fmt"
	"os/exec"
)

func UndoLastCommit() (string, error) {
	out, err := exec.Command("git", "reset", "--soft", "HEAD~1").Output()
	if err != nil {
		return string(out), fmt.Errorf("failed to undo last commit: %w", err)
	}
	return string(out), nil
}
