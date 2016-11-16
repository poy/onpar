package matchers

import "fmt"

type HaveOccurredMatcher struct {
}

func HaveOccurred() HaveOccurredMatcher {
	return HaveOccurredMatcher{}
}

func (m HaveOccurredMatcher) Match(actual interface{}) (interface{}, error) {
	e, ok := actual.(error)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not an error", actual, actual)
	}

	if e == nil {
		return nil, fmt.Errorf("err to not be nil")
	}

	return nil, nil
}
