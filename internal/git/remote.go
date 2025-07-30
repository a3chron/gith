package git

import (
	"os/exec"
)

func GetRemotes() (string, string) {
	out, err := exec.Command("git", "remote", "-v").CombinedOutput()
	if err != nil {
		return "Failed to list remotes", string(out)
	}
	return string(out), ""
}

func GetRemotesClean() (string, string) {
	out, err := exec.Command("git", "remote").CombinedOutput()
	if err != nil {
		return "Failed to list remotes", string(out)
	}
	return string(out), ""
}

func RemoveRemote(remote string) (string, string) {
	out, err := exec.Command("git", "remote", "remove", remote).CombinedOutput()
	if err != nil {
		return "Failed to remove remote", string(out)
	}
	return string(out), ""
}
