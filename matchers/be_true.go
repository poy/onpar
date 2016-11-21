package matchers

import "fmt"

type BeTrueMatcher struct{}

func BeTrue() BeTrueMatcher {
	return BeTrueMatcher{}
}

func (m BeTrueMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(bool)
	if !ok {
		return nil, fmt.Errorf("'%v' (%[1]T) is not a bool", actual)
	}

	if !f {
		return nil, fmt.Errorf("%t is not true", actual)
	}
	return actual, nil
}
