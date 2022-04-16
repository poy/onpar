package matcher

import (
	"fmt"
)

// HaveCapMatcher is a matcher that checks the capacity of a value.
type HaveCapMatcher[T ~[]U | ~chan U, U any] struct {
	expected int
}

// HaveCap returns a HaveCapMatcher with the specified capacity
func HaveCap[T ~[]U | ~chan U, U any](expected int) HaveCapMatcher[T, U] {
	return HaveCapMatcher[T, U]{
		expected: expected,
	}
}

// Match fails if actual has a capacity that is not equal to the expected
// capacity.
func (m HaveCapMatcher[T, U]) Match(actual T) error {
	if cap(actual) != m.expected {
		return fmt.Errorf("%v (cap=%d) does not have a capacity %d", actual, cap(actual), m.expected)
	}

	return nil
}
