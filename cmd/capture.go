/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/lance-security/praetor/internal/engagement"
	"github.com/lance-security/praetor/internal/events"
	"github.com/lance-security/praetor/internal/output"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture [filename]",
	Short: "Capture a tool output through a pipe or filename",
	Long: `Capture will allow you to capture the output of a tool, either
	through pipes (e.g. 'nmap 127.0.0.1 | pt capture') or from a file
	(e.g. 'pt capture tools/nmap_result_27-12_11-44.txt'), and store it
	as a result event in the engagement log.`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get engagement context
		return capture(args)
	},
}

// capture executes the capture command logic
func capture(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	engDir, err := engagement.FindEngagementDir(cwd)
	if err != nil {
		return fmt.Errorf("failed to find engagement directory: %w", err)
	}

	m, err := engagement.ReadMetadata(engDir)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	user := events.GetUser()

	var rawData string
	var sourceLabel string

	if len(args) > 0 {
		rawData, sourceLabel, err = captureFromFile(args)
		if err != nil {
			return err
		}
	} else {
		rawData, sourceLabel, err = captureFromStdin()
		if err != nil {
			return err
		}
	}

	if rawData == "" {
		output.LogWarning("No data captured")
		return fmt.Errorf("no data to capture")
	}

	event := events.NewEvent(
		"result",
		fmt.Sprintf("Captured output from %s", sourceLabel),
		time.Now().UTC().Format(time.RFC3339Nano),
		m.EngagementID,
		filepath.Clean(cwd),
		user,
		rawData,
		nil,
	)

	if err := events.AppendEvent(engagement.EventsPath(engDir), event); err != nil {
		return fmt.Errorf("failed to append event: %w", err)
	}

	output.LogSuccess("Successfully captured and logged result")
	return engagement.TouchLastUsed(engDir)
}

// captureFromFile
func captureFromFile(args []string) (rawData string, sourceLabel string, err error) {
	filename := args[0]
	data, err := readFromFile(filename)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}
	return data, filepath.Base(filename), nil
}

// captureFromStdin
func captureFromStdin() (rawData string, sourceLabel string, err error) {
	output.LogTask("Reading data from stdin")
	stopLoader := output.StartLoader("stdin-read", "Waiting for data...")

	data, err := readFromStdin()
	if err != nil {
		stopLoader(output.LevelError, output.IconReject, fmt.Sprintf("Failed to read from stdin: %v", err))
		return "", "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	defer stopLoader(output.LevelPrimary, output.IconAccept, "Read data from stdin")
	return data, "stdin", nil
}

// readFromFile reads the content of a file with security limits and validation
func readFromFile(filename string) (string, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	const maxSize = 10 * 1024 * 1024 // 10mb
	limitedReader := io.LimitReader(file, maxSize)

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// readFromStdin reads from standard input with security limits
func readFromStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("stdin is a terminal; please pipe data or provide a filename")
	}

	const maxSize = 10 * 1024 * 1024 // 10mb
	limitedReader := io.LimitReader(os.Stdin, maxSize)

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func init() {
	rootCmd.AddCommand(captureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// captureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// captureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
