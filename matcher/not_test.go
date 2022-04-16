//go:generate hel

package matcher_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/matcher"
)

func TestNot(t *testing.T) {
	mockMatcher := newMockMatcher[int](t, time.Second)
	m := matcher.Not[int](mockMatcher)

	pers.Return(mockMatcher.MatchOutput, nil)

	if err := m.Match(101); err == nil {
		t.Error("expected err to not be nil")
	}

	if len(mockMatcher.MatchInput.Actual) != 1 {
		t.Fatal("expected child macther Match() to be invoked once")
	}

	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Fatalf("expected %v does not equal %v", actual, 101)
	}

	pers.Return(mockMatcher.MatchOutput, fmt.Errorf("some-error"))

	if err := m.Match(101); err != nil {
		t.Error("expected err to be nil")
	}
}
