/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package output

import "fmt"

// Basic log utilities

func (of *OutputFormatter) Logf(level OutputLevel, iconType IconType, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	of.Log(level, iconType, message)
}
func (of *OutputFormatter) LogDefault(message string) {
	of.Log(LevelDefault, IconNone, message)
}
func (of *OutputFormatter) LogDefaultf(format string, args ...any) {
	of.Logf(LevelDefault, IconNone, format, args...)
}
func (of *OutputFormatter) LogPrimary(message string) {
	of.Log(LevelPrimary, IconNone, message)
}
func (of *OutputFormatter) LogPrimaryf(format string, args ...any) {
	of.Logf(LevelPrimary, IconNone, format, args...)
}
func (of *OutputFormatter) LogSuccess(message string) {
	of.Log(LevelPrimary, IconAccept, message)
}
func (of *OutputFormatter) LogSuccessf(format string, args ...any) {
	of.Logf(LevelPrimary, IconAccept, format, args...)
}
func (of *OutputFormatter) LogWarning(message string) {
	of.Log(LevelWarning, IconWarning, message)
}
func (of *OutputFormatter) LogWarningf(format string, args ...any) {
	of.Logf(LevelWarning, IconWarning, format, args...)
}
func (of *OutputFormatter) LogError(message string) {
	of.Log(LevelError, IconReject, message)
}
func (of *OutputFormatter) LogErrorf(format string, args ...any) {
	of.Logf(LevelError, IconReject, format, args...)
}
func (of *OutputFormatter) LogTask(message string) {
	of.Log(LevelPrimary, IconArrow, message)
}
func (of *OutputFormatter) LogTaskf(format string, args ...any) {
	of.Logf(LevelPrimary, IconArrow, format, args...)
}

// More advanced utilities

// LogStep outputs a sub-step message with a dash icon at the current indent level
func (of *OutputFormatter) LogStep(message string) {
	level := LevelPrimary
	if of.indentLevel > 0 {
		level = LevelMuted
	}
	of.Log(level, IconDash, message)
}

// LogStepf outputs a formatted sub-step message
func (of *OutputFormatter) LogStepf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	of.LogStep(message)
}

// Indent increases the indentation level for subsequent output
func (of *OutputFormatter) Indent() {
	of.mu.Lock()
	defer of.mu.Unlock()
	of.indentLevel++
}

// Dedent decreases the indentation level
func (of *OutputFormatter) Dedent() {
	of.mu.Lock()
	defer of.mu.Unlock()
	if of.indentLevel > 0 {
		of.indentLevel--
	}
}

// SetIndentLevel sets the indentation level to a specific value
func (of *OutputFormatter) SetIndentLevel(level int) {
	of.mu.Lock()
	defer of.mu.Unlock()
	if level >= 0 {
		of.indentLevel = level
	}
}

// GetIndentLevel returns the current indentation level
func (of *OutputFormatter) GetIndentLevel() int {
	of.mu.RLock()
	defer of.mu.RUnlock()
	return of.indentLevel
}

// WithIndent temporarily increases indentation for a function, then restores it
func (of *OutputFormatter) WithIndent(fn func()) {
	of.Indent()
	defer of.Dedent()
	fn()
}

// WithIndentN temporarily increases indentation by n levels, then restores it
func (of *OutputFormatter) WithIndentN(n int, fn func()) {
	for i := 0; i < n; i++ {
		of.Indent()
	}
	defer func() {
		for i := 0; i < n; i++ {
			of.Dedent()
		}
	}()
	fn()
}

// Global convenience functions that use the default formatter

func Log(level OutputLevel, iconType IconType, message string) {
	GetOutput().Log(level, iconType, message)
}
func Logf(level OutputLevel, iconType IconType, format string, args ...any) {
	GetOutput().Logf(level, iconType, format, args...)
}
func LogDefault(message string) {
	GetOutput().LogDefault(message)
}
func LogDefaultf(format string, args ...any) {
	GetOutput().LogDefaultf(format, args...)
}
func LogPrimary(message string) {
	GetOutput().LogPrimary(message)
}
func LogPrimaryf(format string, args ...any) {
	GetOutput().LogPrimaryf(format, args...)
}
func LogSuccess(message string) {
	GetOutput().LogSuccess(message)
}
func LogSuccessf(format string, args ...any) {
	GetOutput().LogSuccessf(format, args...)
}
func LogWarning(message string) {
	GetOutput().LogWarning(message)
}
func LogWarningf(format string, args ...any) {
	GetOutput().LogWarningf(format, args...)
}
func LogError(message string) {
	GetOutput().LogError(message)
}
func LogErrorf(format string, args ...any) {
	GetOutput().LogErrorf(format, args...)
}
func LogTask(message string) {
	GetOutput().LogTask(message)
}
func LogTaskf(format string, args ...any) {
	GetOutput().LogTaskf(format, args...)
}
func LogStep(message string) {
	GetOutput().LogStep(message)
}
func LogStepf(format string, args ...any) {
	GetOutput().LogStepf(format, args...)
}

// loader utilities

func StartLoader(id string, message string) func(OutputLevel, IconType, string) {
	return GetOutput().StartLoader(id, message)
}

// Indentation utilities

func Indent() {
	GetOutput().Indent()
}
func Dedent() {
	GetOutput().Dedent()
}
func SetIndentLevel(level int) {
	GetOutput().SetIndentLevel(level)
}
func GetIndentLevel() int {
	return GetOutput().GetIndentLevel()
}
func WithIndent(fn func()) {
	GetOutput().WithIndent(fn)
}
func WithIndentN(n int, fn func()) {
	GetOutput().WithIndentN(n, fn)
}

// direct output wrappers

func Println(args ...any) {
	fmt.Fprintln(GetOutput().writer, args...)
}
func Printlnf(format string, args ...any) {
	fmt.Fprintf(GetOutput().writer, format+"\n", args...)
}
