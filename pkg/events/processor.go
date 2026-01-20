// Packages are shared with the closed-source praetor server
/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package events

// ProcessEvents takes a slice of raw events and sequentially applies modify
// and delete operations, returning ProcessedEvents with active events and
// an audit trail. The returned events have their Hash and PrevHash fields cleared.
func ProcessEvents(rawEvents []*Event) *ProcessedEvents {
	deleted := make(map[int]bool)
	eventMap := make(map[int]*Event)

	for _, e := range rawEvents {
		switch e.Type {
		case TypeDelete:
			deleted[e.RefId] = true
		case TypeModify:
			if target, ok := eventMap[e.RefId]; ok {
				target.Content = e.Content
			}
		default:
			eventMap[e.Id] = e
		}
	}

	var active []*Event
	for id, e := range eventMap {
		if !deleted[id] {
			e.Hash, e.PrevHash = "", ""
			active = append(active, e)
		}
	}

	var audit []*Event
	for _, e := range rawEvents {
		if e.Type == TypeModify || e.Type == TypeDelete {
			e.Hash, e.PrevHash = "", ""
			audit = append(audit, e)
		}
	}

	return &ProcessedEvents{Events: active, Audit: audit}
}

// FilterEvents filters events based on provided tags and types.
// If both tags and types are empty, all events are returned.
func FilterEvents(events []*Event, tags, types []string) []*Event {
	if len(tags) == 0 && len(types) == 0 {
		return events
	}

	typeSet := make(map[string]bool, len(types))
	for _, t := range types {
		typeSet[t] = true
	}
	tagSet := make(map[string]bool, len(tags))
	for _, t := range tags {
		tagSet[t] = true
	}

	var result []*Event
	for _, e := range events {
		if len(types) > 0 && !typeSet[e.Type] {
			continue
		}
		if len(tags) > 0 {
			hasTag := false
			for _, t := range e.Tags {
				if tagSet[t] {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		result = append(result, e)
	}
	return result
}
