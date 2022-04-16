package matcher

import (
	"fmt"
	"time"
)

type recvPrefs struct {
	timeout time.Duration
}

// ReceiveOpt is an option that can be passed to the
// ReceiveMatcher constructor.
type ReceiveOpt func(recvPrefs) recvPrefs

// ReceiveWait is an option that makes the ReceiveMatcher
// wait for values for the provided duration before
// deciding that the channel failed to receive.
func ReceiveWait(t time.Duration) ReceiveOpt {
	return func(p recvPrefs) recvPrefs {
		p.timeout = t
		return p
	}
}

// ReceiveMatcher only accepts a readable channel. It will error for anything else.
// It will attempt to receive from the channel but will not block.
// It fails if the channel is closed.
type ReceiveMatcher[T ~chan U, U any] struct {
	sub     Matcher[U]
	timeout time.Duration
}

// Receive will return a ReceiveMatcher
func Receive[T ~chan U, U any](sub Matcher[U], opts ...ReceiveOpt) ReceiveMatcher[T, U] {
	var p recvPrefs
	for _, opt := range opts {
		p = opt(p)
	}
	return ReceiveMatcher[T, U]{
		sub:     sub,
		timeout: p.timeout,
	}
}

// Match receives a value from actual, checking it using the sub matcher.
func (m ReceiveMatcher[T, U]) Match(actual T) error {
	if m.timeout == 0 {
		select {
		case v, ok := <-actual:
			return m.subMatch(v, ok)
		default:
			return fmt.Errorf("expected to receive a value, but no value was available")
		}
	}

	select {
	case v, ok := <-actual:
		return m.subMatch(v, ok)
	case <-time.After(m.timeout):
		return fmt.Errorf("expected to receive a value, but timed out after %v", m.timeout)
	}
}

func (m ReceiveMatcher[T, U]) subMatch(actual U, open bool) error {
	if !open {
		return fmt.Errorf("expected channel to not be closed")
	}
	if err := m.sub.Match(actual); err != nil {
		return fmt.Errorf("received value failed sub-matcher: %v", err)
	}
	return nil
}
