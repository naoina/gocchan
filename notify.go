package gocchan

import "sync"

// Notifier represents a notifier of any event.
type Notifier struct {
	listeners []Listener
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// WaitNotify blocks until the all notifications of global notifier is finished.
func WaitNotify() {
	notifier.Wait()
}

// NotifyAll notify event to all listeners.
func (n *Notifier) NotifyAll(event *Event) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, listener := range n.listeners {
		n.wg.Add(1)
		go func(listener Listener) {
			defer n.wg.Done()
			listener.Listen(event)
		}(listener)
	}
}

// Wait blocks until the all notifications is finished.
func (n *Notifier) Wait() {
	n.wg.Wait()
}

// Listener is an interface of listener of event.
type Listener interface {
	Listen(event *Event)
}

// global notifier.
var notifier = &Notifier{}

// AddEventListener adds a listener of event.
// If listener is nil, it panic.
func AddEventListener(listener Listener) {
	if listener == nil {
		panic("Add Listener is nil")
	}
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	notifier.listeners = append(notifier.listeners, listener)
}
