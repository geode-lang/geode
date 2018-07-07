package util

import (
	"os/exec"
)

// RunCommand executes a command and returns stdout from it.
func RunCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err

	}
	return out, nil
}
