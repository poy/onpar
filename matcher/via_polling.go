package matcher

import (
	"fmt"
	"time"
)

type eventuallyPrefs struct {
	interval time.Duration
	times    int
}

// EventuallyOpt is an option to be passed to the Eventually function to make
// optional changes to behavior.
type EventuallyOpt func(eventuallyPrefs) eventuallyPrefs

// EventuallyTimes sets the number of times that the EventuallyMatcher will poll
// the child matcher for an answer before giving up. The duration of the
// EventuallyMatcher will end up being roughly this number times the
// EventuallyInterval unless the child matcher takes longer than the interval.
//
// The default is 10.
func EventuallyTimes(times int) EventuallyOpt {
	return func(o eventuallyPrefs) eventuallyPrefs {
		o.times = times
		return o
	}
}

// EventuallyInteraval sets the interaval that the EventuallyMatcher will poll
// the child matcher at. 1s means it will be polled once per second.
//
// The default is 10ms.
func EventuallyInterval(d time.Duration) EventuallyOpt {
	return func(o eventuallyPrefs) eventuallyPrefs {
		o.interval = d
		return o
	}
}

// EventuallyMatcher matches by polling the child matcher until
// it returns a success. It will return success the first time
// the child matcher returns a success. If the child matcher
// never returns a nil, then it will return the last error.
type EventuallyMatcher[T Pollable[U], U any] struct {
	matcher  Matcher[U]
	interval time.Duration
	times    int
}

// Eventually returns the default EventuallyMatcher.
func Eventually[T Pollable[U], U any](m Matcher[U], opts ...EventuallyOpt) EventuallyMatcher[T, U] {
	p := eventuallyPrefs{
		interval: 10 * time.Millisecond,
		times:    10,
	}
	for _, o := range opts {
		p = o(p)
	}
	return EventuallyMatcher[T, U]{
		matcher:  m,
		interval: p.interval,
		times:    p.times,
	}
}

// Match takes a Pollable type and polls it on an interval, returning success
// when it passes the child matcher. If no success is returned for the entire
// duration that the EventuallyMatcher is configured with, a timeout error will
// be returned.
func (m EventuallyMatcher[T, U]) Match(actual T) error {
	f := fetchFunc[T, U](actual)

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	times := 0
	var lastErr error
	for range ticker.C {
		if times >= m.times {
			return fmt.Errorf("timed out waiting for sub-matcher to pass - last failure: %w", lastErr)
		}
		lastErr = m.matcher.Match(f())
		if lastErr == nil {
			return nil
		}
		times++
	}
	panic("this should never happen (the ticker does not have a way to stop on its own)")
}
