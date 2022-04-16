package matcher

import (
	"fmt"
	"strings"
)

// StartWithMatcher accepts a string and succeeds
// if the actual string starts with the expected string.
type StartWithMatcher struct {
	prefix string
}

// StartWith returns a StartWithMatcher with the expected prefix.
func StartWith(prefix string) StartWithMatcher {
	return StartWithMatcher{
		prefix: prefix,
	}
}

func (m StartWithMatcher) Match(actual string) error {
	if !strings.HasPrefix(actual, m.prefix) {
		return fmt.Errorf("expected %s to start with %s", actual, m.prefix)
	}

	return nil
}
