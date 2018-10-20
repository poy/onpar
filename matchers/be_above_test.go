package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestAbove(t *testing.T) {
	t.Parallel()

	m := matchers.BeAbove(101)

	_, err := m.Match(101.0)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(99.0)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(int(99))
	if err == nil {
		t.Error("expected err to not be nil")
	}

	v, err := m.Match(103.0)
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(v, 103.0) {
		t.Errorf("expected %v to equal %v", v, 103.0)
	}

	v, err = m.Match(int(103))
	if err != nil {
		t.Error("expected err to be nil")
	}

	_, err = m.Match("invalid")
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
