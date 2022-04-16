package matcher

import (
	"fmt"
)

// BeClosedMatcher succeeds if the channel is closed.
type BeClosedMatcher[T ~<-chan U | ~chan U, U any] struct{}

// BeClosed returns a BeClosedMatcher
func BeClosed[T ~<-chan U | ~chan U, U any]() BeClosedMatcher[T, U] {
	return BeClosedMatcher[T, U]{}
}

func (m BeClosedMatcher[T, U]) Match(actual T) error {
	select {
	case _, ok := <-actual:
		if !ok {
			return nil
		}
	default:
	}
	return fmt.Errorf("expected channel to be closed, but it was not")
}
