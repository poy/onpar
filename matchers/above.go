package matchers

import "fmt"

type AboveMatcher struct {
	expected float64
}

func Above(expected float64) AboveMatcher {
	return AboveMatcher{
		expected: expected,
	}
}

func (m AboveMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(float64)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a float64", actual, actual)
	}

	if f <= m.expected {
		return nil, fmt.Errorf("%f is not above %f", actual, m.expected)
	}

	return nil, nil
}
