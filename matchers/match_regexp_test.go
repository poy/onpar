package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestRegexp(t *testing.T) {
	t.Parallel()

	m := matchers.MatchRegexp("^foo")

	_, err := m.Match("barfoo")
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err := m.Match("foobar")
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(v, "foobar") {
		t.Errorf("expected %v to equal %v", v, "foobar")
	}

	_, err = m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
