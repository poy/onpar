//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	. "github.com/poy/onpar/v2/expect"
)

type diffMatcher struct {
	mockToMatcher
	mockDiffMatcher
}

func TestToRespectsDifferMatchers(t *testing.T) {
	mockT := newMockT(t, time.Second)
	d := newMockDiffMatcher(t, time.Second)
	m := newMockToMatcher(t, time.Second)
	mockMatcher := &diffMatcher{
		mockDiffMatcher: *d,
		mockToMatcher:   *m,
	}
	pers.Return(mockMatcher.MatchOutput, nil, nil)
	mockDiffer := newMockDiffer(t, time.Second)

	f := New(mockT, WithDiffer(mockDiffer))
	f(101).To(mockMatcher)

	Expect(t, mockMatcher).To(pers.HaveMethodExecuted("UseDiffer", pers.WithArgs(mockDiffer)))
}

func TestToPassesActualToMatcher(t *testing.T) {
	mockT := newMockT(t, time.Second)
	mockMatcher := newMockToMatcher(t, time.Second)
	pers.Return(mockMatcher.MatchOutput, nil, nil)

	Expect(mockT, 101).To(mockMatcher)

	select {
	case actual := <-mockMatcher.MatchInput.Actual:
		if !reflect.DeepEqual(actual, 101) {
			t.Errorf("Expected %v to equal %v", actual, 101)
		}
	default:
		t.Errorf("Expected Match() to be invoked")
	}
}

func TestToErrorsIfMatcherFails(t *testing.T) {
	mockT := newMockT(t, time.Second)
	mockMatcher := newMockToMatcher(t, time.Second)
	pers.Return(mockMatcher.MatchOutput, nil, fmt.Errorf("some-error"))

	Expect(mockT, 101).To(mockMatcher)

	if len(mockT.FatalfCalled) != 1 {
		t.Error("expected Fatalf to be invoked 1 time")
	}
}
