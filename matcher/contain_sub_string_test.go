package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestContainSubstring(t *testing.T) {
	t.Parallel()

	m := matcher.ContainSubstring("foo")

	if err := m.Match("bar"); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match("foobar"); err != nil {
		t.Error("expected err to be nil")
	}
}
