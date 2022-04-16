package matcher

import (
	"fmt"
	"strings"
)

// EndWithMatcher accepts a string and succeeds
// if the actual string ends with the expected string.
type EndWithMatcher struct {
	suffix string
}

// EndWith returns an EndWithMatcher with the expected suffix.
func EndWith(suffix string) EndWithMatcher {
	return EndWithMatcher{
		suffix: suffix,
	}
}

func (m EndWithMatcher) Match(actual string) error {
	if !strings.HasSuffix(actual, m.suffix) {
		return fmt.Errorf("expected %s to end with %s", actual, m.suffix)
	}

	return nil
}
