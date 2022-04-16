package matcher

import "fmt"

// BeFalseMatcher will succeed if actual is false.
type BeFalseMatcher struct{}

// BeFalse will return a BeFalseMatcher
func BeFalse() BeFalseMatcher {
	return BeFalseMatcher{}
}

func (m BeFalseMatcher) Match(actual bool) error {
	if actual {
		return fmt.Errorf("expected %t to be false", actual)
	}

	return nil
}
