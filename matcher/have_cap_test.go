package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestCap(t *testing.T) {
	sliceMatcher := matcher.HaveCap[[]int, int](5)

	if err := sliceMatcher.Match(make([]int, 0, 3)); err == nil {
		t.Error("expected err to not be nil")
	}

	x := make([]int, 0, 5)
	if err := sliceMatcher.Match(x); err != nil {
		t.Error("expected err to be nil")
	}

	chanMatcher := matcher.HaveCap[chan int, int](5)
	if err := chanMatcher.Match(make(chan int, 3)); err == nil {
		t.Error("expected err to not be nil")
	}

	c := make(chan int, 5)
	if err := chanMatcher.Match(c); err != nil {
		t.Error("expected err to be nil")
	}
}
