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

func renderTerminal(eventList, audit []*events.Event, opts RenderOptions) (string, error) {
	var b strings.Builder

	for _, e := range eventList {
		b.WriteString(renderTerminalEventLine(e))
		b.WriteString("\n")
	}

	if opts.IncludeAudit && len(audit) > 0 {
		b.WriteString("\n")
		b.WriteString(utils.Muted("Audit:\n"))
		for _, e := range audit {
			b.WriteString(renderTerminalAuditLine(e))
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}

func renderTerminalEventLine(e *events.Event) string {
	ts := formatTime(e.Timestamp)
	typeLabel := strings.ToUpper(e.Type)

	var b strings.Builder
	b.WriteString(utils.Mutedf("%s ", ts))
	b.WriteString(events.StyleType(typeLabel))
	b.WriteString(utils.Muted(" "))
	b.WriteString(e.Content)

	if len(e.Tags) > 0 {
		b.WriteString(utils.Mutedf(" [%s]", strings.Join(e.Tags, ", ")))
	}

	return b.String()
}

func renderTerminalAuditLine(e *events.Event) string {
	ts := formatTime(e.Timestamp)

	var b strings.Builder
	b.WriteString(utils.Mutedf("%s ", ts))
	b.WriteString(utils.Muted(strings.ToUpper(e.Type)))
	b.WriteString(utils.Muted(": "))
	b.WriteString(e.Content)

	if e.RefId > 0 {
		b.WriteString(utils.Mutedf(" (event_id=%d)", e.RefId))
	}

	return b.String()
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
