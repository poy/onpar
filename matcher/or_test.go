package matcher_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matcher"
)

func TestOrFailureUsesEachMatcher(t *testing.T) {
	t.Parallel()
	mockMatcherA := newMockMatcher[int](t, time.Second)
	mockMatcherB := newMockMatcher[int](t, time.Second)
	mockMatcherC := newMockMatcher[int](t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, fmt.Errorf("some-error"))
	pers.Return(mockMatcherB.MatchOutput, fmt.Errorf("some-error"))
	pers.Return(mockMatcherC.MatchOutput, fmt.Errorf("some-error"))

	m := matcher.Or[int](mockMatcherA, mockMatcherB, mockMatcherC)

	if err := m.Match(101); err == nil {
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
	mockMatcherA := newMockMatcher[int](t, time.Second)
	mockMatcherB := newMockMatcher[int](t, time.Second)
	mockMatcherC := newMockMatcher[int](t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, fmt.Errorf("some-error"))
	pers.Return(mockMatcherB.MatchOutput, nil)
	pers.Return(mockMatcherC.MatchOutput, nil)

	m := matcher.Or[int](mockMatcherA, mockMatcherB, mockMatcherC)

	if err := m.Match(101); err != nil {
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
}
