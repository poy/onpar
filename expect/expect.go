package expect

import (
	"runtime"

	"github.com/apoydence/onpar/matchers"
)

// T is a type that we can perform assertions with.
type T interface {
	Fatalf(format string, args ...interface{})
}

type Expectation func(actual interface{}) *To

func New(t T) Expectation {
	return func(actual interface{}) *To {
		return &To{
			actual: actual,
			t:      t,
		}
	}
}

func Expect(t T, actual interface{}) *To {
	return &To{
		actual: actual,
		t:      t,
	}
}

type To struct {
	actual    interface{}
	parentErr error

	t T
}

func (t *To) To(matcher matchers.Matcher) {
	_, err := matcher.Match(t.actual)
	if err != nil {
		_, fileName, lineNumber, _ := runtime.Caller(1)
		t.t.Fatalf("%s\n%s:%d", err.Error(), fileName, lineNumber)
	}
}
