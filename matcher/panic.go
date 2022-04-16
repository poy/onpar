package matcher

import "errors"

// PanicMatcher accepts a function. It succeeds if the function panics.
type PanicMatcher struct {
}

// Panic returns a Panic matcher.
func Panic() PanicMatcher {
	return PanicMatcher{}
}

// Match errors if the actual function does not panic.
func (m PanicMatcher) Match(actual func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			err = errors.New("expected to panic")
		}
	}()

	actual()

	return nil
}
