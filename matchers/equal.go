package matchers

import (
	"fmt"
	"reflect"

	"github.com/poy/onpar/diff"
)

// EqualMatcher performs a DeepEqual between the actual and expected.
type EqualMatcher struct {
	expected interface{}
	diffOpts []diff.Opt
}

// Equal returns an EqualMatcher with the expected value
func Equal(expected interface{}) *EqualMatcher {
	return &EqualMatcher{
		expected: expected,
	}
}

func (m *EqualMatcher) UseDiffOpts(opts ...diff.Opt) {
	m.diffOpts = opts
}

func (m EqualMatcher) Match(actual interface{}) (interface{}, error) {
	if !reflect.DeepEqual(actual, m.expected) {
		return nil, fmt.Errorf("%v to equal %v\ndiff: %s", actual, m.expected, diff.Show(actual, m.expected, m.diffOpts...))
	}

	return actual, nil
}
