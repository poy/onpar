package matcher

// AndMatcher is a matcher which combines child matchers, only returning success
// if all children also return success.
type AndMatcher[T any] struct {
	Children []Matcher[T]
}

// And constructs an AndMatcher from the passed in child matchers. The child
// matchers will be called in the order they are passed in until the first
// non-nil error is returned.
func And[T any](a, b Matcher[T], ms ...Matcher[T]) AndMatcher[T] {
	return AndMatcher[T]{
		Children: append([]Matcher[T]{a, b}, ms...),
	}
}

// Match returns a nil error if all children return a nil error, otherwise it
// returns the first error it finds.
//
// Execution of child matchers will stop at the first non-nil error.
func (m AndMatcher[T]) Match(actual T) error {
	for _, child := range m.Children {
		if err := child.Match(actual); err != nil {
			return err
		}
	}
	return nil
}
