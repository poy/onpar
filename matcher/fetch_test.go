package matcher_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestFetchMatcher(t *testing.T) {
	t.Parallel()

	var i int
	m := matcher.FetchTo(&i)

	if err := m.Match(99); err != nil {
		t.Fatal("expected err to be nil")
	}

	if !reflect.DeepEqual(i, 99) {
		t.Fatalf("expected %v to equal %v", i, 99)
	}
}
