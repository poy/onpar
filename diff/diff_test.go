package diff_test

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/nelsam/hel/v2/pers"
	"github.com/poy/onpar"
	"github.com/poy/onpar/diff"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
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

	o.Spec("it calls Sprinters", func(t *testing.T) {
		s := newMockSprinter()
		pers.Return(s.SprintOutput, "foo")
		out := diff.New(diff.WithSprinter(s)).Diff("first", "second")
		expect.Expect(t, s).To(pers.HaveMethodExecuted("Sprint", pers.WithArgs("firstsecond")))
		expect.Expect(t, out).To(matchers.Equal("foo"))
	})
}

func ExampleDiffer_Diff() {
	fmt.Println(diff.New().Diff("some string with foo in it", "some string with bar in it"))
	// Output:
	// some string with >foo!=bar< in it
}

func ExampleWithFormat() {
	style := diff.WithFormat("LOOK-->%s<--!!!")
	fmt.Println(diff.New(style).Diff("some string with foo in it", "some string with bar in it"))
	// Output:
	// some string with LOOK-->foobar<--!!! in it
}

func ExampleActual() {
	styles := []diff.Opt{
		diff.Actual( // styles passed to this will only apply to actual values
			diff.WithFormat("(%s|"),
		),
		diff.Expected( // styles passed to this will only apply to expected values
			diff.WithFormat("%s)"),
		),
	}
	fmt.Println(diff.New(styles...).Diff("some string with foo in it", "some string with bar in it"))
	// Output:
	// some string with (foo|bar) in it
}

func ExampleWithSprinter_Color() {
	// WithSprinter is provided for integration with any type
	// that has an `Sprint(...interface{}) string` method.  Here
	// we use github.com/fatih/color.
	styles := []diff.Opt{
		diff.Actual(diff.WithSprinter(color.New(color.CrossedOut, color.FgRed))),
		diff.Expected(diff.WithSprinter(color.New(color.FgYellow))),
	}
	fmt.Println(diff.New(styles...).Diff("some string with foo in it", "some string with bar in it"))
}
