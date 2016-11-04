package onpar_test

import (
	"testing"

	"github.com/apoydence/onpar"
)

func TestPassBeforeEachOutputToIt(t *testing.T) {
	onpar.AfterEach(func(t *testing.T) {
	})

	onpar.Describe("DA", func() {
		onpar.BeforeEach(func(t *testing.T) (int, string) {
			return 99, "something"
		})

		onpar.AfterEach(func(t *testing.T, i int, s string) {
			if i != 99 {
				t.Errorf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Errorf("expected %s = %s", s, "something")
			}
		})

		onpar.It("A", func(t *testing.T, i int, s string) {
			if i != 99 {
				t.Errorf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Errorf("expected %s = %s", s, "something")
			}
		})

		onpar.Describe("DB", func() {
			onpar.BeforeEach(func(t *testing.T, i int, s string) float64 {
				return 101
			})

			onpar.It("B", func(t *testing.T, i int, s string, f float64) {
				if i != 99 {
					t.Errorf("expected %d = %d", i, 99)
				}

				if s != "something" {
					t.Errorf("expected %s = %s", s, "something")
				}

				if f != 101 {
					t.Errorf("expected %f = %f", f, 101.0)
				}
			})
		})
	})

	onpar.Run(t)
}
