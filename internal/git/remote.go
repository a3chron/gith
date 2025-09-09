package git

import (
	"os/exec"
	"strings"
)

func GetRemotes() (string, string) {
	out, err := exec.Command("git", "remote", "-v").CombinedOutput()
	if err != nil {
		return "Failed to list remotes", string(out)
	}
	return string(out), ""
}

func GetNumRemotes() (int, string) {
	out, err := exec.Command("git", "remote").CombinedOutput()
	if err != nil {
		return 0, string(out)
	}
	return len(strings.Split(string(out), "\n")), ""
}

func FormatFancyRemotes(remotes string) string {
	if remotes == "" {
		return remotes
	}

	lines := strings.Split(remotes, "\n")
	remoteInfo := make(map[string]struct {
		url   string
		fetch bool
		push  bool
	})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		name := parts[0]
		url := parts[1]
		op := strings.Trim(parts[2], "()")

		info := remoteInfo[name]
		info.url = url

		switch op {
		case "fetch":
			info.fetch = true
		case "push":
			info.push = true
		}

		remoteInfo[name] = info
	}

	var result []string
	for name, info := range remoteInfo {
		line := "\\ct" + name + "\n" + info.url

		if info.fetch && info.push {
			// Both operations - no suffix
		} else if info.fetch {
			line += "\n\\cwreadonly"
		} else if info.push {
			line += "\n\\cwwriteonly"
		}

		result = append(result, line+"\n-")
	}

	return strings.Join(result, "\n")
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
