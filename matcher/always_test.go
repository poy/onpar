package matcher_test

import (
	"fmt"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matcher"
)

func TestFailsWhenMatcherFails(t *testing.T) {
	t.Parallel()

	match := newMockMatcher[int](t, time.Second)
	m := matcher.Always[func() int, int](match)
	pers.Return(match.MatchOutput, fmt.Errorf("some-error"))

	callCount := make(chan bool, 200)
	f := func() int {
		callCount <- true
		return 99
	}

	if err := m.Match(f); err == nil {
		t.Error("expected err to not be nil")
	}

	if len(callCount) != 1 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 1)
	}
}

func TestPolls10TimesForSuccess(t *testing.T) {
	t.Parallel()
	match := newMockMatcher[int](t, time.Second)
	m := matcher.Always[func() int, int](match)

	i := 0
	callCount := make(chan bool, 200)
	f := func() int {
		i++
		callCount <- true
		pers.Return(match.MatchOutput, nil)
		return i
	}

	if err := m.Match(f); err != nil {
		t.Fatal("expected err to be nil")
	}

	if len(callCount) != 10 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 10)
	}
}

func TestPollsEach10msForSuccess(t *testing.T) {
	t.Parallel()
	match := newMockMatcher[int](t, time.Second)
	m := matcher.Always[func() int, int](match)

	var ts []int64
	f := func() int {
		ts = append(ts, time.Now().UnixNano())
		pers.Return(match.MatchOutput, nil)
		return 101
	}
	m.Match(f)

	for i := 0; i < len(ts)-1; i++ {
		if ts[i+1]-ts[i] < int64(8*time.Millisecond) || ts[i+1]-ts[i] > int64(15*time.Millisecond) {
			t.Fatalf("expected %d to be within 8ms and 15ms", ts[i+1]-ts[i])
		}
	}
}
