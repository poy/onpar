package async_test

import (
	"testing"
	"time"

	. "github.com/poy/onpar/v2/expect"
	. "github.com/poy/onpar/v2/matcher"
)

func TestChannel(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 0; i < 100; i++ {
			c <- i
		}
	}()

	Expect(t, c).To(Eventually[chan int, int](Equal(50), EventuallyTimes(100), EventuallyInterval(time.Microsecond)))
}
