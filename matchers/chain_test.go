package matchers_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matchers"
)

func TestChainPassesResultOnSuccess(t *testing.T) {
	mockMatcherA := newMockMatcher(t, time.Second)
	mockMatcherB := newMockMatcher(t, time.Second)
	mockMatcherC := newMockMatcher(t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, 1, nil)
	pers.Return(mockMatcherB.MatchOutput, 2, nil)
	pers.Return(mockMatcherC.MatchOutput, 3, nil)

	m := matchers.Chain(mockMatcherA, mockMatcherB, mockMatcherC)
	result, err := m.Match(101)

	if len(mockMatcherB.MatchInput.Actual) != 1 {
		t.Fatalf("expected Match() to be invoked")
	}

	if len(mockMatcherC.MatchInput.Actual) != 1 {
		t.Fatalf("expected Match() to be invoked")
	}

	bActual := <-mockMatcherB.MatchInput.Actual
	if !reflect.DeepEqual(bActual, 1) {
		t.Errorf("expected %v to equal %v", bActual, 1)
	}

	cActual := <-mockMatcherC.MatchInput.Actual
	if !reflect.DeepEqual(cActual, 2) {
		t.Errorf("expected %v to equal %v", cActual, 2)
	}

	if err != nil {
		t.Errorf("expected err to be nil")
	}

	if !reflect.DeepEqual(result, 3) {
		t.Errorf("expected %v to equal %v", result, 3)
	}
}

func TestChainStopsOnFailuire(t *testing.T) {
	mockMatcherA := newMockMatcher(t, time.Second)
	mockMatcherB := newMockMatcher(t, time.Second)
	mockMatcherC := newMockMatcher(t, time.Second)

	pers.Return(mockMatcherA.MatchOutput, 1, nil)
	pers.Return(mockMatcherB.MatchOutput, 2, fmt.Errorf("some-error"))
	pers.Return(mockMatcherC.MatchOutput, 3, nil)

	m := matchers.Chain(mockMatcherA, mockMatcherB, mockMatcherC)
	_, err := m.Match(101)

	if len(mockMatcherB.MatchInput.Actual) != 1 {
		t.Fatalf("expected Match() to be invoked")
	}

	if len(mockMatcherC.MatchInput.Actual) != 0 {
		t.Fatalf("expected Match() to not be invoked")
	}

	bActual := <-mockMatcherB.MatchInput.Actual
	if !reflect.DeepEqual(bActual, 1) {
		t.Errorf("expected %v to equal %v", bActual, 1)
	}

	if err == nil {
		t.Errorf("expected err to not be nil")
	}
}
