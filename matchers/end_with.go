package matchers

import (
	"fmt"
	"strings"
)

type EndWithMatcher struct {
	suffix string
}

func EndWith(suffix string) EndWithMatcher {
	return EndWithMatcher{
		suffix: suffix,
	}
}

func (m EndWithMatcher) Match(actual interface{}) (interface{}, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasSuffix(s, m.suffix) {
		return nil, fmt.Errorf("%s does not end with %s", s, m.suffix)
	}

	return actual, nil
}
