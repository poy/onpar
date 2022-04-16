package matcher

import (
	"fmt"
)

type haveKeyPrefs[T any] struct {
	valueMatcher Matcher[T]
}

// HaveKeyOpt is an option type that may be passed to HaveKey to make optional
// changes to behavior (e.g. checking the value at the specified key).
type HaveKeyOpt[T any] func(haveKeyPrefs[T]) haveKeyPrefs[T]

// WithValue takes a matcher and applies the matcher against the value at the
// specified key, if it exists.
func WithValue[T any](m Matcher[T]) HaveKeyOpt[T] {
	return func(p haveKeyPrefs[T]) haveKeyPrefs[T] {
		p.valueMatcher = m
		return p
	}
}

// HaveKeyMatcher succeeds if the map contains the specified key.
type HaveKeyMatcher[T ~map[K]V, K comparable, V any] struct {
	key        K
	valMatcher Matcher[V]
}

// HaveKey returns a HaveKeyMatcher with the specified key.
func HaveKey[T ~map[K]V, K comparable, V any](key K, opts ...HaveKeyOpt[V]) HaveKeyMatcher[T, K, V] {
	var p haveKeyPrefs[V]
	for _, o := range opts {
		p = o(p)
	}
	return HaveKeyMatcher[T, K, V]{
		key:        key,
		valMatcher: p.valueMatcher,
	}
}

// Match looks up the specified key in actual, erroring if it doesn't exist. If
// WithValue was specified, then the resulting value at the specified key will
// be checked against the value matcher as well.
func (m HaveKeyMatcher[T, K, V]) Match(actual T) error {
	v, ok := actual[m.key]
	if !ok {
		return fmt.Errorf("expected map %v to contain key %v", actual, m.key)
	}

	if m.valMatcher != nil {
		if err := m.valMatcher.Match(v); err != nil {
			return fmt.Errorf("expected value %v (from key %v) to pass sub-matcher, but it failed: %w", v, m.key, err)
		}
	}

	return nil
}
