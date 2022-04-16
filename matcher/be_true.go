package matcher

import "fmt"

// BeTrueMatcher will succeed if actual is true.
type BeTrueMatcher struct{}

// BeTrue will return a BeTrueMatcher
func BeTrue() BeTrueMatcher {
	return BeTrueMatcher{}
}

func (m BeTrueMatcher) Match(actual bool) error {
	if !actual {
		return fmt.Errorf("%t is not true", actual)
	}
	return nil
}
