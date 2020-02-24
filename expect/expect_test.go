//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nelsam/hel/pers"
	"github.com/poy/onpar/diff"
	. "github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
)

type diffMatcher struct {
	mockDiffer
	mockMatcher
}

func TestToRespectsDifferMatchers(t *testing.T) {
	mockT := newMockT()
	m := newMockMatcher()
	d := newMockDiffer()
	mockMatcher := &diffMatcher{
		mockDiffer:  *d,
		mockMatcher: *m,
	}
	done, err := pers.ConsistentlyReturn(mockMatcher.MatchOutput, nil, nil)
	if err != nil {
		t.Fatalf("faild to consistenly return")
	}
	defer done()

	actualOpt := diff.Actual(diff.WithFGColor(diff.Red))
	expectedOpt := diff.Expected(diff.WithFGColor(diff.Green))

	f := New(mockT, WithDiffOpts(actualOpt, expectedOpt))
	f(101).To(mockMatcher)

	select {
	case <-mockMatcher.UseDiffOptsCalled:
		opts := <-mockMatcher.UseDiffOptsInput.Opts
		Expect(t, opts).To(matchers.HaveLen(2))
		Expect(t, fmt.Sprintf("%p", opts[0])).To(matchers.Equal(fmt.Sprintf("%p", actualOpt)))
		Expect(t, fmt.Sprintf("%p", opts[1])).To(matchers.Equal(fmt.Sprintf("%p", expectedOpt)))
	default:
		t.Fatalf("Expected matcher.UseDiffOpts to be called, but it was not.")
	}
}

func TestToPassesActualToMatcher(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
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
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	Expect(mockT, 101).To(mockMatcher)

	if len(mockT.FatalfCalled) != 1 {
		t.Error("expected Fatalf to be invoked 1 time")
	}
}
