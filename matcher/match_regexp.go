package matcher

import (
	"fmt"
	"regexp"
)

type MatchRegexpMatcher struct {
	pattern string
}

func MatchRegexp(pattern string) MatchRegexpMatcher {
	return MatchRegexpMatcher{
		pattern: pattern,
	}
}

func (m MatchRegexpMatcher) Match(actual string) error {
	r, err := regexp.Compile(m.pattern)
	if err != nil {
		return fmt.Errorf("MatchRegexp was passed an invalid pattern (%v): %w", m.pattern, err)
	}

	if !r.MatchString(actual) {
		return fmt.Errorf("expected %s to match pattern %s", actual, m.pattern)
	}

	return nil
}
