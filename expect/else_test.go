package expect_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/apoydence/onpar/expect"
)

func TestElseFailNowNotInvokedOnSuccess(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	close(mockMatcher.MatchOutput.Err)

	Expect(mockT, 101).To(mockMatcher).Else.FailNow()

	if len(mockT.FatalCalled) != 0 {
		t.Fatal("Expected Fatal() to not be invoked")
	}
}

func TestElseFailNowKillsTestOnFailure(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	Expect(mockT, 101).To(mockMatcher).Else.FailNow()

	if len(mockT.FatalInput.Arg0) != 1 {
		t.Fatal("Expected Fatal() to be invoked")
	}

	msg := <-mockT.FatalInput.Arg0
	if !reflect.DeepEqual(msg, []interface{}{"some-error"}) {
		t.Errorf("Expected %s to equal %s", msg, "some-error")
	}
}

func TestElseToNotInvokedOnSuccess(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	close(mockMatcher.MatchOutput.Err)

	Expect(mockT, 101).To(mockMatcher).Else.To(mockMatcher)

	if len(mockMatcher.MatchCalled) != 1 {
		t.Fatal("Expected Match to not be invoked via Else")
	}

}

func TestElseToInvokedOnFailure(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")
	mockMatcher.MatchOutput.Err <- nil

	Expect(mockT, 101).To(mockMatcher).Else.To(mockMatcher)
	if len(mockMatcher.MatchInput.Actual) != 2 {
		t.Fatal("Expected Match() to be invoked via Else")
	}

	<-mockMatcher.MatchInput.Actual
	actual := <-mockMatcher.MatchInput.Actual
	if !reflect.DeepEqual(actual, 101) {
		t.Errorf("Expected %v to equal %v", actual, 101)
	}
}
