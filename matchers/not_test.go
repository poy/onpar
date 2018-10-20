//go:generate hel

package matchers_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestNot(t *testing.T) {
	mockMatcher := newMockMatcher()
	m := matchers.Not(mockMatcher)

	mockMatcher.MatchOutput.ResultValue <- 103
	mockMatcher.MatchOutput.Err <- nil

	_, err := m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	if len(mockMatcher.MatchInput.Actual) != 1 {
		t.Fatal("expected child macther Match() to be invoked once")
	}

	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Fatalf("expected %v does not equal %v", actual, 101)
	}

	mockMatcher.MatchOutput.ResultValue <- 103
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	v, err := m.Match(101)
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(v, 103) {
		t.Errorf("expected %v to equal %v", v, 103)
	}
}
