package util

import (
	"os/exec"
	"strings"

	"github.com/nickwanninger/geode/pkg/util/log"
)

// RunCommand executes a command and returns stdout from it.
func RunCommand(command string, args ...string) []byte {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("failed to run command `%s %s`: `%s`\n\n%s",
			command, strings.Join(args, " "),
			err.Error(), string(out))
	}
	return out
}
