package matchers_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestOrFailureUsesEachMatcher(t *testing.T) {
	t.Parallel()
	mockMatcherA := newMockMatcher()
	mockMatcherB := newMockMatcher()
	mockMatcherC := newMockMatcher()

	mockMatcherA.MatchOutput.ResultValue <- 1
	mockMatcherB.MatchOutput.ResultValue <- 2
	mockMatcherC.MatchOutput.ResultValue <- 3
	mockMatcherA.MatchOutput.Err <- fmt.Errorf("some-error")
	mockMatcherB.MatchOutput.Err <- fmt.Errorf("some-error")
	mockMatcherC.MatchOutput.Err <- fmt.Errorf("some-error")

	m := matchers.Or(mockMatcherA, mockMatcherB, mockMatcherC)

	_, err := m.Match(101)

	if err == nil {
		t.Error("expected err to not be nil")
	}

	actual := <-mockMatcherA.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("expecte %v to equal %v", actual, 101)
	}

	actual = <-mockMatcherB.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("expecte %v to equal %v", actual, 101)
	}

	actual = <-mockMatcherC.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("expecte %v to equal %v", actual, 101)
	}
}

func TestOrStopsOnSuccess(t *testing.T) {
	t.Parallel()
	mockMatcherA := newMockMatcher()
	mockMatcherB := newMockMatcher()
	mockMatcherC := newMockMatcher()

	mockMatcherA.MatchOutput.ResultValue <- 1
	mockMatcherB.MatchOutput.ResultValue <- 2
	mockMatcherC.MatchOutput.ResultValue <- 3
	mockMatcherA.MatchOutput.Err <- fmt.Errorf("some-error")
	close(mockMatcherB.MatchOutput.Err)
	close(mockMatcherC.MatchOutput.Err)

	m := matchers.Or(mockMatcherA, mockMatcherB, mockMatcherC)

	v, err := m.Match(101)

	if err != nil {
		t.Error("expected err to be nil")
	}

	if len(mockMatcherA.MatchCalled) != 1 {
		t.Errorf("expected Match() to be invoked 1 time")
	}

	if len(mockMatcherB.MatchCalled) != 1 {
		t.Errorf("expected Match() to be invoked 1 time")
	}

	if len(mockMatcherC.MatchCalled) != 0 {
		t.Errorf("expected Match() to be invoked 0 times")
	}

	if !reflect.DeepEqual(v, 101) {
		t.Errorf("expecte %v to equal %v", v, 101)
	}
}
