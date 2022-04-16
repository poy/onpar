package matcher_test

import (
	"testing"
	"time"

	"github.com/poy/onpar/v2/matcher"
)

func TestReceiveSucceedsForABufferedChannel(t *testing.T) {
	t.Parallel()
	c := make(chan bool, 1)
	m := matcher.Receive[chan bool, bool](matcher.Anything[bool]())

	if err := m.Match(c); err == nil {
		t.Error("expected err to not be nil")
	}

	c <- true
	if err := m.Match(c); err != nil {
		t.Error("expected err to be nil")
	}
}

func TestReceiveWaitSucceedsForAnUnbufferedChannelInALoop(t *testing.T) {
	t.Parallel()
	c := make(chan bool)
	done := make(chan struct{})
	defer close(done)
	go func() {
		for {
			select {
			case c <- true:
			case <-done:
				return
			default:
				// Receive should still be able to get a value from c,
				// even in this buggy loop.  Onpar is for detecting
				// bugs, after all.
			}
		}
	}()

	m := matcher.Receive[chan bool, bool](matcher.Anything[bool](), matcher.ReceiveWait(500*time.Millisecond))
	if err := m.Match(c); err != nil {
		t.Errorf("expected err to be nil; got %s", err)
	}
}

func TestReceiveWaitSucceedsEventually(t *testing.T) {
	t.Parallel()
	c := make(chan bool)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c <- true
	}()

	m := matcher.Receive[chan bool, bool](matcher.Anything[bool](), matcher.ReceiveWait(500*time.Millisecond))
	if err := m.Match(c); err != nil {
		t.Errorf("expected err to be nil; got %s", err)
	}
}

func TestReceiveFailsForClosedChannel(t *testing.T) {
	t.Parallel()
	c := make(chan bool, 1)
	m := matcher.Receive[chan bool, bool](matcher.Anything[bool]())
	c <- true

	if err := m.Match(c); err != nil {
		t.Error("expected err to be nil")
	}

	close(c)
	if err := m.Match(c); err == nil {
		t.Error("expected err to not be nil")
	}
}
