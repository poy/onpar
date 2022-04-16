package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestLen(t *testing.T) {
	sliceMatcher := matcher.HaveLen[[]int, int, int](5)

	if err := sliceMatcher.Match([]int{1, 2, 3}); err == nil {
		t.Error("expected err to not be nil")
	}

	x := []int{1, 2, 3, 4, 5}
	if err := sliceMatcher.Match(x); err != nil {
		t.Error("expected err to be nil")
	}

	chanMatcher := matcher.HaveLen[chan int, int, int](5)
	if err := chanMatcher.Match(make(chan int, 3)); err == nil {
		t.Error("expected err to not be nil")
	}

	c := make(chan int, 10)
	for i := 0; i < 5; i++ {
		c <- i
	}
	if err := chanMatcher.Match(c); err != nil {
		t.Error("expected err to be nil")
	}

	mapMatcher := matcher.HaveLen[map[int]bool, bool, int](5)
	if err := mapMatcher.Match(map[int]bool{1: true, 2: true, 3: true}); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := mapMatcher.Match(map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}); err != nil {
		t.Error("expected err to be nil")
	}
}
