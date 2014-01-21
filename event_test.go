package gocchan

import (
	"reflect"
	"testing"
)

func Test_EventType_String(t *testing.T) {
	for expected, ev := range map[string]EventType{
		"EventFeatureHasNotBeenAdded":                EventFeatureHasNotBeenAdded,
		"EventFeatureMethodMissing":                  EventFeatureMethodMissing,
		"EventFeatureWasFault":                       EventFeatureWasFault,
		"EventFeatureMethodInvalidNumberOfArguments": EventFeatureMethodInvalidNumberOfArguments,
		"EventFeatureMethodSignatureMismatch":        EventFeatureMethodSignatureMismatch,
		"unknown": -1,
	} {
		actual := ev.String()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}
}

func Test_NewEvent(t *testing.T) {
	actual := NewEvent(EventFeatureWasFault, "testerr1")
	expected := &Event{Type: EventFeatureWasFault, Err: "testerr1"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	actual = NewEvent(EventFeatureWasFault, "testerr2")
	expected = &Event{Type: EventFeatureWasFault, Err: "testerr2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	actual = NewEvent(EventFeatureMethodMissing, "testerr3")
	expected = &Event{Type: EventFeatureMethodMissing, Err: "testerr3"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}
}
