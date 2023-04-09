package matchers_test

import (
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestHaveKey(t *testing.T) {
	m := matchers.HaveKey(99)

	_, err := m.Match(map[int]string{1: "a"})
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err := m.Match(map[int]string{99: "a"})
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	if v != "a" {
		t.Errorf("expected %v to equal %v", v, "a")
	}

	_, err = m.Match(map[string]string{"foo": "a"})
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
