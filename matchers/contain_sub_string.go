package matchers

import (
	"fmt"
	"strings"
)

type ContainSubstringMatcher struct {
	substr string
}

func ContainSubstring(substr string) ContainSubstringMatcher {
	return ContainSubstringMatcher{
		substr: substr,
	}
}

func (m ContainSubstringMatcher) Match(actual interface{}) (interface{}, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.Contains(s, m.substr) {
		return nil, fmt.Errorf("%s does not contain %s", s, m.substr)
	}

	return actual, nil
}
