/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lachlanharrisdev/praetor/internal/formats"
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
	formats.Info("Setting up sandboxed Bastion environment")

	b := Bastion{
		ProjectDir: ".",
		Command:    args,
		AllowNet:   false,
	}

	formats.Info("Checking bubblewrap installation")
	if err := CheckAndInstallBubblewrap(); err != nil {
		formats.Error(err.Error())
		return err
	}
	formats.Success("Bubblewrap is installed")

	formats.Info("Resolving project directory")
	absProjectDir, err := filepathAbs(b.ProjectDir)
	if err != nil {
		formats.Errorf("Failed to resolve project directory: %v", err)
		return fmt.Errorf("failed to resolve project dir: %w", err)
	}
	formats.Success("Project directory resolved")

	bwrapArgs := configureSandbox(b, absProjectDir)

	bwrapArgs = append(bwrapArgs, b.Command...)
	cmd := execCommand("bwrap", bwrapArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if len(b.Command) > 0 {
		formats.Infof("Executing in bastion: %s", b.Command[0])
	} else {
		formats.Info("Executing in bastion")
	}

	if err := cmd.Run(); err != nil {
		formats.Errorf("Bastion command failed: %v", err)
		return fmt.Errorf("bastion command failed: %w", err)
	}

	formats.Success("Bastion command completed")
	return nil
}

func configureSandbox(b Bastion, absProjectDir string) []string {
	formats.Info("Configuring sandbox environment")
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
		formats.Info("Network isolation enabled")
	}

	formats.Info("Bastion configuration complete")

	return bwrapArgs
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
