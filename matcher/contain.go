package matcher

import (
	"fmt"
)

type ContainMatcher[T ~[]U, U comparable] struct {
	values T
}

func Contain[T ~[]U, U comparable](values ...U) ContainMatcher[T, U] {
	return ContainMatcher[T, U]{
		values: values,
	}
}

func (m ContainMatcher[T, U]) Match(actual T) error {
	var missing []U
	for _, expected := range m.values {
		found := false
		for _, v := range actual {
			if expected == v {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expected)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("expected %v to contain %v; missing elements: %v", actual, m.values, missing)
}
