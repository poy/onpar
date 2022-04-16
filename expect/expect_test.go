//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	. "github.com/poy/onpar/v2/expect"
	"github.com/poy/onpar/v2/matcher"
)

type diffMatcher[T any] struct {
	mockMatcher[T]
	mockDiffMatcher
}

func TestToRespectsDifferMatchers(t *testing.T) {
	mockT := newMockT(t, time.Second)
	d := newMockDiffMatcher(t, time.Second)
	m := newMockMatcher[int](t, time.Second)
	mMatcher := &diffMatcher[int]{
		mockDiffMatcher: *d,
		mockMatcher:     *m,
	}
	pers.Return(mMatcher.MatchOutput, nil)
	mockDiffer := newMockDiffer(t, time.Second)

	Expect(mockT, 101, WithDiffer(mockDiffer)).To(mMatcher)

	Expect(t, mMatcher).To(matcher.Reflect[*diffMatcher[int]](pers.HaveMethodExecuted("UseDiffer", pers.WithArgs(mockDiffer))))
}

func TestToPassesActualMatcher(t *testing.T) {
	mockT := newMockT(t, time.Second)
	mockMatcher := newMockMatcher[int](t, time.Second)
	pers.Return(mockMatcher.MatchOutput, nil)

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
	mockMatcher := newMockMatcher[int](t, time.Second)
	pers.Return(mockMatcher.MatchOutput, fmt.Errorf("some-error"))

	Expect(mockT, 101).To(mockMatcher)

	if len(mockT.FatalfCalled) != 1 {
		t.Error("expected Fatalf to be invoked 1 time")
	}
}
