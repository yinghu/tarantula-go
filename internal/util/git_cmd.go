package util

import "os/exec"

func GitCurBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
