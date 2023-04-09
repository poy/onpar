package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestBeClosed(t *testing.T) {
	t.Parallel()
	m := matchers.BeClosed()

	c := make(chan int)
	close(c)
	v, err := m.Match(c)
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(v, c) {
		t.Errorf("expected %v to equal %v", v, c)
	}

	c = make(chan int)
	_, err = m.Match(c)
	if err == nil {
		t.Errorf("expected err to not be nil")
	}

	_, err = m.Match(99)
	if err == nil {
		t.Errorf("expected err to not be nil")
	}
}
