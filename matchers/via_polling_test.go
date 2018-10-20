package matchers_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/poy/onpar/matchers"
)

func TestViaPollingFailsPolls100Times(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	callCount := make(chan bool, 200)
	f := func() int {
		callCount <- true
		matcher.MatchOutput.ResultValue <- nil
		matcher.MatchOutput.Err <- fmt.Errorf("still wrong")
		return 99
	}
	m.Match(f)

	if len(callCount) != 100 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 100)
	}
}

func TestViaPollingPollsEvery10ms(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	var ts []int64
	f := func() int {
		ts = append(ts, time.Now().UnixNano())
		matcher.MatchOutput.ResultValue <- nil
		matcher.MatchOutput.Err <- fmt.Errorf("still wrong")
		return 99
	}
	m.Match(f)

	for i := 0; i < len(ts)-1; i++ {
		if ts[i+1]-ts[i] < int64(10*time.Millisecond) || ts[i+1]-ts[i] > int64(30*time.Millisecond) {
			t.Fatalf("expected %d to be within 10ms and 30ms", ts[i+1]-ts[i])
		}
	}
}

func TestViaPollingStopsAfterSuccess(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

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

	if !reflect.DeepEqual(v, 1) {
		t.Errorf("expected %v to equal %v", v, 1)
	}

	if len(callCount) != 1 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 1)
	}
}

func TestViaPollingUsesGivenProperties(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPollingMatcher{
		Matcher:  matcher,
		Duration: time.Millisecond,
		Interval: 100 * time.Microsecond,
	}

	callCount := make(chan bool, 200)
	f := func() int {
		callCount <- true
		matcher.MatchOutput.ResultValue <- nil
		matcher.MatchOutput.Err <- fmt.Errorf("still wrong")
		return 99
	}
	m.Match(f)

	if len(callCount) != 10 {
		t.Errorf("expected callCount (len=%d) to have a len of %d", len(callCount), 10)
	}
}

func TestViaPollingPassesAlongChan(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)
	matcher.MatchOutput.ResultValue <- nil
	matcher.MatchOutput.Err <- nil

	c := make(chan int)
	m.Match(c)

	if len(matcher.MatchInput.Actual) != 1 {
		t.Errorf("expected Actual (len=%d) to have a len of %d", len(matcher.MatchInput.Actual), 1)
	}

	actual := <-matcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, c) {
		t.Errorf("expected %v to equal %v", actual, c)
	}
}

func TestViaPollingFailsForNonChanOrFunc(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	_, err := m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}

func TestViaPollingFailsForFuncWithArgs(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	_, err := m.Match(func(int) int { return 101 })
	if err == nil {
		t.Error("expected err to not be nil")
	}
}

func TestViaPollingFailsForFuncWithWrongReturns(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	_, err := m.Match(func() (int, int) { return 101, 103 })
	if err == nil {
		t.Error("expected err to not be nil")
	}
}

func TestViaPollingFailsForSendOnlyChan(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	_, err := m.Match(make(chan<- int))
	if err == nil {
		t.Error("expected err to not be nil")
	}
}

func TestViaPollingUsesChildMatcherErr(t *testing.T) {
	t.Parallel()
	matcher := newMockMatcher()
	m := matchers.ViaPolling(matcher)

	f := func() int {
		matcher.MatchOutput.Err <- fmt.Errorf("some-message")
		matcher.MatchOutput.ResultValue <- nil
		return 99
	}
	_, err := m.Match(f)

	if err == nil {
		t.Fatal("expected err to not be nil")
	}

	if err.Error() != "some-message" {
		t.Errorf("expected err to have message: %s", "some-message")
	}
}
