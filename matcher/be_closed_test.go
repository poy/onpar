package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestBeClosed(t *testing.T) {
	t.Parallel()
	m := matcher.BeClosed[chan int]()

	c := make(chan int)
	close(c)
	if err := m.Match(c); err != nil {
		t.Error("expected err to be nil")
	}

	c = make(chan int)
	if err := m.Match(c); err == nil {
		t.Errorf("expected err to not be nil")
	}
}
