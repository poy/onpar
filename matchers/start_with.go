package matchers

import (
	"fmt"
	"strings"
)

type StartWithMatcher struct {
	prefix string
}

func StartWith(prefix string) StartWithMatcher {
	return StartWithMatcher{
		prefix: prefix,
	}
}

func (m StartWithMatcher) Match(actual interface{}) (interface{}, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasPrefix(s, m.prefix) {
		return nil, fmt.Errorf("%s does not start with %s", s, m.prefix)
	}

	return actual, nil
}
