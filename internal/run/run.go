/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package run

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lance-security/praetor/internal/output"
)

// RunCmd will run the provided command from the arguments
// will soon have support for a "modules" system designed to read specific CLI tools (e.g nmap) and append a "result" event to engagement log
func RunCmd(args []string) error {
	output.LogTask("Starting command execution")
	stopLoader := output.StartLoader("cmd-execution", fmt.Sprintf("Executing: %s", args[0]))

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		stopLoader(output.LevelError, output.IconReject, fmt.Sprintf("Command failed: %v", err))
		return fmt.Errorf("command failed: %w", err)
	}

	stopLoader(output.LevelPrimary, output.IconAccept, "Command execution completed")
	return nil
}
