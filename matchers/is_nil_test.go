package matchers_test

import (
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestBeNil(t *testing.T) {
	t.Parallel()

	m := matchers.BeNil()

	_, err := m.Match(nil)
	if err != nil {
		t.Error("expected err to be nil")
	}

	v, err := m.Match(42)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	if v != 42 {
		t.Errorf("expected %v to equal %d", v, 42)
	}

	var x map[string]string
	_, err = m.Match(x)
	if err != nil {
		t.Error("expected err to be nil")
	}

	x = map[string]string{}

	_, err = m.Match(x)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
