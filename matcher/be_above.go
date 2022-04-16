package matcher

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// BeAboveMatcher succeeds if the actual is greater than the expected.
type BeAboveMatcher[T constraints.Ordered] struct {
	expected T
}

// BeAbove returns a BeAboveMatcher with the expected value.
func BeAbove[T constraints.Ordered](expected T) BeAboveMatcher[T] {
	return BeAboveMatcher[T]{
		expected: expected,
	}
}

func (m BeAboveMatcher[T]) Match(actual T) error {
	if actual <= m.expected {
		return fmt.Errorf("expected %v to be above %v", actual, m.expected)
	}

	return nil
}
