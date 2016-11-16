package matchers

import (
	"fmt"
	"strings"
)

type StartsWithMatcher struct {
	prefix string
}

func StartsWith(prefix string) StartsWithMatcher {
	return StartsWithMatcher{
		prefix: prefix,
	}
}

func (m StartsWithMatcher) Match(actual interface{}) (interface{}, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasPrefix(s, m.prefix) {
		return nil, fmt.Errorf("%s does not start with %s", s, m.prefix)
	}

	return actual, nil
}
