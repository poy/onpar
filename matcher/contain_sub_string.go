package matcher

import (
	"fmt"
	"strings"
)

// ContainSubstringMatcher accepts a string and succeeds
// if the actual string contains the expected string.
type ContainSubstringMatcher struct {
	substr string
}

// ContainSubstring returns a ContainSubstringMatcher with the
// expected substring.
func ContainSubstring(substr string) ContainSubstringMatcher {
	return ContainSubstringMatcher{
		substr: substr,
	}
}

func (m ContainSubstringMatcher) Match(actual string) error {
	if !strings.Contains(actual, m.substr) {
		return fmt.Errorf("expected %s to contain %s", actual, m.substr)
	}

	return nil
}
