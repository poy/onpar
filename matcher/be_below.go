package matcher

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// BeBelowMatcher succeeds if the actual is less than the expected.
type BeBelowMatcher[T constraints.Ordered] struct {
	expected T
}

// BeBelow returns a BeBelowMatcher with the expected value.
func BeBelow[T constraints.Ordered](expected T) BeBelowMatcher[T] {
	return BeBelowMatcher[T]{
		expected: expected,
	}
}

func (m BeBelowMatcher[T]) Match(actual T) error {
	if actual >= m.expected {
		return fmt.Errorf("expected %v to be below %v", actual, m.expected)
	}

	return nil
}
