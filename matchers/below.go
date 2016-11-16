package matchers

import "fmt"

type BelowMatcher struct {
	expected float64
}

func Below(expected float64) BelowMatcher {
	return BelowMatcher{
		expected: expected,
	}
}

func (m BelowMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(float64)
	if !ok {
		return nil, fmt.Errorf("%v (%T) is not a float64", actual, actual)
	}

	if f > m.expected {
		return nil, fmt.Errorf("%f is not below %f", actual, m.expected)
	}

	return nil, nil
}
