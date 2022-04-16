package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestEndWith(t *testing.T) {
	t.Parallel()

	m := matcher.EndWith("foo")

	if err := m.Match("bar"); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match("barfoo"); err != nil {
		t.Error("expected err to be nil")
	}
}
