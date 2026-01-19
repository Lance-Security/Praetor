package formats

import (
	"strings"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/utils"
)

func init() {
	RegisterMessageRenderer(FormatTerminal, renderMessagesTerminal)
}

func renderTerminal(processed *events.ProcessedEvents) (string, error) {
	var b strings.Builder

	for _, e := range processed.Events {
		ts := formatTime(e.Timestamp)
		typeLabel := strings.ToUpper(e.Type)

		b.WriteString(utils.Mutedf("%s ", ts))
		b.WriteString(events.StyleType(typeLabel))
		b.WriteString(utils.Muted(" "))
		b.WriteString(e.Content)

		// Show event ID for notes
		if e.Type == "note" && e.Id > 0 {
			b.WriteString(utils.Mutedf(" (id: %d)", e.Id))
		}

		if len(e.Tags) > 0 {
			b.WriteString(utils.Mutedf(" [%s]", strings.Join(e.Tags, ", ")))
		}

		b.WriteString("\n")
	}

	return b.String(), nil
}

func renderMessagesTerminal(messages []Message, opts Options) (string, error) {
	var b strings.Builder
	for _, m := range messages {
		if m.Event != nil {
			b.WriteString(renderMessageEventLine(*m.Event))
			b.WriteString("\n")
			continue
		}

		line := strings.TrimSpace(m.Text)
		if line == "" && len(m.Fields) == 0 {
			continue
		}

		if opts.UseTimestamp {
			ts := m.Timestamp
			if ts == "" {
				ts = time.Now().UTC().Format(time.RFC3339Nano)
			}
			b.WriteString(utils.Mutedf("[%s] ", formatTime(ts)))
		}

		b.WriteString(styleByLevel(m.Level, line))

		if len(m.Fields) > 0 {
			b.WriteString(utils.Muted(" "))
			b.WriteString(utils.Muted(formatFields(m.Fields)))
		}

		b.WriteString("\n")
	}
	return b.String(), nil
}

func renderMessageEventLine(ev events.Event) string {
	timestamp := formatTime(ev.Timestamp)
	typeLabel := strings.ToUpper(ev.Type)
	var b strings.Builder
	b.WriteString(utils.Mutedf("[%s] ", timestamp))
	b.WriteString(events.StyleType(typeLabel))
	b.WriteString(utils.Mutedf(" by %s in %s", ev.User, events.ShortenCwd(ev.Cwd, 30)))
	if ev.Id != 0 {
		b.WriteString(utils.Mutedf(" id %d | ", ev.Id))
	} else {
		b.WriteString(utils.Muted(" | "))
	}
	b.WriteString(utils.Default(ev.Content))
	if len(ev.Tags) > 0 {
		b.WriteString(utils.Mutedf(" [%s]", strings.Join(ev.Tags, ", ")))
	}
	return b.String()
}

func styleByLevel(level Level, s string) string {
	switch level {
	case LevelSuccess:
		return utils.Accept(s)
	case LevelWarn:
		return utils.Warning(s)
	case LevelError:
		return utils.Error(s)
	default:
		return utils.Default(s)
	}
}
