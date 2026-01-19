package formats

import (
	"fmt"
	"strings"
	"time"
)

func formatFields(fields map[string]any) string {
	parts := make([]string, 0, len(fields))
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, " ")
}

func formatTime(ts string) string {
	parsed, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return ts
	}
	return parsed.Format("2006-01-02 15:04")
}
