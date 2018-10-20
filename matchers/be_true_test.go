package matchers_test

import (
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestTruthiness(t *testing.T) {
	t.Parallel()

	m := matchers.BeTrue()

	_, err := m.Match(42)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err := m.Match(true)
	if err != nil {
		t.Error("expected err to be nil")
	}
	if v != true {
		t.Errorf("expected %v to be true", v)
	}

	_, err = m.Match(false)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err = m.Match(2 == 2)
	if err != nil {
		t.Error("expected err to be nil")
	}
	if v != true {
		t.Errorf("expected %v to be true", v)
	}
}
