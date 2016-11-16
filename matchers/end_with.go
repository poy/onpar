package matchers

import (
	"fmt"
	"strings"
)

type EndsWithMatcher struct {
	suffix string
}

func EndsWith(suffix string) EndsWithMatcher {
	return EndsWithMatcher{
		suffix: suffix,
	}
}

func (m EndsWithMatcher) Match(actual interface{}) (interface{}, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasSuffix(s, m.suffix) {
		return nil, fmt.Errorf("%s does not end with %s", s, m.suffix)
	}

	return nil, nil
}
