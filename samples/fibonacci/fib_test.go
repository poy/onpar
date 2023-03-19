package fibonacci_test

import (
	"testing"

	"github.com/poy/onpar/v3"
	. "github.com/poy/onpar/v3/expect"
	. "github.com/poy/onpar/v3/matchers"
	"github.com/poy/onpar/v3/samples/fibonacci"
)

func TestDifferentInputs(t *testing.T) {
	o := onpar.New(t)

	o.Group("when n is 0", func() {
		o.Spec("it returns 1", func(t *testing.T) {
			result := fibonacci.Fibonacci(0)
			Expect(t, result).To(Equal(1))
		})
	})

	o.Group("when n is 1", func() {
		o.Spec("it returns 1", func(t *testing.T) {
			result := fibonacci.Fibonacci(1)
			Expect(t, result).To(Equal(1))
		})
	})

	o.Group("when n is greater than 1", func() {
		o.Spec("it returns 8 for n=5", func(t *testing.T) {
			result := fibonacci.Fibonacci(5)
			Expect(t, result).To(Equal(8))
		})

		o.Spec("it returns 55 for n=9", func(t *testing.T) {
			result := fibonacci.Fibonacci(9)
			Expect(t, result).To(Equal(55))
		})
	})
}
