package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MatchJSONMatcher converts both expected and actual to a map[string]interface{}
// and does a reflect.DeepEqual between them
type MatchJSONMatcher[T ~string | ~[]byte] struct {
	expected T
}

// MatchJSON returns an MatchJSONMatcher with the expected value
func MatchJSON[T ~string | ~[]byte](expected T) MatchJSONMatcher[T] {
	return MatchJSONMatcher[T]{
		expected: expected,
	}
}

// Match returns nil if actual matches the expected value, after performing a
// json unmarshal. If either value fails to unmarshal, the unmarshal error will
// be returned instead.
func (m MatchJSONMatcher[T]) Match(actual T) error {
	a, sa, err := m.unmarshal(actual)
	if err != nil {
		return fmt.Errorf("expected actual to unmarshal as json; json.Unmarshal returned: %w", err)
	}

	e, se, err := m.unmarshal(m.expected)
	if err != nil {
		return fmt.Errorf("expected the expected value to unmarshal as json; json.Unmarshal returned: %w", err)
	}

	if !reflect.DeepEqual(a, e) {
		return fmt.Errorf("expected %s to equal %s", sa, se)
	}

	return nil
}

func (m MatchJSONMatcher[T]) unmarshal(x T) (interface{}, string, error) {
	var result interface{}
	switch x := interface{}(x).(type) {
	case []byte:
		if err := json.Unmarshal(x, &result); err != nil {
			return nil, string(x), err
		}
		return result, string(x), nil
	case string:
		if err := json.Unmarshal([]byte(x), &result); err != nil {
			return nil, x, err
		}
		return result, x, nil
	default:
		panic(fmt.Errorf("unexpected type from type union: %T", x))
	}
}
