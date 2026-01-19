package formats

import (
	"strings"

	"github.com/lachlanharrisdev/praetor/internal/events"
)

type Format int

const (
	FormatTerminal Format = iota
	FormatJSON
)

type RenderOptions struct {
	IncludeMetadata bool
	IncludeAudit    bool
	Tags            []string
	Types           []string
	MaxEvents       int
	Metadata        map[string]string
}

func Render(format Format, processed *events.ProcessedEvents, opts RenderOptions) (string, error) {
	eventList := processed.Events

	if len(opts.Tags) > 0 || len(opts.Types) > 0 {
		eventList = events.FilterEvents(eventList, opts.Tags, opts.Types)
	}

	if opts.MaxEvents > 0 && len(eventList) > opts.MaxEvents {
		eventList = eventList[len(eventList)-opts.MaxEvents:]
	}

	switch format {
	case FormatJSON:
		return renderEventReportJSON(eventList, processed.Audit, opts)
	default:
		return renderTerminal(eventList, processed.Audit, opts)
	}
}

func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json", "j", "js", "jsn":
		return FormatJSON
	case "terminal", "term", "t":
		return FormatTerminal
	default:
		return FormatTerminal
	}
}

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatTerminal:
		return "terminal"
	default:
		return "unknown"
	}
}
