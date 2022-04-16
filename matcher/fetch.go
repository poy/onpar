package matcher

import (
	"fmt"
)

// FetchToMatcher may be used to store the actual value to an address.
type FetchToMatcher[T any] struct {
	target *T
}

// FetchTo returns a FetchToMatcher that will store the actual to target.
func FetchTo[T any](target *T) FetchToMatcher[T] {
	return FetchToMatcher[T]{
		target: target,
	}
}

// Match stores the actual value to m's target. Match always returns a nil
// error.
func (m FetchToMatcher[T]) Match(actual T) error {
	if m.target == nil {
		return fmt.Errorf("cannot store actual to a nil pointer")
	}
	*m.target = actual
	return nil
}
