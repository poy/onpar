package matchers

import (
	"fmt"
	"reflect"
)

// ReceiveMatcher only accepts a readable channel. It will error for anything else.
// It will attempt to receive from the channel but will not block.
// It fails if the channel is closed.
type ReceiveMatcher struct{}

// Receive will return a ReceiveMatcher
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
