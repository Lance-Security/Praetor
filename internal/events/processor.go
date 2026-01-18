package events

type ProcessedEvents struct {
	Events []*Event
	Audit  []*Event
}

// PrepareEvents takes a raw `events.jsonl` and sequentially applies modify
// and delete operations, and returns an *[]Event with the tamper-proofing
// removed*.
func PrepareEvents(path string) (*ProcessedEvents, error) {
	rawEvents, err := GetAllEvents(path)
	if err != nil {
		return nil, err
	}

	deleted, eventMap := processRawEvents(rawEvents)
	active := buildActiveEvents(eventMap, deleted)
	audit := buildAuditTrail(rawEvents)

	return &ProcessedEvents{
		Events: active,
		Audit:  audit,
	}, nil
}

func processRawEvents(rawEvents []*Event) (map[int]bool, map[int]*Event) {
	deleted := make(map[int]bool)
	eventMap := make(map[int]*Event)

	for _, e := range rawEvents {
		switch e.Type {
		case "delete":
			deleted[e.RefId] = true
		case "modify":
			if target, ok := eventMap[e.RefId]; ok {
				target.Content = e.Content
			}
		default:
			eventMap[e.Id] = e
		}
	}

	return deleted, eventMap
}

func buildActiveEvents(eventMap map[int]*Event, deleted map[int]bool) []*Event {
	var active []*Event
	for id, e := range eventMap {
		if deleted[id] {
			continue
		}
		e.Hash = ""
		e.PrevHash = ""
		active = append(active, e)
	}
	return active
}

func buildAuditTrail(rawEvents []*Event) []*Event {
	var auditEvents []*Event
	for _, e := range rawEvents {
		if e.Type == "modify" || e.Type == "delete" {
			e.Hash = ""
			e.PrevHash = ""
			auditEvents = append(auditEvents, e)
		}
	}
	return auditEvents
}

// FilterEvents filters events based on provided tags and types.
func FilterEvents(events []*Event, tags, types []string) []*Event {
	if len(tags) == 0 && len(types) == 0 {
		return events
	}

	typeSet := makeSet(types)
	tagSet := makeSet(tags)

	var result []*Event
	for _, e := range events {
		if matchesFilters(e, typeSet, tagSet, len(types) > 0, len(tags) > 0) {
			result = append(result, e)
		}
	}

	return result
}

func makeSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

func matchesFilters(e *Event, typeSet, tagSet map[string]bool, hasTypeFilter, hasTagFilter bool) bool {
	if hasTypeFilter && !typeSet[e.Type] {
		return false
	}
	if hasTagFilter && !hasEventTag(e, tagSet) {
		return false
	}
	return true
}

func hasEventTag(e *Event, tagSet map[string]bool) bool {
	for _, t := range e.Tags {
		if tagSet[t] {
			return true
		}
	}
	return false
}
