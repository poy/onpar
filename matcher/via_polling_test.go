package matcher_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matcher"
)

func TestEventuallyFailsPolls10Times(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[func() int, int](mockMatcher)

	callCount := 0
	f := func() int {
		callCount++
		pers.Return(mockMatcher.MatchOutput, fmt.Errorf("still wrong"))
		return 99
	}
	m.Match(f)

	if callCount != 10 {
		t.Errorf("expected callCount (%d) to equal %d", callCount, 10)
	}
}

func TestEventuallyPollsEvery10ms(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[func() int, int](mockMatcher)

	var ts []int64
	f := func() int {
		ts = append(ts, time.Now().UnixNano())
		pers.Return(mockMatcher.MatchOutput, fmt.Errorf("still wrong"))
		return 99
	}
	m.Match(f)

	for i := 0; i < len(ts)-1; i++ {
		if ts[i+1]-ts[i] < int64(8*time.Millisecond) || ts[i+1]-ts[i] > int64(15*time.Millisecond) {
			t.Fatalf("expected %v to be within 8ms and 15ms", time.Duration(ts[i+1]-ts[i]))
		}
	}
}

func TestEventuallyStopsAfterSuccess(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[func() int, int](mockMatcher)

	i := 0
	callCount := make(chan bool, 200)
	f := func() int {
		i++
		callCount <- true
		pers.Return(mockMatcher.MatchOutput, nil)
		return i
	}

	if err := m.Match(f); err != nil {
		t.Fatal("expected err to be nil")
	}

	if len(callCount) != 1 {
		t.Errorf("expected callCount (%d) to equal %d", len(callCount), 1)
	}
}

func TestEventuallyUsesGivenProperties(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[func() int, int](mockMatcher, matcher.EventuallyTimes(100), matcher.EventuallyInterval(10*time.Microsecond))

	callCount := 0
	f := func() int {
		callCount++
		pers.Return(mockMatcher.MatchOutput, fmt.Errorf("still wrong"))
		return 99
	}
	m.Match(f)

	if callCount != 100 {
		t.Errorf("expected callCount (%d) to equal %d", callCount, 100)
	}
}

func TestEventuallyPassesAlongChan(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[chan int, int](mockMatcher)
	pers.Return(mockMatcher.MatchOutput, nil)

	c := make(chan int)
	m.Match(c)

	if len(mockMatcher.MatchInput.Actual) != 1 {
		t.Errorf("expected Actual (len=%d) to have a len of %d", len(mockMatcher.MatchInput.Actual), 1)
	}

	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 0) {
		t.Errorf("expected %v to equal %v", actual, 0)
	}
}

func TestEventuallyUsesChildMatcherErr(t *testing.T) {
	t.Parallel()
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Eventually[func() int, int](mockMatcher)

	childErr := fmt.Errorf("some-message")
	f := func() int {
		pers.Return(mockMatcher.MatchOutput, childErr)
		return 99
	}

	err := m.Match(f)
	if err == nil {
		t.Fatal("expected err to not be nil")
	}

	if !errors.Is(err, childErr) {
		t.Errorf("expected errors.Is(`%v`, `%v`) to return true", err, childErr)
	}
}
