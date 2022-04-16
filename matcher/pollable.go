package matcher

import (
	"fmt"
)

// Pollable is any type which may be used in polling matchers.
type Pollable[T any] interface {
	~func() T | ~chan T | ~<-chan T
}

func fetchFunc[T Pollable[U], U any](actual T) func() U {
	// Since we can't do branching logic depending on actual type in generics,
	// yet, we need to convert to an interface type so that we can type-assert
	// and type-switch.
	inter := interface{}(actual)
	if f, ok := inter.(func() U); ok {
		return f
	}
	var rcv <-chan U
	switch v := inter.(type) {
	case chan U:
		rcv = v
	case <-chan U:
		rcv = v
	default:
		panic(fmt.Errorf("unhandled Pollable type from type union: %T", actual))
	}
	return func() U {
		select {
		case e := <-rcv:
			return e
		default:
			var empty U
			return empty
		}
	}
}
