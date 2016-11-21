package matchers

import "fmt"

type BeFalseMatcher struct{}

func BeFalse() BeFalseMatcher {
	return BeFalseMatcher{}
}

func (m BeFalseMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(bool)
	if !ok {
		return nil, fmt.Errorf("'%v' (%[1]T) is not a bool", actual)
	}

	if f {
		return nil, fmt.Errorf("%t is not false", actual)
	}

	return actual, nil
}
