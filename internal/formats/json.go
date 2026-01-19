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

func renderEventReportJSON(eventList, audit []*events.Event, opts RenderOptions) (string, error) {
	out := map[string]any{
		"events": eventList,
	}

	if opts.IncludeAudit && len(audit) > 0 {
		out["audit"] = audit
	}
	if opts.IncludeMetadata && len(opts.Metadata) > 0 {
		out["metadata"] = opts.Metadata
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b) + "\n", nil
}
