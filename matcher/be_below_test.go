package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestBelow(t *testing.T) {
	t.Parallel()

	m := matcher.BeBelow(101)

	if err := m.Match(99.0); err != nil {
		t.Error("expected err to be nil")
	}

	if err := m.Match(int(99)); err != nil {
		t.Error("expected err to be nil")
	}

	if err := m.Match(101.0); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match(103.0); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match(int(103)); err == nil {
		t.Error("expected err to not be nil")
	}
}
