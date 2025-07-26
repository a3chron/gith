package git

import "os/exec"

func GetTags() (string, error) {
	out, err := exec.Command("git", "tag", "-l", "--sort=v:refname").Output()
	return string(out), err
}
