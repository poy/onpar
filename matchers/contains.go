package matchers

import (
	"fmt"
	"reflect"
)

type ContainsMatcher struct {
	values []interface{}
}

func Contains(values ...interface{}) ContainsMatcher {
	return ContainsMatcher{
		values: values,
	}
}

func (m ContainsMatcher) Match(actual interface{}) (interface{}, error) {
	actualType := reflect.TypeOf(actual)
	if actualType.Kind() != reflect.Slice && actualType.Kind() != reflect.Array {
		return nil, fmt.Errorf("%s is not a Slice or Array", actualType.Kind())
	}

	actualValue := reflect.ValueOf(actual)
	for _, elem := range m.values {
		if !m.containsElem(actualValue, elem) {
			return nil, fmt.Errorf("%v does not contain %v", actual, elem)
		}
	}

	return actual, nil
}

func (m ContainsMatcher) containsElem(actual reflect.Value, elem interface{}) bool {
	for i := 0; i < actual.Len(); i++ {
		if reflect.DeepEqual(actual.Index(i).Interface(), elem) {
			return true
		}
	}

	return false
}
