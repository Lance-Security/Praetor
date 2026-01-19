/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package run

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lachlanharrisdev/praetor/internal/formats"
)

// RunCmd will run the provided command from the arguments
// TODO: add support for a "modules" system designed to read specific CLI tools (e.g nmap) and append a "result" event to engagement log
func RunCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided to run")
	}

	formats.Infof("Executing: %s", args[0])

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		formats.Errorf("Command failed: %v", err)
		return fmt.Errorf("command failed: %w", err)
	}

	formats.Success("Command completed")
	return nil
}
