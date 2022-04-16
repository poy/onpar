package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	m := matcher.Panic()

	if err := m.Match(func() { panic("error") }); err != nil {
		t.Error("expected err to be nil")
	}

	if err := m.Match(func() {}); err == nil {
		t.Error("expected err to not be nil")
	}
}
