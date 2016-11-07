package onpar_test

import (
	"testing"

	"github.com/apoydence/onpar"
)

func TestPassBeforeEachOutputToIt(t *testing.T) {
	c := make(chan string, 100)
	onpar.AfterEach(func(t *testing.T) {
	})

	onpar.Group("DA", func() {
		onpar.BeforeEach(func(t *testing.T) (int, string, chan string) {
			c <- "DA-BeforeEach"
			return 99, "something", c
		})

		onpar.AfterEach(func(t *testing.T, i int, s string, c chan string) {
			c <- "DA-AfterEach"

			if i != 99 {
				t.Errorf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Errorf("expected %s = %s", s, "something")
			}
		})

		onpar.Spec("A", func(t *testing.T, i int, s string, c chan string) {
			c <- "DA-A"
			if i != 99 {
				t.Errorf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Errorf("expected %s = %s", s, "something")
			}
		})

		onpar.Group("DB", func() {
			onpar.BeforeEach(func(t *testing.T, i int, s string, c chan string) float64 {
				c <- "DB-BeforeEach"
				return 101
			})

			onpar.Spec("B", func(t *testing.T, i int, s string, c chan string, f float64) {
				c <- "DB-BEach"
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

	t.Run("", func(tt *testing.T) {
		onpar.Run(tt)
	})

	if len(c) != 7 {
		t.Errorf("expected c (len=%d) to have len %d", len(c), 7)
	}
}
