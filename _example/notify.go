package main

import (
	"fmt"

	"github.com/naoina/gocchan"
)

type NotifyFeature struct{}

func (f *NotifyFeature) ActiveIf(context interface{}, options ...interface{}) bool {
	return true
}

func (f *NotifyFeature) Panic(context interface{}) {
	panic("expected panic")
}

type NotifyListen struct{}

func (n *NotifyListen) Listen(event *gocchan.Event) {
	switch event.Type {
	case gocchan.EventFeatureHasNotBeenAdded:
		fallthrough
	case gocchan.EventFeatureMethodMissing:
		fallthrough
	case gocchan.EventFeatureWasFault:
		fmt.Println(event.Type)
		fmt.Println(event.Err)
	default:
		fmt.Println("unknown event: %v", event.Type)
	}
}

func main() {
	gocchan.AddFeature("notify", &NotifyFeature{})
	gocchan.AddEventListener(&NotifyListen{})
	gocchan.Invoke("", "notify", "Panic", nil)
	gocchan.WaitNotify()
}
