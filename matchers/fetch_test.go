package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestFetchMatcher(t *testing.T) {
	t.Parallel()

	var i int
	m := matchers.Fetch(&i)

	v, err := m.Match(99)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	if !reflect.DeepEqual(v, 99) {
		t.Fatalf("expected %v to equal %v", v, 99)
	}

	if !reflect.DeepEqual(i, 99) {
		t.Fatalf("expected %v to equal %v", i, 99)
	}

	_, err = m.Match("invalid")
	if err == nil {
		t.Fatal("expected err to not equal nil")
	}

	m = matchers.Fetch(i)
	_, err = m.Match(99)
	if err == nil {
		t.Fatal("expected err to not equal nil")
	}
}
