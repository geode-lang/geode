package util

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/nickwanninger/geode/pkg/util/log"
)

// RunCommand executes a command and returns stdout from it.
func RunCommand(command string, args ...string) ([]byte, error) {
	var out []byte
	var err error

	tmpcmd := command + " " + strings.Join(args, " ")
	maxLen := 500
	if len(tmpcmd) > maxLen {
		tmpcmd = tmpcmd[:maxLen-3] + "..."
	}
	title := fmt.Sprintf("Command Execution (%s)", tmpcmd)
	log.Timed(title, func() {
		cmd := exec.Command(command, args...)
		out, err = cmd.CombinedOutput()
	})

	if err != nil {
		return out, err

	}
	return out, nil
}
