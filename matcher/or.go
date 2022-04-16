package matcher

// OrMatcher is a matcher that returns success if any of its children return success.
type OrMatcher[T any] struct {
	Children []Matcher[T]
}

// Or constructs an OrMatcher from the passed in child matchers. The child
// matchers will be called in order until the first nil error is returned.
func Or[T any](a, b Matcher[T], ms ...Matcher[T]) OrMatcher[T] {
	return OrMatcher[T]{
		Children: append([]Matcher[T]{a, b}, ms...),
	}
}

// Match returns a nil error if any child matcher returns a nil error.
// Otherwise, it returns the error from the last child.
//
// Execution of child matchers will stop at the first nil error.
func (m OrMatcher[T]) Match(actual T) error {
	var err error
	for _, child := range m.Children {
		err = child.Match(actual)
		if err == nil {
			return nil
		}
	}
	return err
}
