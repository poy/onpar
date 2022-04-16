package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestTruthiness(t *testing.T) {
	t.Parallel()

	m := matcher.BeTrue()

	if err := m.Match(true); err != nil {
		t.Error("expected err to be nil")
	}

	if err := m.Match(false); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match(2 == 2); err != nil {
		t.Error("expected err to be nil")
	}
}
