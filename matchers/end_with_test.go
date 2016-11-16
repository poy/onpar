package matchers_test

import (
	"reflect"
	"testing"

	"github.com/apoydence/onpar/matchers"
)

func TestEndsWith(t *testing.T) {
	t.Parallel()

	m := matchers.EndsWith("foo")

	_, err := m.Match("bar")
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err := m.Match("barfoo")
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(v, "barfoo") {
		t.Errorf("expected %v to equal %v", v, "barfoo")
	}

	_, err = m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
