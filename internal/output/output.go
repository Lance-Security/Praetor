/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package output

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/utils"
)

type IconType int

const (
	IconNone IconType = iota
	IconDash
	IconArrow
	IconAccept
	IconReject
	IconWarning
	IconLoader
)

type OutputLevel int

const (
	// standard output with default colour (white)
	LevelDefault OutputLevel = iota
	// primary/top-level output in primary colour (blue)
	LevelPrimary
	// sub-level output in muted / dim colour (gray)
	LevelMuted
	// warning output in warning colour (yellow)
	LevelWarning
	// error output in error colour (red)
	LevelError
)

// for loader animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var iconMap = map[IconType]string{
	IconNone:    "",
	IconDash:    "–",
	IconArrow:   "→",
	IconAccept:  "✓",
	IconReject:  "✗",
	IconWarning: "⚠",
}

// loaderState tracks the state of an active loader
type loaderState struct {
	frameIndex int
	stopChan   chan struct{}
	stopped    bool
}

// OutputFormatter manages standardized console output with Docker-like formatting
type OutputFormatter struct {
	mu             sync.RWMutex
	writer         io.Writer
	indentLevel    int
	activeLoaders  map[string]*loaderState
	loaderTicker   *time.Ticker
	loaderStopChan chan struct{}
	loaderWg       sync.WaitGroup
	initOnce       sync.Once
}

// NewOutputFormatter creates a new OutputFormatter instance
func NewOutputFormatter(writer io.Writer) *OutputFormatter {
	if writer == nil {
		writer = os.Stdout
	}
	return &OutputFormatter{
		writer:        writer,
		indentLevel:   0,
		activeLoaders: make(map[string]*loaderState),
	}
}

// startLoaderTicker initializes the global loader animation ticker
func (of *OutputFormatter) startLoaderTicker() {
	of.initOnce.Do(func() {
		of.loaderTicker = time.NewTicker(100 * time.Millisecond)
		of.loaderStopChan = make(chan struct{})

		of.loaderWg.Add(1)
		go func() {
			defer of.loaderWg.Done()
			for {
				select {
				case <-of.loaderTicker.C:
					of.mu.Lock()
					// increment all active loader frames
					for _, state := range of.activeLoaders {
						state.frameIndex = (state.frameIndex + 1) % len(spinnerFrames)
					}
					of.mu.Unlock()
				case <-of.loaderStopChan:
					return
				}
			}
		}()
	})
}

// getIndentation returns the indentation string based on the current level
/*
func (of *OutputFormatter) getIndentation() string {
	const indentSpace = "  " // 2 spaces per level
	result := ""
	for i := 0; i < of.indentLevel; i++ {
		result += indentSpace
	}
	return result
}*/

// getIcon returns the icon string for the given type
func (of *OutputFormatter) getIcon(iconType IconType, frameIndex int) string {
	if iconType == IconLoader {
		if frameIndex >= 0 && frameIndex < len(spinnerFrames) {
			return spinnerFrames[frameIndex]
		}
		return spinnerFrames[0]
	}
	if icon, exists := iconMap[iconType]; exists {
		return icon
	}
	return ""
}

// formatOutput constructs the formatted output line (must be called with indentLevel already known)
func (of *OutputFormatter) formatOutput(level OutputLevel, iconType IconType, message string, frameIndex int, indentLevel int) string {
	const indentSpace = "  "
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent += indentSpace
	}

	icon := of.getIcon(iconType, frameIndex)

	var colorFunc func(...any) string
	switch level {
	case LevelPrimary:
		colorFunc = utils.Primary
	case LevelWarning:
		colorFunc = utils.Warning
	case LevelError:
		colorFunc = utils.Error
	case LevelMuted:
		colorFunc = utils.Muted
	case LevelDefault:
		fallthrough
	default:
		colorFunc = utils.Default
	}

	var output string
	if icon != "" {
		output = fmt.Sprintf("%s%s %s", indent, icon, message)
	} else {
		output = fmt.Sprintf("%s%s", indent, message)
	}

	return colorFunc(output)
}

// clearLine sends ANSI escape sequence to clear the current line
func clearLine() {
	w := io.Writer(os.Stdout)
	if defaultFormatter != nil && defaultFormatter.writer != nil {
		w = defaultFormatter.writer
	}
	fmt.Fprint(w, "\r\033[K")
}

// Log outputs a message with the specified level and icon
func (of *OutputFormatter) Log(level OutputLevel, iconType IconType, message string) {
	of.mu.RLock()
	indentLevel := of.indentLevel
	of.mu.RUnlock()

	output := of.formatOutput(level, iconType, message, 0, indentLevel)
	fmt.Fprintln(of.writer, output)
}

// StartLoader begins displaying a loading spinner with the given message
// Returns a function that should be called to stop the loader
func (of *OutputFormatter) StartLoader(id string, message string) func(OutputLevel, IconType, string) {
	of.startLoaderTicker()

	// Create new loader state outside the lock
	state := &loaderState{
		frameIndex: 0,
		stopChan:   make(chan struct{}),
	}

	// Register the loader
	of.mu.Lock()
	of.activeLoaders[id] = state
	of.mu.Unlock()

	// Start loader render goroutine
	// The global ticker in startLoaderTicker() handles frameIndex updates
	// This goroutine only renders the current frame
	go func() {
		for {
			select {
			case <-of.loaderTicker.C:
				of.mu.Lock()
				currentState, exists := of.activeLoaders[id]
				if !exists {
					of.mu.Unlock()
					return
				}
				frameIndex := currentState.frameIndex
				indentLevel := of.indentLevel
				of.mu.Unlock()

				clearLine()
				output := of.formatOutput(LevelDefault, IconLoader, message, frameIndex, indentLevel)
				fmt.Fprint(of.writer, "\r"+output)

			case <-state.stopChan:
				return
			}
		}
	}()

	// Return stop function
	return func(finalLevel OutputLevel, finalIcon IconType, finalMessage string) {
		of.mu.Lock()
		loaderState, exists := of.activeLoaders[id]
		if exists && !loaderState.stopped {
			loaderState.stopped = true
			close(loaderState.stopChan)
			delete(of.activeLoaders, id)
		}
		of.mu.Unlock()

		// Clear the line and print final message
		clearLine()
		of.Log(finalLevel, finalIcon, finalMessage)
	}
}

// Close cleans up resources used by the formatter
func (of *OutputFormatter) Close() {
	of.mu.Lock()
	defer of.mu.Unlock()

	if of.loaderTicker != nil {
		of.loaderTicker.Stop()
	}

	if of.loaderStopChan != nil {
		close(of.loaderStopChan)
	}

	of.loaderWg.Wait()
}

// Global formatter instance
var defaultFormatter *OutputFormatter
var defaultOnce sync.Once

// GetOutput returns the default output formatter, initializing if needed
func GetOutput() *OutputFormatter {
	defaultOnce.Do(func() {
		defaultFormatter = NewOutputFormatter(os.Stdout)
	})
	return defaultFormatter
}
