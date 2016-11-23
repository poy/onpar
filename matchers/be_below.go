package matchers

import "fmt"

// BeBelowMatcher accepts a float64. It succeeds if the actual is
// less than the expected.
type BeBelowMatcher struct {
	expected float64
}

// BeBelow returns a BeBelowMatcher with the expected value.
func BeBelow(expected float64) BeBelowMatcher {
	return BeBelowMatcher{
		expected: expected,
	}
}

func (m BeBelowMatcher) Match(actual interface{}) (interface{}, error) {
	f, ok := actual.(float64)
	if !ok {
		return nil, fmt.Errorf("%v (%T) is not a float64", actual, actual)
	}

	if f > m.expected {
		return nil, fmt.Errorf("%f is not below %f", actual, m.expected)
	}

	return actual, nil
}
