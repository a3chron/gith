package git

import (
	"fmt"
	"os/exec"
)

func UndoLastCommit() (string, error) {
	out, err := exec.Command("git", "reset", "--soft", "HEAD~1").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("failed to undo last commit: %w", err)
	}
	return string(out), nil
}

func CommitStaged(commitMessage string) (string, error) {
	out, err := exec.Command("git", "commit", "-m", commitMessage).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("failed to commit staged: %w", err)
	}
	return string(out), nil
}

func CommitAll(commitMessage string) (string, error) {
	out, err := exec.Command("git", "commit", "-am", commitMessage).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("failed to commit staged: %w", err)
	}
	return string(out), nil
}
