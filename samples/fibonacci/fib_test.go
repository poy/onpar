package fibonacci_test

import (
	"testing"

	. "github.com/apoydence/onpar"
	"github.com/apoydence/onpar/samples/fibonacci"
)

func TestLessThanOne(t *testing.T) {
	Group("when n is 0", func() {
		Spec("it returns 1", func(tt *testing.T) {
			result := fibonacci.Fibonacci(0)

			if result != 1 {
				tt.Errorf("expected result (%d) to == %d ", result, 1)
			}
		})
	})

	Group("when n is 1", func() {
		Spec("it returns 1", func(tt *testing.T) {
			result := fibonacci.Fibonacci(1)

			if result != 1 {
				tt.Errorf("expected result (%d) to == %d ", result, 1)
			}
		})
	})

	Group("when n is greater than 1", func() {
		Spec("it returns 8 for n=5", func(tt *testing.T) {
			result := fibonacci.Fibonacci(5)

			if result != 8 {
				tt.Errorf("expected result (%d) to == %d ", result, 8)
			}
		})

		Spec("it returns 55 for n=9", func(tt *testing.T) {
			result := fibonacci.Fibonacci(9)

			if result != 55 {
				tt.Errorf("expected result (%d) to == %d ", result, 55)
			}
		})
	})

	Run(t)
}
