package gocchan

// EventType represents a type of event.
type EventType int

const (
	EventFeatureHasNotBeenAdded EventType = iota + 1
	EventFeatureMethodMissing
	EventFeatureWasFault
)

// String returns a name of event type.
func (typ EventType) String() string {
	switch typ {
	case EventFeatureHasNotBeenAdded:
		return "EventFeatureHasNotBeenAdded"
	case EventFeatureMethodMissing:
		return "EventFeatureMethodMissing"
	case EventFeatureWasFault:
		return "EventFeatureWasFault"
	}
	return "unknown"
}

// Event represents a event.
type Event struct {
	// type of event.
	Type EventType

	// additional information of event.
	Err interface{}
}

// NewEvent returns a new event.
func NewEvent(typ EventType, err interface{}) *Event {
	return &Event{
		Type: typ,
		Err:  err,
	}
}
