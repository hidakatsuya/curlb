package runner

import (
	"os"
	"os/exec"
)

func RunCurl(args []string) error {
	cmd := exec.Command("curl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
