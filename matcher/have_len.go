package matcher

import (
	"fmt"
)

// HaveLenMatcher succeeds if the type has the specified length.
type HaveLenMatcher[T ~[]U | ~map[K]U | ~chan U, U any, K comparable] struct {
	expected int
}

// HaveLen returns a HaveLenMatcher with the specified length.
func HaveLen[T ~[]U | ~map[K]U | ~chan U, U any, K comparable](expected int) HaveLenMatcher[T, U, K] {
	return HaveLenMatcher[T, U, K]{
		expected: expected,
	}
}

func (m HaveLenMatcher[T, U, K]) Match(actual T) error {
	if len(actual) != m.expected {
		return fmt.Errorf("expected %v (len=%d) to have a length of %d", actual, len(actual), m.expected)
	}

	return nil
}
