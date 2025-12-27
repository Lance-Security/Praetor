/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lance-security/praetor/internal/output"
)

// declaring these as variables to allow mocking in tests in the future
var (
	lookPath    = exec.LookPath
	execCommand = exec.Command
	osStat      = os.Stat
	filepathAbs = filepath.Abs
)

type Bastion struct {
	ProjectDir string
	Command    []string
	AllowNet   bool
}

// RunInBastion checks for installation of bubblewrap, asks user
// to install it if not found, sets up a bubblewrap environment,
// then runs the command inside
func RunInBastion(args []string) error {
	output.LogTask("Setting up sandboxed Bastion environment")

	b := Bastion{
		ProjectDir: ".",
		Command:    args,
		AllowNet:   false,
	}

	output.Indent()
	defer output.Dedent()

	output.LogStep("Checking bubblewrap installation")
	if err := CheckAndInstallBubblewrap(); err != nil {
		output.LogError(err.Error())
		return err
	}
	output.LogSuccess("Bubblewrap is installed")

	output.LogStep("Resolving project directory")
	absProjectDir, err := filepathAbs(b.ProjectDir)
	if err != nil {
		output.LogErrorf("Failed to resolve project directory: %v", err)
		return fmt.Errorf("failed to resolve project dir: %w", err)
	}
	output.LogSuccess("Project directory resolved")

	output.LogStep("Configuring sandbox environment")
	bwrapArgs := []string{
		"--ro-bind", "/usr", "/usr",
		"--ro-bind", "/lib", "/lib",
		"--ro-bind", "/bin", "/bin",
		"--proc", "/proc",
		"--dev", "/dev",
		"--tmpfs", "/tmp",
		"--bind", absProjectDir, "/engagement", // mount engagement folder
		"--chdir", "/engagement",
	}

	// many distros use /lib64 and/or /usr/lib64 for the dynamic loader.
	if _, err := osStat("/lib64"); err == nil {
		bwrapArgs = append(bwrapArgs, "--ro-bind", "/lib64", "/lib64")
	}
	if _, err := osStat("/usr/lib64"); err == nil {
		bwrapArgs = append(bwrapArgs, "--ro-bind", "/usr/lib64", "/usr/lib64")
	}

	if _, err := osStat("/etc"); err == nil {
		bwrapArgs = append(bwrapArgs, "--ro-bind", "/etc", "/etc")
	}

	if !b.AllowNet {
		bwrapArgs = append(bwrapArgs, "--unshare-net")
		output.LogStep("Network isolation enabled")
	}

	output.LogStep("Bastion configuration complete")
	output.Dedent()

	bwrapArgs = append(bwrapArgs, b.Command...)
	cmd := execCommand("bwrap", bwrapArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	output.LogTask("Starting command execution in bastion")
	stopLoader := output.StartLoader("bastion-execution", fmt.Sprintf("Executing: %s", b.Command[0]))

	if err := cmd.Run(); err != nil {
		stopLoader(output.LevelError, output.IconReject, fmt.Sprintf("Bastion command failed: %v", err))
		return fmt.Errorf("bastion command failed: %w", err)
	}

	stopLoader(output.LevelPrimary, output.IconAccept, "Bastion command completed successfully")
	return nil
}

// CheckAndInstallBubblewrap checks if bubblewrap is installed,
// and if not, prompts the user to install it.
// TODO: implement auto installation
func CheckAndInstallBubblewrap() error {
	_, err := lookPath("bwrap")
	if err == nil {
		// bubblewrap is installed
		return nil
	}
	return fmt.Errorf("bubblewrap is not installed. please install with `sudo apt install bubblewrap` or `sudo pacman -S bubblewrap` and try again.")
}
