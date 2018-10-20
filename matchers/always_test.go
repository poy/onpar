package matchers_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/poy/onpar/matchers"
)

func TestFailsWhenMatcherFails(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.Always(matcher)
	matcher.MatchOutput.ResultValue <- nil
	matcher.MatchOutput.Err <- fmt.Errorf("some-error")

	callCount := make(chan bool, 200)
	f := func() int {
		callCount <- true
		return 99
	}

	_, err := m.Match(f)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	if len(callCount) != 1 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 1)
	}
}

func TestPolls10TimesForSuccess(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.Always(matcher)

	i := 0
	callCount := make(chan bool, 200)
	f := func() int {
		i++
		callCount <- true
		matcher.MatchOutput.ResultValue <- i
		matcher.MatchOutput.Err <- nil
		return i
	}

	v, err := m.Match(f)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	if !reflect.DeepEqual(v, 10) {
		t.Errorf("expected %v to equal %v", v, 10)
	}

	if len(callCount) != 10 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 10)
	}
}

func TestPollsEach10msForSuccess(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.Always(matcher)

	var ts []int64
	f := func() int {
		ts = append(ts, time.Now().UnixNano())
		matcher.MatchOutput.ResultValue <- nil
		matcher.MatchOutput.Err <- nil
		return 101
	}
	m.Match(f)

	for i := 0; i < len(ts)-1; i++ {
		if ts[i+1]-ts[i] < int64(10*time.Millisecond) || ts[i+1]-ts[i] > int64(30*time.Millisecond) {
			t.Fatalf("expected %d to be within 10ms and 30ms", ts[i+1]-ts[i])
		}
	}
}
