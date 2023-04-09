package diff_test

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/fatih/color"
	"github.com/poy/onpar"
	"github.com/poy/onpar/diff"
	"github.com/poy/onpar/diff/str"
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
	o := onpar.New(t)
	defer o.Run()

	for _, tt := range []struct {
		name  string
		value any
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

	o.Spec("it shows the value of pointers when compared against nil", func(t *testing.T) {
		a := "foo"
		out := diff.New().Diff(&a, nil)
		if !strings.Contains(out, "foo") {
			t.Fatalf("output %s should have contained actual value '%s'", out, a)
		}
		out = diff.New().Diff(nil, &a)
		if !strings.Contains(out, "foo") {
			t.Fatalf("output %s should have contained expected value '%s'", out, a)
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
		name             string
		actual, expected any
		output           string
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
			out := diff.New().Diff(tt.actual, tt.expected)
			if out != tt.output {
				t.Fatalf("expected the diff between %v and %v to be %s; got %s", tt.actual, tt.expected, tt.output, out)
			}
		})
	}

	o.Spec("it calls Sprinters", func(t *testing.T) {
		s := newMockSprinter(t, time.Second)
		pers.Return(s.SprintOutput, "foo")
		out := diff.New(diff.WithSprinter(s)).Diff("first", "second")
		expect.Expect(t, s).To(pers.HaveMethodExecuted("Sprint", pers.WithArgs("firstsecond")))
		expect.Expect(t, out).To(matchers.Equal("foo"))
	})

	o.Spec("it does not hang on strings mentioned in issue 30", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			defer close(done)
			diff.New().Diff(
				`{"current":[{"kind":0,"at":{"seconds":1596288863,"nanos":13000000},"msg":"Something bad happened."}]}`,
				`{"current": [{"kind": "GENERIC", "at": "2020-08-01T13:34:23.013Z", "msg": "Something bad happened."}], "history": []}`,
			)
		}()
		select {
		case <-done:
			// This diff is not "stable" - with concurrent differs (like the
			// str.CharDiff differ that we're exercising), we can't guarantee
			// the same output every time.
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for diff to finish")
		}
	})

	o.Spec("it always includes a basic string diff", func(t *testing.T) {
		out := diff.New(diff.WithStringAlgos()).Diff("foobarbaz", "foobaconbaz")
		expect.Expect(t, out).To(matchers.Equal(">foobarbaz!=foobaconbaz<"))
	})

	o.Group("with custom string algorithms", func() {
		type testCtx struct {
			t            *testing.T
			algo1, algo2 *mockStringDiffAlgorithm
			timeout      time.Duration
			differ       *diff.Differ
		}

		o := onpar.BeforeEach(o, func(t *testing.T) testCtx {
			algo1 := newMockStringDiffAlgorithm(t, time.Second)
			algo2 := newMockStringDiffAlgorithm(t, time.Second)
			timeout := time.Second
			return testCtx{
				t:       t,
				algo1:   algo1,
				algo2:   algo2,
				timeout: timeout,
				differ:  diff.New(diff.WithStringAlgos(algo1, algo2), diff.WithTimeout(timeout)),
			}
		})

		o.Spec("it returns the basic diff when no better diffs are returned", func(tt testCtx) {
			ch := make(chan str.Diff)
			close(ch)
			pers.Return(tt.algo1.DiffsOutput, ch)
			pers.Return(tt.algo2.DiffsOutput, ch)
			out := tt.differ.Diff("foobar", "foobaz")
			expect.Expect(tt.t, out).To(matchers.Equal(">foobar!=foobaz<"))
		})

		o.Spec("it returns a better diff from an algorithm", func(tt testCtx) {
			ch := make(chan str.Diff, 1)
			diff := newMockDiff(tt.t, time.Second)
			ch <- diff
			close(ch)

			pers.Return(tt.algo1.DiffsOutput, ch)
			pers.Return(tt.algo2.DiffsOutput, ch)

			pers.Return(diff.CostOutput, 2)
			pers.Return(diff.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("foob"), Expected: []rune("foob")},
				{Type: str.TypeReplace, Actual: []rune("ar"), Expected: []rune("az")},
			})

			out := tt.differ.Diff("foobar", "foobaz")
			expect.Expect(tt.t, out).To(matchers.Equal("foob>ar!=az<"))
		})

		o.Spec("it chooses the best diff regardless of order", func(tt testCtx) {
			ch := make(chan str.Diff, 2)
			diff1 := newMockDiff(tt.t, time.Second)
			diff2 := newMockDiff(tt.t, time.Second)
			ch <- diff1
			ch <- diff2
			close(ch)

			pers.Return(tt.algo1.DiffsOutput, ch)
			pers.Return(tt.algo2.DiffsOutput, ch)

			pers.ConsistentlyReturn(tt.t, diff1.CostOutput, 2)
			pers.ConsistentlyReturn(tt.t, diff1.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("foob"), Expected: []rune("foob")},
				{Type: str.TypeReplace, Actual: []rune("ar"), Expected: []rune("az")},
			})

			pers.ConsistentlyReturn(tt.t, diff2.CostOutput, 5)
			pers.ConsistentlyReturn(tt.t, diff2.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("f"), Expected: []rune("f")},
				{Type: str.TypeReplace, Actual: []rune("oobar"), Expected: []rune("oobaz")},
			})

			out := tt.differ.Diff("foobar", "foobaz")
			expect.Expect(tt.t, out).To(matchers.Equal("foob>ar!=az<"))
		})

		o.Spec("it respects the timeout", func(tt testCtx) {
			ch := make(chan str.Diff, 1)
			diff := newMockDiff(tt.t, time.Second)
			ch <- diff

			pers.Return(tt.algo1.DiffsOutput, ch)
			pers.Return(tt.algo2.DiffsOutput, ch)

			pers.Return(diff.CostOutput, 2)
			pers.Return(diff.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("foob"), Expected: []rune("foob")},
				{Type: str.TypeReplace, Actual: []rune("ar"), Expected: []rune("az")},
			})

			stop := make(chan struct{})
			defer close(stop)
			out := make(chan string)
			go func() {
				defer close(out)
				select {
				case <-stop:
				case out <- tt.differ.Diff("foobar", "foobaz"):
				}
			}()

			gracePeriod := 10 * time.Millisecond
			timeout := time.After(tt.timeout + gracePeriod)
			select {
			case <-timeout:
				tt.t.Fatalf("did not receive output within %v", tt.timeout)
			case out := <-out:
				expect.Expect(tt.t, out).To(matchers.Equal("foob>ar!=az<"))
			}
		})

		o.Spec("it runs algorithms concurrently and returns the lowest cost available", func(tt testCtx) {
			ch1 := make(chan str.Diff, 2)
			ch2 := make(chan str.Diff, 1)
			diff1 := newMockDiff(tt.t, time.Second)
			diff2 := newMockDiff(tt.t, time.Second)
			diff3 := newMockDiff(tt.t, time.Second)
			ch1 <- diff1
			ch2 <- diff2
			ch1 <- diff3

			pers.Return(tt.algo1.DiffsOutput, ch1)
			pers.Return(tt.algo2.DiffsOutput, ch2)

			pers.ConsistentlyReturn(tt.t, diff1.CostOutput, 4)
			pers.Return(diff1.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("fo"), Expected: []rune("fo")},
				{Type: str.TypeReplace, Actual: []rune("obar"), Expected: []rune("obaz")},
			})

			pers.ConsistentlyReturn(tt.t, diff2.CostOutput, 1)
			pers.Return(diff2.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("fooba"), Expected: []rune("fooba")},
				{Type: str.TypeReplace, Actual: []rune("r"), Expected: []rune("z")},
			})

			pers.ConsistentlyReturn(tt.t, diff3.CostOutput, 2)
			pers.Return(diff3.SectionsOutput, []str.DiffSection{
				{Type: str.TypeMatch, Actual: []rune("foob"), Expected: []rune("foob")},
				{Type: str.TypeReplace, Actual: []rune("ar"), Expected: []rune("az")},
			})

			out := tt.differ.Diff("foobar", "foobaz")
			expect.Expect(tt.t, out).To(matchers.Equal("fooba>r!=z<"))
		})
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

func ExampleWithSprinter_color() {
	// WithSprinter is provided for integration with any type
	// that has an `Sprint(...any) string` method.  Here
	// we use github.com/fatih/color.
	styles := []diff.Opt{
		diff.Actual(diff.WithSprinter(color.New(color.CrossedOut, color.FgRed))),
		diff.Expected(diff.WithSprinter(color.New(color.FgYellow))),
	}
	fmt.Println(diff.New(styles...).Diff("some string with foo in it", "some string with bar in it"))
}
