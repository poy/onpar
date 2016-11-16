package matchers_test

import (
	"testing"

	"github.com/apoydence/onpar/matchers"
)

func TestStartsWith(t *testing.T) {
	t.Parallel()

	m := matchers.StartsWith("foo")

	_, err := m.Match("bar")
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match("foobar")
	if err != nil {
		t.Error("expected err to be nil")
	}

	_, err = m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
