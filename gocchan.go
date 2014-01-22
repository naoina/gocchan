package gocchan

import (
	"errors"
	"fmt"
	"reflect"
)

var featureStatus = make(map[string]*status)

type status struct {
	feature Feature

	// Whether the feature was fault.
	fault bool
}

var (
	ErrInvokeDefault = errors.New("invoke defualt")
)

// Feature is an interface of feature.
type Feature interface {
	// ActiveIf returns whether the feature is active.
	ActiveIf(context interface{}, options ...interface{}) bool
}

// ActiveIf returns true if ActiveIf of Feature returns true, otherwise returns false.
func ActiveIf(featureName string, context interface{}, options ...interface{}) bool {
	status := featureStatus[featureName]
	if status == nil || status.fault {
		return false
	}
	return status.feature.ActiveIf(context, options...)
}

// AddFeature adds feature with name.
// If feature is nil, it panic.
func AddFeature(name string, feature Feature) {
	if feature == nil {
		panic("Add Feature is nil")
	}
	featureStatus[name] = &status{
		feature: feature,
		fault:   false,
	}
}

// Invoke invokes function of added feature.
// context and options are passed to ActiveIf() method of the Feature associated with featureName.
// Will invoke the method named funcName if defined in Feature associated with featureName.
// When featureName hasn't been added, funcName hasn't been defined, or any errors occurred,
// will invoke the defaultFunc with given context if defaultFunc isn't nil.
// Also if any errors occurred at least once, next invoking will always invoke the defaultFunc.
func Invoke(context interface{}, featureName, funcName string, defaultFunc func(), options ...interface{}) {
	status := featureStatus[featureName]
	defer func() {
		if err := recover(); err != nil {
			if err != ErrInvokeDefault {
				status.fault = true
				notifier.NotifyAll(NewEvent(EventFeatureWasFault, err))
			}
			if defaultFunc != nil {
				defaultFunc()
			}
		}
	}()
	if status == nil {
		err := fmt.Errorf("feature has not been added: `%s`", featureName)
		event := NewEvent(EventFeatureHasNotBeenAdded, err)
		notifier.NotifyAll(event)
		panic(ErrInvokeDefault)
	}
	if status.fault {
		panic(ErrInvokeDefault)
	}
	f := reflect.ValueOf(status.feature).MethodByName(funcName)
	if !f.IsValid() {
		err := fmt.Errorf("method is not found: `%s` in feature `%s`", funcName, featureName)
		event := NewEvent(EventFeatureMethodMissing, err)
		notifier.NotifyAll(event)
		panic(ErrInvokeDefault)
	}
	ftype := f.Type()
	if ftype.NumIn() != 1 {
		err := fmt.Errorf("number of arguments must be one: method `%s` in feature `%s`", funcName, featureName)
		event := NewEvent(EventFeatureMethodInvalidNumberOfArguments, err)
		notifier.NotifyAll(event)
		panic(ErrInvokeDefault)
	}
	cvalue := reflect.ValueOf(context)
	if !cvalue.IsValid() {
		cvalue = reflect.ValueOf(&context).Elem()
	}
	if !cvalue.Type().AssignableTo(ftype.In(0)) {
		err := fmt.Errorf("method signature mismatch: context is a type `%T`, but type `%s` is an argument type of the method `%s` in feature `%s`", context, ftype.In(0), funcName, featureName)
		event := NewEvent(EventFeatureMethodSignatureMismatch, err)
		notifier.NotifyAll(event)
		panic(ErrInvokeDefault)
	}
	if !status.feature.ActiveIf(context, options...) {
		panic(ErrInvokeDefault)
	}
	f.Call([]reflect.Value{cvalue})
}

// IsActive returns true if feature is active, otherwise returns false.
func IsActive(featureName string) bool {
	status := featureStatus[featureName]
	if status == nil {
		return false
	}
	return !status.fault
}
