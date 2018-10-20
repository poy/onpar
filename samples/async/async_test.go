package async_test

import (
	"testing"

	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestChannel(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 0; i < 100; i++ {
			c <- i
		}
	}()

	Expect(t, c).To(ViaPolling(
		Chain(Receive(), Equal(50)),
	))
}
