package util

import (
	"os/exec"
)

// RunCommand executes a command and returns stdout from it.
func RunCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
		// log.Fatal("failed to run command `%s %s`: `%s`\n\n%s",
		// 	command, strings.Join(args, " "),
		// 	err.Error(), string(out))
	}
	return out, nil
}
