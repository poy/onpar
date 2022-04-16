package matcher

import (
	"time"
)

type alwaysPrefs struct {
	interval time.Duration
	times    int
}

// AlwaysOpt is an option to be passed to the Always function to make optional
// changes to behavior.
type AlwaysOpt func(alwaysPrefs) alwaysPrefs

// AlwaysTimes sets the number of times that the AlwaysMatcher will poll the
// child matcher for an answer before giving up. The duration of the
// AlwaysMatcher will end up being roughly this number times the AlwaysInterval
// unless the child matcher takes longer than the interval.
//
// The default is 10.
func AlwaysDuration(times int) AlwaysOpt {
	return func(o alwaysPrefs) alwaysPrefs {
		o.times = times
		return o
	}
}

// AlwaysInteraval sets the interaval that the AlwaysMatcher will poll the child
// matcher at.  1s means it will be polled once per second.
//
// The default is 10ms
func AlwaysInteraval(d time.Duration) AlwaysOpt {
	return func(o alwaysPrefs) alwaysPrefs {
		o.interval = d
		return o
	}
}

// AlwaysMatcher matches by polling the child matcher until it returns an error.
// It will return an error the first time the child matcher returns an error. If
// the child matcher never returns an error, then it will return nil.
type AlwaysMatcher[T Pollable[U], U any] struct {
	matcher  Matcher[U]
	interval time.Duration
	times    int
}

// Always returns a default AlwaysMatcher.
func Always[T Pollable[U], U any](m Matcher[U], opts ...AlwaysOpt) AlwaysMatcher[T, U] {
	p := alwaysPrefs{
		interval: 10 * time.Millisecond,
		times:    10,
	}
	for _, o := range opts {
		p = o(p)
	}
	return AlwaysMatcher[T, U]{
		matcher:  m,
		interval: p.interval,
		times:    p.times,
	}
}

// Match takes a value that can change over time. Therefore, the only
// two valid options are a function with no arguments and a single return
// type, or a readable channel. Anything else will return an error.
//
// Channels will be polled by performing a select with a default on an interval.
//
// Child matchers will be passed the values read from polling either the
// function or the channel.
//
// TODO: decide if it's useful to poll channels. The receive matcher already
// supports a timeout.
func (m AlwaysMatcher[T, U]) Match(actual T) error {
	f := fetchFunc[T, U](actual)

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	times := 0
	for range ticker.C {
		if times >= m.times {
			return nil
		}
		if err := m.matcher.Match(f()); err != nil {
			return err
		}
		times++
	}
	panic("this should never happen (the ticker does not have a way to stop on its own)")
}
