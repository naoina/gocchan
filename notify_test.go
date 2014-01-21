package gocchan

import (
	"reflect"
	"testing"
)

type testListener struct {
	name   string
	listen *Event
}

func (listener *testListener) Listen(event *Event) {
	listener.listen = event
}

func Test_Notifier_NotifyAll(t *testing.T) {
	listeners := []Listener{
		&testListener{name: "test1"},
		&testListener{name: "test2"},
	}
	runTest := func(event *Event) {
		notifier := &Notifier{listeners: listeners}
		notifier.NotifyAll(event)
		notifier.Wait()
		for _, listener := range listeners {
			l := listener.(*testListener)
			actual := l.listen
			expected := event
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("listener %v expect %q, but %q", l.name, expected, actual)
			}
		}
	}
	runTest(&Event{Err: "test1event"})
	runTest(&Event{Err: "test2event"})
}

func Test_AddEventListener(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't occurred by listener is nil")
			}
		}()
		AddEventListener(nil)
	}()

	if len(notifier.listeners) > 0 {
		t.Fatalf("listeners has already been added")
	}
	listener := &testListener{name: "listener1"}
	AddEventListener(listener)
	actual := notifier.listeners
	expected := []Listener{listener}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	listener2 := &testListener{name: "listener2"}
	AddEventListener(listener2)
	actual = notifier.listeners
	expected = []Listener{listener, listener2}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}
}
