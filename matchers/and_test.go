package matchers_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matchers"
)

func TestAndSuccessUsesEachMatcher(t *testing.T) {
	mockMatcherA := newMockMatcher(t, time.Second)
	mockMatcherB := newMockMatcher(t, time.Second)
	mockMatcherC := newMockMatcher(t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, 1, nil)
	pers.Return(mockMatcherB.MatchOutput, 2, nil)
	pers.Return(mockMatcherC.MatchOutput, 3, nil)

	m := matchers.And(mockMatcherA, mockMatcherB, mockMatcherC)

	v, err := m.Match(101)

	if err != nil {
		t.Error("expected err to be nil")
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

	if !reflect.DeepEqual(v, 101) {
		t.Errorf("expecte %v to equal %v", v, 101)
	}
}

func TestAndStopsOnFailure(t *testing.T) {
	mockMatcherA := newMockMatcher(t, time.Second)
	mockMatcherB := newMockMatcher(t, time.Second)
	mockMatcherC := newMockMatcher(t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, 1, nil)
	pers.Return(mockMatcherB.MatchOutput, 2, fmt.Errorf("some-error"))
	pers.Return(mockMatcherC.MatchOutput, 3, nil)

	m := matchers.And(mockMatcherA, mockMatcherB, mockMatcherC)

	_, err := m.Match(101)

	if err == nil {
		t.Error("expected err to not be nil")
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
