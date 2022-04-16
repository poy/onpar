package matcher

import "fmt"

// HaveOccurredMatcher will succeed if the actual value is a non-nil error.
type HaveOccurredMatcher struct {
}

// HaveOccurred returns a HaveOccurredMatcher
func HaveOccurred() HaveOccurredMatcher {
	return HaveOccurredMatcher{}
}

func (m HaveOccurredMatcher) Match(actual error) error {
	if actual == nil {
		return fmt.Errorf("expected a non-nil error; got %v", actual)
	}

	return nil
}
