package matchers_test

import (
	"testing"

	"github.com/apoydence/onpar/matchers"
)

func TestCap(t *testing.T) {
	m := matchers.HaveCap(5)

	_, err := m.Match(make([]int, 0, 3))
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(make([]int, 0, 5))
	if err != nil {
		t.Error("expected err to be nil")
	}

	_, err = m.Match(make(chan int, 3))
	if err == nil {
		t.Error("expected err to not be nil")
	}

	c := make(chan int, 5)
	_, err = m.Match(c)
	if err != nil {
		t.Error("expected err to be nil")
	}

	_, err = m.Match(123)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
