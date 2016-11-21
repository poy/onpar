package matchers

import "fmt"

type beTrueMatcher struct{}

func BeTrue() beTrueMatcher {
	return beTrueMatcher{}
}

func (m beTrueMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(bool)
	if !ok {
		return nil, fmt.Errorf("'%v' (%[1]T) is not a bool", actual)
	}

	if !f {
		return nil, fmt.Errorf("%t is not true", actual)
	}
	return actual, nil
}
