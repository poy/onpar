package matchers

import "fmt"

type BeAboveMatcher struct {
	expected float64
}

func BeAbove(expected float64) BeAboveMatcher {
	return BeAboveMatcher{
		expected: expected,
	}
}

func (m BeAboveMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(float64)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a float64", actual, actual)
	}

	if f <= m.expected {
		return nil, fmt.Errorf("%f is not above %f", actual, m.expected)
	}

	return actual, nil
}
