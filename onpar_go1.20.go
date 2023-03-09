//go:build go1.20

package onpar

import "testing"

func init() {
	panic("onpar: go 1.20.x introduced a breaking change that broke onpar v2, and we could not solve it without a breaking change of our own. Please upgrade to onpar v3.")
}

type prefs struct{}

// Opt v2 is broken on go 1.20! Please use onpar v3!
type Opt func(prefs) prefs

// Onpar v2 is broken on go 1.20! Please use onpar v3!
//
// This type is provided to avoid introducing difficult-to-understand
// compilation errors in go 1.20. Importing onpar v2 on go 1.20 will just panic.
type Onpar[T, U any] struct{}

// New returns a useless Onpar, because onpar v2 is broken on go 1.20!
func New(t *testing.T, opts ...Opt) *Onpar[*testing.T, *testing.T] {
	return &Onpar[*testing.T, *testing.T]{}
}

// BeforeEach v2 is broken on go 1.20! Please use onpar v3!
func BeforeEach[T, U, V any](_ *Onpar[T, U], _ func(U) V) *Onpar[U, V] {
	return &Onpar[U, V]{}
}

// Spec v2 is broken on go 1.20! Please use onpar v3!
func (o *Onpar[T, U]) Spec(_ string, _ func(U)) {
}

// SerialSpec v2 is broken on go 1.20! Please use onpar v3!
func (o *Onpar[T, U]) SerialSpec(name string, f func(U)) {
}

// Group v2 is broken on go 1.20! Please use onpar v3!
func (o *Onpar[T, U]) Group(_ string, _ func()) {
}

// AfterEach v2 is broken on go 1.20! Please use onpar v3!
func (o *Onpar[T, U]) AfterEach(_ func(U)) {
}
