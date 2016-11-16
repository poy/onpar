package matchers

import (
	"fmt"
	"reflect"
)

type ReceiveMatcher struct {
}

func Receive() ReceiveMatcher {
	return ReceiveMatcher{}
}

func (m ReceiveMatcher) Match(actual interface{}) (interface{}, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Chan || t.ChanDir() == reflect.SendDir {
		return nil, fmt.Errorf("%s is not a readable channel", t.String())
	}

	v := reflect.ValueOf(actual)
	rxValue, ok := v.TryRecv()

	if !ok {
		return nil, fmt.Errorf("did not receive")
	}

	return rxValue.Interface(), nil
}
