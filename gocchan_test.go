package gocchan

import (
	"fmt"
	"reflect"
	"testing"
)

type TestFeature struct {
	t                *testing.T
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

func (f *TestFeature) Func3(context string) {
	f.calledBy = append(f.calledBy, fmt.Sprintf("Func3:%v", context))
}

func (f *TestFeature) Func4(context interface{}, other interface{}) {
	f.t.Errorf("Func4 is never called")
}

func (f *TestFeature) Func5() {
	f.t.Errorf("Func5 is never called")
}

func (f *TestFeature) FuncPanic(context interface{}) {
	panic("expected panic")
}

func Test_ActiveIf(t *testing.T) {
	func() {
		defer func() {
			featureStatus = make(map[string]*status)
		}()
		if featureStatus["test"] != nil {
			t.Fatalf("feature has already been added")
		}
		var actual interface{} = ActiveIf("test", "ctx1", "opt1")
		var expected interface{} = false
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	for _, v := range []struct {
		fault, active, expected bool
		calledBy                []string
	}{
		{false, true, true, []string{"ctx1:[opt1]"}},
		{false, false, false, []string{"ctx1:[opt1]"}},
		{true, true, false, nil},
		{true, false, false, nil},
	} {
		func() {
			defer func() {
				featureStatus = make(map[string]*status)
			}()
			if featureStatus["test"] != nil {
				t.Fatalf("feature has already been added")
			}
			feature := &TestFeature{t, "test1", v.active, nil, nil}
			featureStatus["test"] = &status{
				feature: feature,
				fault:   v.fault,
			}
			var actual interface{} = ActiveIf("test", "ctx1", "opt1")
			var expected interface{} = v.expected
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("fault is %#v, active is %#v expect %q, but %q", v.fault, v.active, expected, actual)
			}
			actual = feature.activeIfCalledBy
			expected = v.calledBy
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("fault is %#v, active is %#v expect %q, but %q", v.fault, v.active, expected, actual)
			}
		}()
	}
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
		feature := &TestFeature{t, "test1", true, nil, nil}
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
		feature := &TestFeature{t, name, active, nil, nil}
		featureStatus["testfeature"] = &status{
			feature: feature,
			fault:   false,
		}
		return feature
	}

	func() {
		init("test1", true)
		called := false
		Invoke("", "unknown", "Func1", func() {
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
		Invoke("", "unknown", "Func1", nil)
	}()

	func() {
		init("test3", true)
		called := false
		Invoke("", "testfeature", "unknown", func() {
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
		Invoke("", "testfeature", "unknown", nil)
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
			Invoke("", "testfeature", fname, func() {
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

	func() {
		feature := init("test8", true)
		Invoke("test", "testfeature", "Func3", func() {
			t.Errorf("defaultFunc has been called")
		})
		actual := feature.calledBy
		expected := []string{"Func3:test"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		feature := init("test9", true)
		called := false
		Invoke(1, "testfeature", "Func3", func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called by argument type mismatch")
		}
		actual := feature.calledBy
		expected := []string(nil)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		feature := init("test9", true)
		called := false
		Invoke("test", "testfeature", "Func4", func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called by too many arguments definition")
		}
		actual := feature.calledBy
		expected := []string(nil)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		feature := init("test10", true)
		called := false
		Invoke("test", "testfeature", "Func5", func() {
			called = true
		})
		if !called {
			t.Errorf("defaultFunc hasn't been called by too few arguments definition")
		}
		actual := feature.calledBy
		expected := []string(nil)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		for _, fname := range []string{"Func1", "Func2"} {
			feature := init("test11", true)
			Invoke(nil, "testfeature", fname, func() {
				t.Errorf("defaultFunc has been called")
			})
			actual = feature.calledBy
			expected = []string{fmt.Sprintf("%s:<nil>", fname)}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Expect %q, but %q", expected, actual)
			}
		}
	}()
}

func Test_IsActive(t *testing.T) {
	func() {
		defer func() {
			featureStatus = make(map[string]*status)
		}()
		st := featureStatus["testIsActive"]
		if st != nil {
			t.Fatalf("feature has already been added")
		}
		featureStatus["testIsActive"] = &status{
			feature: &TestFeature{t, "test", true, nil, nil},
			fault:   false,
		}
		actual := IsActive("testIsActive")
		expected := true
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		defer func() {
			featureStatus = make(map[string]*status)
		}()
		st := featureStatus["testIsActive"]
		if st != nil {
			t.Fatalf("feature has already been added")
		}
		actual := IsActive("testIsActive")
		expected := false
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()

	func() {
		defer func() {
			featureStatus = make(map[string]*status)
		}()
		st := featureStatus["testIsActive"]
		if st != nil {
			t.Fatalf("feature has already been added")
		}
		featureStatus["testIsActive"] = &status{
			feature: &TestFeature{t, "test", true, nil, nil},
			fault:   true,
		}
		actual := IsActive("testIsActive")
		expected := false
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %q, but %q", expected, actual)
		}
	}()
}
