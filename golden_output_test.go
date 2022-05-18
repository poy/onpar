//go:build goldenoutput

package onpar_test

import (
	"testing"

	"github.com/poy/onpar/v2"
)

func TestNestedStructure(t *testing.T) {
	// This test is used to generate verbose output to ensure that nested
	// structure creates the expected nested calls to t.Run. We have some more
	// unit-like tests that assert on these things too, but they can only really
	// test that the test names use the expected paths. They cannot prove that
	// each group has its own t.Run.
	top := onpar.New(t)

	top.Spec("foo", func(*testing.T) {})

	b := onpar.BeforeEach(top, func(t *testing.T) *testing.T { return t })

	b.Spec("bar", func(*testing.T) {})

	b.Group("baz", func() {
		b := onpar.BeforeEach(b, func(t *testing.T) string { return "foo" })

		b.Spec("foo", func(string) {})

		b.Group("bar", func() {
			b.Spec("foo", func(string) {})
		})

		b.Group("", func() {
			b.Spec("baz", func(string) {})
		})
	})
}
