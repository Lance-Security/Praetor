package formats

import (
	"encoding/json"
	"strings"

	"github.com/lachlanharrisdev/praetor/internal/events"
)

func init() {
	RegisterMessageRenderer(FormatJSON, renderMessagesJSON)
}

func RenderJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func renderJSON(processed *events.ProcessedEvents) (string, error) {
	out := map[string]any{
		"events": processed.Events,
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b) + "\n", nil
}

func renderMessagesJSON(msgs []Message, _ Options) (string, error) {
	var b strings.Builder
	for _, m := range msgs {
		line, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String(), nil
}
