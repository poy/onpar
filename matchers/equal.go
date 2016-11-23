package matchers

import (
	"fmt"
	"reflect"
)

// EqualMatcher performs a DeepEqual between the actual and expected.
type EqualMatcher struct {
	expected interface{}
}

// Equal returns an EqualMatcher with the expected value
func Equal(expected interface{}) EqualMatcher {
	return EqualMatcher{
		expected: expected,
	}
}

func (m EqualMatcher) Match(actual interface{}) (interface{}, error) {
	if !reflect.DeepEqual(actual, m.expected) {
		return nil, fmt.Errorf("%v to equal %v", actual, m.expected)
	}

	return actual, nil
}
