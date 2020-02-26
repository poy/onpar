package diff_test

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/poy/onpar"
	"github.com/poy/onpar/diff"
)

type testNestedStruct struct {
	testStruct

	T testStruct
}

type testStruct struct {
	Foo string
	Bar int

	// unexported is not explicitly tested, but if our code attempts to compare
	// unexported fields, it will panic.  This is here to ensure that we notice
	// that.
	unexported *testStruct
}

func TestDiff(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	for _, tt := range []struct {
		name  string
		value interface{}
	}{
		{"string", "this is a string"},
		{"int", 21},
		{"slice", []string{"this", "is", "a", "slice", "of", "strings"}},
		{"map", map[string]string{"maps": "should", "work": "too"}},
		{"struct", testStruct{Foo: "foo", Bar: 42}},
		{"pointer", &testStruct{Foo: "foo", Bar: 42}},
		{"nested structs", testNestedStruct{testStruct: testStruct{Foo: "foo", Bar: 42}, T: testStruct{Foo: "baz", Bar: 42069}}},
	} {
		tt := tt
		o.Spec(fmt.Sprintf("it does not call option functions for matching %s types", tt.name), func(t *testing.T) {
			out := diff.New(diff.WithFormat("!!!FAIL!!!%s!!!FAIL!!!")).Diff(tt.value, tt.value)
			if strings.Contains(out, "!!!FAIL!!!") {
				t.Fatalf("expected matching output to return without formatting; got %s", out)
			}
		})
	}

	o.Spec("it doesn't care if pointer values are different", func(t *testing.T) {
		out := diff.New(diff.WithFormat("!!!FAIL!!!%s!!!FAIL!!!")).Diff(&testStruct{}, &testStruct{})
		if strings.Contains(out, "!!!FAIL!!!") {
			t.Fatalf("expected different pointer values to recursively compare")
		}
	})

	o.Spec("it can handle nil values", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("nil values panicked (*diff.Differ).Diff: %v\n%s", r, string(debug.Stack()))
			}
		}()
		diff.New().Diff(nil, nil)
		diff.New().Diff(nil, 42)
		diff.New().Diff(42, nil)
	})

	o.Spec("it diffs maps", func(t *testing.T) {
		a := map[string]string{"foo": "bar", "baz": "bazinga"}
		b := map[string]string{"foo": "baz", "bazinga": "baz"}
		expectedSubstrs := []string{"foo: ba>r!=z<", ">missing key bazinga!=bazinga: baz<", ">extra key baz: bazinga!=baz: nil<"}

		out := diff.New().Diff(a, b)
		for _, s := range expectedSubstrs {
			if !strings.Contains(out, s) {
				t.Fatalf("expected substring '%s' to exist in '%s'", s, out)
			}
		}
	})

	for _, tt := range []struct {
		name     string
		a, b     interface{}
		expected string
	}{
		{"different strings", "foo", "bar", ">foo!=bar<"},
		{"different ints", 12, 14, ">12!=14<"},
		{"different substrings", "foobarbaz", "fooeggbaz", "foo>bar!=egg<baz"},
		{"longer expected string", "foobarbaz", "foobazingabaz", "fooba>r!=zinga<baz"},
		{"longer actual string", "foobazingabaz", "foobarbaz", "fooba>zinga!=r<baz"},
		{"multiple different substrings", "pythonfooeggsbazingabacon", "gofoobarbazingabaz", ">pyth!=g<o>n!=<foo>eggs!=bar<bazingaba>con!=z<"},
		{"different struct fields", testStruct{Foo: "foo", Bar: 42}, testStruct{Foo: "bar", Bar: 42}, "diff_test.testStruct{Foo: >foo!=bar<, Bar: 42}"},
		{"different anonymous fields",
			testNestedStruct{testStruct: testStruct{Foo: "foo", Bar: 42}, T: testStruct{Foo: "baz", Bar: 42069}},
			testNestedStruct{testStruct: testStruct{Foo: "bar", Bar: 42}, T: testStruct{Foo: "baz", Bar: 42069}},
			"diff_test.testNestedStruct{testStruct: diff_test.testStruct{Foo: >foo!=bar<, Bar: 42}, T: diff_test.testStruct{Foo: baz, Bar: 42069}}"},
		{"different nested fields",
			testNestedStruct{testStruct: testStruct{Foo: "foo", Bar: 42}, T: testStruct{Foo: "baz", Bar: 42069}},
			testNestedStruct{testStruct: testStruct{Foo: "foo", Bar: 42}, T: testStruct{Foo: "bazinga", Bar: 42069}},
			"diff_test.testNestedStruct{testStruct: diff_test.testStruct{Foo: foo, Bar: 42}, T: diff_test.testStruct{Foo: baz>!=inga<, Bar: 42069}}"},
	} {
		tt := tt
		o.Spec(fmt.Sprintf("it shows diffs for %s", tt.name), func(t *testing.T) {
			out := diff.New().Diff(tt.a, tt.b)
			if out != tt.expected {
				t.Fatalf("expected the diff between %v and %v to be %s; got %s", tt.a, tt.b, tt.expected, out)
			}
		})
	}
}
