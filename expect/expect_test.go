//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/apoydence/onpar/expect"
)

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

	select {
	case msg := <-mockT.ErrorfInput.Format:
		if msg != "some-error" {
			t.Errorf("Expected %v to equal %v", msg, "some-error")
		}
	default:
		t.Errorf("Expected Errorf() to be invoked")
	}
}

func TestToReturnsChainedToOnSuccess(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	close(mockMatcher.MatchOutput.Err)

	Expect(mockT, 101).To(mockMatcher).And.To(mockMatcher)

	if len(mockMatcher.MatchInput.Actual) != 2 {
		t.Fatal("Expected Match() to be invoked 2 times")
	}

	<-mockMatcher.MatchInput.Actual
	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("Expected %v to equal %v", actual, 101)
	}
}

func TestToDoesNotInvokeAndMatcherOnFailure(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	Expect(mockT, 101).To(mockMatcher).And.To(mockMatcher)

	if len(mockMatcher.MatchInput.Actual) != 1 {
		t.Fatal("Expected Match() to be invoked 1 time")
	}
}

func TestPassesNewValuesAlong(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	mockMatcher.MatchOutput.ResultValue <- 99
	mockMatcher.MatchOutput.ResultValue <- 103
	close(mockMatcher.MatchOutput.Err)

	Expect(mockT, 101).To(mockMatcher).AndForThat.To(mockMatcher)

	if len(mockMatcher.MatchInput.Actual) != 2 {
		t.Fatal("Expected Match() to be invoked 2 times")
	}

	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("Expected %v to equal %v", actual, 101)
	}

	actual = <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 99) {
		t.Errorf("Expected %v to equal %v", actual, 99)
	}
}
