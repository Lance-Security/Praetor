/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

// internal/events/formatter.go
// Reads events and formats them for common displays
// including terminal output and markdown

package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/utils"
)

// ShowEventTerminal formats a single event for terminal output
//
// [{time}] {type} by {user} in {cwd} | {content}
// time in muted colour, type in caps bold, "by {user} in cwd" in dim, content normal
// time formatted from RFC3339Nano to "2006-01-02 15:04"
// TODO: get configuration for format & colours
func ShowEventTerminal(e Event) string {
	timestamp := e.Timestamp
	if ts, err := time.Parse(time.RFC3339Nano, e.Timestamp); err == nil {
		timestamp = ts.Format("2006-01-02 15:04")
	}

	typeLabel := strings.ToUpper(e.Type)
	var typeStyled string
	switch strings.ToLower(e.Type) {
	case "note":
		typeStyled = utils.Primary(typeLabel)
	case "command":
		typeStyled = utils.Warning(typeLabel)
	case "result":
		typeStyled = utils.Accept(typeLabel)
	case "error":
		typeStyled = utils.Error(typeLabel)
	default:
		typeStyled = utils.Primary(typeLabel)
	}

	var b strings.Builder
	b.WriteString(utils.Mutedf("[%s] ", timestamp))
	b.WriteString(fmt.Sprintf("%s ", typeStyled))
	//b.WriteString(utils.Mutedf("by %s in %s | ", e.User, e.Cwd))
	b.WriteString(utils.Mutedf("by "))
	b.WriteString(fmt.Sprintf("%s ", e.User))
	b.WriteString(utils.Mutedf("in "))
	b.WriteString(fmt.Sprintf("%s ", ShortenCwd(e.Cwd, 30)))
	b.WriteString(utils.Mutedf("id %v | ", e.Id))
	b.WriteString(e.Content) // normal content
	return b.String()
}

// ShowEventsTerminal formats multiple events for terminal output
func ShowEventsTerminal(events []*Event) string {
	var b strings.Builder
	for _, e := range events {
		b.WriteString(ShowEventTerminal(*e))
		b.WriteString("\n")
	}
	return b.String()
}

// ShortenCwd shortens a directory path to fit within maxLength
func ShortenCwd(cwd string, maxLength int) string {
	if len(cwd) <= maxLength {
		// /root/middle/final-path
		return cwd
	}
	parts := strings.Split(cwd, "/")
	if len(parts) < 3 {
		// /root/final-path
		return cwd
	}
	start := parts[0]
	if start == "" && len(parts) > 1 {
		start = "/" + parts[1]
		parts = parts[1:]
	}
	end := parts[len(parts)-1]
	// /root/.../final-path
	shortened := start + "/.../" + end

	if len(shortened) > maxLength {
		// .../final-path
		return "..." + end
	}
	return shortened
}
