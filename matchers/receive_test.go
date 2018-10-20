package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestReceiveSucceedsForABufferedChannel(t *testing.T) {
	t.Parallel()
	c := make(chan bool, 1)
	m := matchers.Receive()

	_, err := m.Match(c)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	c <- true
	result, err := m.Match(c)
	if err != nil {
		t.Error("expected err to be nil")
	}

	if !reflect.DeepEqual(result, true) {
		t.Errorf("expected %v to equal %v", result, true)
	}
}

func TestReceiveFailsNotReadableChan(t *testing.T) {
	t.Parallel()
	c := make(chan int, 10)
	m := matchers.Receive()

	_, err := m.Match(101)
	if err == nil {
		t.Error("expected err to not be nil")
	}

	_, err = m.Match(chan<- int(c))
	if err == nil {
		t.Error("expected err to not be nil")
	}
}

func TestReceiveFailsForClosedChannel(t *testing.T) {
	t.Parallel()
	c := make(chan bool, 1)
	m := matchers.Receive()
	c <- true

	_, err := m.Match(c)
	if err != nil {
		t.Error("expected err to be nil")
	}

	close(c)
	_, err = m.Match(c)
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
