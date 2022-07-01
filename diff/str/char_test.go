package str_test

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/poy/onpar/v2"
	"github.com/poy/onpar/v2/diff/str"
	"github.com/poy/onpar/v2/expect"
	"github.com/poy/onpar/v2/matchers"
)

func TestCharDiff(t *testing.T) {
	o := onpar.New(t)

	matchReplace := func(t *testing.T, start str.Type, l ...string) []str.DiffSection {
		// matchReplace is used to generate match/replace cadence for expected
		// values, since by default the char diff will always return diffs that
		// have a match followed by a replace followed by a match, and so on.

		t.Helper()

		var sections []str.DiffSection
		curr := start
		for _, v := range l {
			section := str.DiffSection{
				Type:     curr,
				Actual:   []rune(v),
				Expected: []rune(v),
			}
			if curr == str.TypeReplace {
				// NOTE: if any of the diffs we need to test contain a pipe
				// character, this will break.
				parts := strings.Split(v, "|")
				if len(parts) != 2 {
					t.Fatalf("test error: expected replace string (%v) to use a | character to separate actual and expected values, but got %d values when splitting", v, len(parts))
				}
				section.Actual = []rune(parts[0])
				section.Expected = []rune(parts[1])
			}
			sections = append(sections, section)
			curr = 1 - curr
		}
		return sections
	}

	replace := func(t *testing.T, a, b string) string {
		// replace is used to generate replacement strings for the matchReplace
		// function. This helps make the test read more clearly, while also
		// checking for separator characters in the source strings.
		t.Helper()

		if strings.Contains(a, "|") {
			t.Fatalf("replace source string %v contains pipe character", a)
		}
		if strings.Contains(b, "|") {
			t.Fatalf("replace source string %v contains pipe character", b)
		}
		return fmt.Sprintf("%s|%s", a, b)
	}

	o.Group("exhaustive results", func() {
		finalDiff := func(t *testing.T, timeout time.Duration, diffs <-chan str.Diff) str.Diff {
			t.Helper()

			var final str.Diff
			tCh := time.After(timeout)
			for {
				select {
				case next, ok := <-diffs:
					if !ok {
						return final
					}
					final = next
				case <-tCh:
					t.Fatalf("failed to exhaust results within %v", timeout)
				}
			}
		}

		for _, tt := range []struct {
			name             string
			actual, expected string
			output           []str.DiffSection
		}{
			{"different strings", "foo", "bar", []str.DiffSection{{Type: str.TypeReplace, Actual: []rune("foo"), Expected: []rune("bar")}}},
			{"different substrings", "foobarbaz", "fooeggbaz", matchReplace(t, str.TypeMatch, "foo", replace(t, "bar", "egg"), "baz")},
			{"longer expected string", "foobarbaz", "foobazingabaz", matchReplace(t, str.TypeMatch, "fooba", replace(t, "r", "zinga"), "baz")},
			{"longer actual string", "foobazingabaz", "foobarbaz", matchReplace(t, str.TypeMatch, "fooba", replace(t, "zinga", "r"), "baz")},
			{"multiple different substrings", "pythonfooeggsbazingabacon", "gofoobarbazingabaz",
				matchReplace(t, str.TypeReplace,
					replace(t, "pyth", "g"),
					"o",
					replace(t, "n", ""),
					"foo",
					replace(t, "eggs", "bar"),
					"bazingaba",
					replace(t, "con", "z"),
				),
			},
		} {
			tt := tt
			o.Spec(tt.name, func(t *testing.T) {
				ch := str.NewCharDiff().Diffs(context.Background(), []rune(tt.actual), []rune(tt.expected))
				final := finalDiff(t, time.Second, ch)

				exp := readableSections(tt.output)
				for i, v := range readableSections(final.Sections()) {
					// NOTE: these matches will be checked again below;
					// this is just to provide more detail about which
					// indexes failed. We could use a differ, but since
					// this test is testing differs, a bug in the differ
					// might break the test (rather than simply failing
					// it).
					if i > len(exp) {
						t.Errorf("actual (length %d) was longer than expected (length %d)", len(final.Sections()), len(exp))
						break
					}
					if !reflect.DeepEqual(v, exp[i]) {
						t.Errorf("%#v did not match %#v", v, exp[i])
					}
				}
				expect.Expect(t, readableSections(final.Sections())).To(matchers.Equal(readableSections(tt.output)))
			})
		}

		o.Spec("it doesn't hang on strings mentioned in issue 30", func(t *testing.T) {
			// This diff has multiple options for the "best" result (more than
			// one diff at the lowest possible cost). So we can't very well
			// assert on the exact diff returned, but we can assert on the cost
			// and the total actual and expected strings.
			actual := `{"current":[{"kind":0,"at":{"seconds":1596288863,"nanos":13000000},"msg":"Something bad happened."}]}`
			expected := `{"current": [{"kind": "GENERIC", "at": "2020-08-01T13:34:23.013Z", "msg": "Something bad happened."}], "history": []}`
			diffs := str.NewCharDiff().Diffs(context.Background(), []rune(actual), []rune(expected))
			final := finalDiff(t, time.Second, diffs)

			expectedCost := float64(57)
			expect.Expect(t, final.Cost()).To(matchers.Equal(expectedCost))

			var retActual, retExpected string
			for _, v := range final.Sections() {
				retActual += string(v.Actual)
				retExpected += string(v.Expected)
			}
			expect.Expect(t, retActual).To(matchers.Equal(actual))
			expect.Expect(t, retExpected).To(matchers.Equal(expected))
		})
	})

	o.Spec("it returns decreasingly costly results until the context is done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		actual := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"
		expected := "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
		diffs := str.NewCharDiff().Diffs(ctx, []rune(actual), []rune(expected))

		// We want to check a few results, ensuring that the cost is lower each
		// time, before cancelling the context to ensure that the results stop
		// when the context finishes.
		const toCheck = 5
		best := math.MaxFloat64
		for i := 0; i < toCheck; i++ {
			v := <-diffs
			if v.Cost() > best {
				t.Fatalf("new cost (%f) is greater than old cost (%f)", v.Cost(), best)
			}
			best = v.Cost()
		}

		cancel()
		// Ensure that the logic has time to shut down before we resume reading
		// from the channel.
		time.Sleep(100 * time.Millisecond)
		if _, ok := <-diffs; ok {
			t.Fatalf("results channel is still open after cancelling the context")
		}
	})
}

type readableSection struct {
	typ              str.Type
	actual, expected string
}

func readableSections(s []str.DiffSection) []readableSection {
	var rs []readableSection
	for _, v := range s {
		rs = append(rs, readableSection{
			typ:      v.Type,
			actual:   string(v.Actual),
			expected: string(v.Expected),
		})
	}
	return rs
}
