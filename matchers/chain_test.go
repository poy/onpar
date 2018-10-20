package matchers_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestChainPassesResultOnSuccess(t *testing.T) {
	mockMatcherA := newMockMatcher()
	mockMatcherB := newMockMatcher()
	mockMatcherC := newMockMatcher()
	close(mockMatcherA.MatchOutput.Err)
	close(mockMatcherB.MatchOutput.Err)
	close(mockMatcherC.MatchOutput.Err)
	mockMatcherA.MatchOutput.ResultValue <- 1
	mockMatcherB.MatchOutput.ResultValue <- 2
	mockMatcherC.MatchOutput.ResultValue <- 3

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
	mockMatcherA := newMockMatcher()
	mockMatcherB := newMockMatcher()
	mockMatcherC := newMockMatcher()
	close(mockMatcherA.MatchOutput.Err)
	mockMatcherB.MatchOutput.Err <- fmt.Errorf("some-error")
	close(mockMatcherC.MatchOutput.Err)
	mockMatcherA.MatchOutput.ResultValue <- 1
	mockMatcherB.MatchOutput.ResultValue <- 2
	mockMatcherC.MatchOutput.ResultValue <- 3

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
