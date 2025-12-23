package bastions

import (
	"fmt"
	"os/exec"
)

type Bastion struct {
	ProjectDir string
	Command    []string
	AllowNet   bool
}

// RunInBastion checks for installation of bubblewrap, asks to
// install it if it is missing (y/n, then runs automatic install command),
// then runs the command inside the isolated environment.
func RunInBastion(args []string) error {
	b := Bastion{
		ProjectDir: ".",
		Command:    args,
		AllowNet:   false,
	}

	if err := CheckAndInstallBubblewrap(); err != nil {
		return err
	}

	bwrapArgs := []string{
		"--ro-bind", "/usr", "/usr",
		"--ro-bind", "/lib", "/lib",
		"--ro-bind", "/bin", "/bin",
		"--proc", "/proc",
		"--dev", "/dev",
		"--bind", b.ProjectDir, "/work", // Mount engagement folder
		"--chdir", "/work",
	}

	if !b.AllowNet {
		bwrapArgs = append(bwrapArgs, "--unshare-net")
	}

	// Append the actual command the user wants to run
	bwrapArgs = append(bwrapArgs, b.Command...)
	cmd := exec.Command("bwrap", bwrapArgs...)
	// ... handle PTY and logging here ...
	return cmd.Run()
}

// CheckAndInstallBubblewrap checks if bubblewrap is installed,
// and if not, prompts the user to install it.
// If the user agrees, it runs the installation command.
func CheckAndInstallBubblewrap() error {
	_, err := exec.LookPath("bwrap")
	if err == nil {
		// bubblewrap is installed
		return nil
	}
	return fmt.Errorf("Bubblewrap is not installed. Please install with `sudo apt install bubblewrap` or `sudo pacman -S bubblewrap` and try again.")
}
