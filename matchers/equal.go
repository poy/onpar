package matchers

import (
	"fmt"
	"reflect"
)

type EqualMatcher struct {
	expected interface{}
}

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
