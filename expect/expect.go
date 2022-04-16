package expect

import (
	"path"
	"runtime"

	"github.com/poy/onpar/v2/matcher"
)

// ToMatcher is an expectation which can perform matches against actual values.
type ToMatcher[A any] interface {
	Match(actual A) error
}

// Differ is a type of matcher that will need to diff its expected and
// actual values.
type DiffMatcher interface {
	UseDiffer(matcher.Differ)
}

// T is a type that we can perform assertions with.
type T interface {
	Fatalf(format string, args ...interface{})
}

// THelper has the method that tells the testing framework that it can declare
// itself a test helper.
type THelper interface {
	Helper()
}

type prefs struct {
	differ matcher.Differ
}

// Opt is an option that can be passed to New to modify Expectations.
//
// NOTE: generics do not yet infer types in functional options, so these option
// functions accept a non-generic unexported struct instead. Preferences are
// then loaded from the struct into the Expectation.
//
// TODO: if/when generics can infer types in functional options, it's probably
// best (for documentation purposes) to make options apply directly to To types
// again.
type Opt func(prefs) prefs

// WithDiffer stores the diff.Differ to be used when displaying diffs between
// actual and expected values.
func WithDiffer(d matcher.Differ) Opt {
	return func(p prefs) prefs {
		p.differ = d
		return p
	}
}

// Expect takes a T (usually *testing.T) and an actual value, returning a To
// which may be used to apply matchers against the actual value.
func Expect[A any](t T, actual A, opts ...Opt) *To[A] {
	to := To[A]{
		actual: actual,
		t:      t,
	}
	var p prefs
	for _, opt := range opts {
		p = opt(p)
	}
	if p.differ != nil {
		to.differ = p.differ
	}
	return &to
}

// To is a type that stores actual values prior to running them through
// matchers.
type To[A any] struct {
	actual    A
	parentErr error

	t      T
	differ matcher.Differ
}

// To takes a matcher and passes it the actual value, failing t's T value
// if the matcher returns an error.
func (t *To[A]) To(matcher ToMatcher[A]) {
	if helper, ok := t.t.(THelper); ok {
		helper.Helper()
	}

	if d, ok := matcher.(DiffMatcher); ok && t.differ != nil {
		d.UseDiffer(t.differ)
	}

	if err := matcher.Match(t.actual); err != nil {
		_, fileName, lineNumber, _ := runtime.Caller(1)
		t.t.Fatalf("%s\n%s:%d", err.Error(), path.Base(fileName), lineNumber)
	}
}
