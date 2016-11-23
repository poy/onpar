package matchers

import (
	"fmt"
	"reflect"
)

type BeClosedMatcher struct {
}

func BeClosed() BeClosedMatcher {
	return BeClosedMatcher{}
}

func (m BeClosedMatcher) Match(actual interface{}) (interface{}, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Chan || t.ChanDir() == reflect.SendDir {
		return nil, fmt.Errorf("%s is not a readable channel", t.String())
	}

	v := reflect.ValueOf(actual)

	winnerIndex, _, open := reflect.Select([]reflect.SelectCase{
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: v},
		reflect.SelectCase{Dir: reflect.SelectDefault},
	})

	if winnerIndex == 0 && !open {
		return actual, nil
	}

	return nil, fmt.Errorf("channel open")
}
