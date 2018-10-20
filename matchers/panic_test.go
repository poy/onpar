package matchers_test

import (
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	m := matchers.Panic()

	_, err := m.Match(func() { panic("error") })
	if err != nil {
		t.Error("expected err to be nil")
	}

	_, err = m.Match(101.0)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(func() {})
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match("invalid")
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
