//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nelsam/hel/v2/pers"
	. "github.com/poy/onpar/expect"
)

type diffMatcher struct {
	mockToMatcher
	mockDiffMatcher
}

func TestToRespectsDifferMatchers(t *testing.T) {
	mockT := newMockT()
	d := newMockDiffMatcher()
	m := newMockToMatcher()
	mockMatcher := &diffMatcher{
		mockDiffMatcher: *d,
		mockToMatcher:   *m,
	}
	pers.Return(mockMatcher.MatchOutput, nil, nil)
	mockDiffer := newMockDiffer()

	f := New(mockT, WithDiffer(mockDiffer))
	f(101).To(mockMatcher)

	Expect(t, mockMatcher).To(pers.HaveMethodExecuted("UseDiffer", pers.WithArgs(mockDiffer)))
}

func TestToPassesActualToMatcher(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockToMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	close(mockMatcher.MatchOutput.Err)

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
	mockT := newMockT()
	mockMatcher := newMockToMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	Expect(mockT, 101).To(mockMatcher)

	if len(mockT.FatalfCalled) != 1 {
		t.Error("expected Fatalf to be invoked 1 time")
	}
}
