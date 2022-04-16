package matcher

import (
	"fmt"
)

// EqualMatcher performs a DeepEqual between the actual and expected.
type EqualMatcher[T comparable] struct {
	expected T
	differ   Differ
}

// Equal returns an EqualMatcher with the expected value
func Equal[T comparable](expected T) *EqualMatcher[T] {
	return &EqualMatcher[T]{
		expected: expected,
	}
}

func (m *EqualMatcher[T]) UseDiffer(d Differ) {
	m.differ = d
}

func (m EqualMatcher[T]) Match(actual T) error {
	if actual != m.expected {
		if m.differ == nil {
			return fmt.Errorf("expected %+[1]v (%[1]T) to equal %+[2]v (%[2]T)", actual, m.expected)
		}
		return fmt.Errorf("expected %v to equal %v\ndiff: %s", actual, m.expected, m.differ.Diff(actual, m.expected))
	}

	return nil
}
