package gocchan

import (
	"fmt"
	"reflect"
	"testing"
)

type TestFeature struct {
	name             string
	active           bool
	calledBy         []string
	activeIfCalledBy []string
}

func (f *TestFeature) ActiveIf(context interface{}, options ...interface{}) bool {
	f.activeIfCalledBy = append(f.activeIfCalledBy, fmt.Sprintf("%v:%v", context, options))
	return f.active
}

func (f *TestFeature) Func1(context interface{}) {
	f.calledBy = append(f.calledBy, fmt.Sprintf("Func1:%v", context))
}

func (f *TestFeature) Func2(context interface{}) {
	f.calledBy = append(f.calledBy, fmt.Sprintf("Func2:%v", context))
}

func (f *TestFeature) FuncPanic(context interface{}) {
	panic("expected panic")
}

func Test_AddFeature(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't occurred")
			}
		}()
		AddFeature("test", nil)
	}()

	func() {
		defer func() {
			featureStatus = make(map[string]*status)
		}()
		name := "test"
		if _, exists := featureStatus[name]; exists {
			t.Fatalf("Feature %v has already been added", name)
		}
		feature := &TestFeature{"test1", true, nil, nil}
		AddFeature("test", feature)
		actual := featureStatus[name]
		expected := &status{feature: feature, fault: false}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()
}

func Test_Invoke(t *testing.T) {
	init := func(name string, active bool) *TestFeature {
		feature := &TestFeature{name, active, nil, nil}
		featureStatus["testfeature"] = &status{
			feature: feature,
			fault:   false,
		}
		return feature
	}

	func() {
		init("test1", true)
		called := false
		Invoke(nil, "unknown", "Func1", func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called with unknown feature name")
		}
	}()

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("panic occurred by defaultFunc of nil")
			}
		}()
		init("test2", true)
		Invoke(nil, "unknown", "Func1", nil)
	}()

	func() {
		init("test3", true)
		called := false
		Invoke(nil, "testfeature", "unknown", func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called with unknown function name")
		}
	}()

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("panic occurred by defaultFunc of nil")
			}
		}()
		init("test4", true)
		Invoke(nil, "testfeature", "unknown", nil)
	}()

	func() {
		feature := init("test5", true)
		if len(feature.activeIfCalledBy) != 0 {
			t.Errorf("ActiveIf has already been called")
		}
		Invoke("ctx", "testfeature", "Func1", nil, "opt1", "opt2")
		actual := feature.activeIfCalledBy
		expected := []string{"ctx:[opt1 opt2]"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		init("test6", false)
		for _, fname := range []string{"Func1", "Func2"} {
			called := false
			Invoke(nil, "testfeature", fname, func() {
				called = true
			})
			if !called {
				t.Errorf("defaultFunc hasn't been called by not active: %v", fname)
			}
		}
	}()

	feature := init("test7", true)
	if len(feature.calledBy) != 0 {
		t.Fatalf("feature function has already been called")
	}
	Invoke("testctx1", "testfeature", "Func1", func() {
		t.Errorf("defaultFunc has been called")
	})
	actual := feature.calledBy
	expected := []string{"Func1:testctx1"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	Invoke("testctx2", "testfeature", "Func1", func() {
		t.Errorf("defaultFunc has been called")
	})
	actual = feature.calledBy
	expected = []string{"Func1:testctx1", "Func1:testctx2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	Invoke("testctx3", "testfeature", "Func2", func() {
		t.Errorf("defaultFunc has been called")
	})
	actual = feature.calledBy
	expected = []string{"Func1:testctx1", "Func1:testctx2", "Func2:testctx3"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	called := false
	Invoke("testctx4", "testfeature", "FuncPanic", func() {
		called = true
	})
	if !called {
		t.Errorf("defaultFunc hasn't been called by panic")
	}
	actual = feature.calledBy
	expected = []string{"Func1:testctx1", "Func1:testctx2", "Func2:testctx3"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	for _, fname := range []string{"Func1", "Func2"} {
		called = false
		Invoke("testctx5", "testfeature", fname, func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called by after panic: %v", fname)
		}
		actual = feature.calledBy
		expected = []string{"Func1:testctx1", "Func1:testctx2", "Func2:testctx3"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}
}
