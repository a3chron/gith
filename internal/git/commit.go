package git

import "os/exec"

func UndoLastCommit() (string, error) {
	out, err := exec.Command("git", "reset", "--soft HEAD~1").Output()
	return string(out), err
}
