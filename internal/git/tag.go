package git

import (
	"os/exec"
	"strings"
)

func GetAllTags() (string, error) {
	out, err := exec.Command("git", "tag", "-n", "-l", "--sort=-v:refname").Output()
	return string(out), err
}

func GetNLatestTags(n int) (string, error) {
	out, err := GetAllTags()
	if err != nil {
		return string(out), err
	}

	outArr := strings.Split(out, "\n")

	if len(outArr) <= n {
		return string(out), err
	}

	outArrCopy := make([]string, n)
	copy(outArrCopy, outArr)

	result := strings.Join(outArrCopy, "\n")

	return result, err
}
