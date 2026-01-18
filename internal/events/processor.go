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

	var active []*Event
	for id, e := range eventMap {
		if deleted[id] {
			continue
		}
		e.Hash = ""
		e.PrevHash = ""
		active = append(active, e)
	}

	// *tamper-proofing artefacts are removed due to the fact we are actively modifying
	// and deleting parts of the "chain"; exported data is not the source of truth
	var auditEvents []*Event
	for _, e := range rawEvents {
		if e.Type == "modify" || e.Type == "delete" {
			e.Hash = ""
			e.PrevHash = ""
			auditEvents = append(auditEvents, e)
		}
	}

	return &ProcessedEvents{
		Events: active,
		Audit:  auditEvents,
	}, nil
}

// FilterEvents filters events based on provided tags and types.
func FilterEvents(events []*Event, tags, types []string) []*Event {
	if len(tags) == 0 && len(types) == 0 {
		return events
	}

	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[t] = true
	}

	typeSet := make(map[string]bool)
	for _, t := range types {
		typeSet[t] = true
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
