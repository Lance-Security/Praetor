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

func Render(format Format, processed *events.ProcessedEvents) (string, error) {
	switch format {
	case FormatJSON:
		return renderJSON(processed)
	default:
		return renderTerminal(processed)
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
