package matcher

// AnythingMatcher is a matcher that literally matches anything. It's intended
// for use as a sub-matcher when you don't (yet) care about the value being
// matched against, for example in a ReceiveMatcher.
type AnythingMatcher[T any] struct{}

// Anything returns a matcher that matches anything.
func Anything[T any]() AnythingMatcher[T] {
	return AnythingMatcher[T]{}
}

// Match returns nil.
func (m AnythingMatcher[T]) Match(v T) error {
	return nil
}
