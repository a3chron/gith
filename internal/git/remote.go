package git

import "os/exec"

func ListRemotes() (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()
	return string(out), err
}
