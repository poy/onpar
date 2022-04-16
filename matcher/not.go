package matcher

import "fmt"

// NotMatcher accepts a matcher and will succeed if the specified matcher fails.
type NotMatcher[T any] struct {
	child Matcher[T]
}

// Not returns a NotMatcher with the specified child matcher.
func Not[T any](child Matcher[T]) NotMatcher[T] {
	return NotMatcher[T]{
		child: child,
	}
}

// Match returns an error if the child matcher does not return an error, or nil
// otherwise.
func (m NotMatcher[T]) Match(actual T) error {
	err := m.child.Match(actual)
	if err == nil {
		return fmt.Errorf("%+v (%[1]T) was expected to fail matcher %#v", actual, m.child)
	}

	return nil
}
