package expect

import (
	"path"
	"runtime"

	"github.com/poy/onpar/matchers"
)

// T is a type that we can perform assertions with.
type T interface {
	Fatalf(format string, args ...interface{})
}

// THelper has the method that tells the testing framework that it can declare
// itself a test helper.
type THelper interface {
	Helper()
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
	if helper, ok := t.t.(THelper); ok {
		helper.Helper()
	}

	_, err := matcher.Match(t.actual)
	if err != nil {
		_, fileName, lineNumber, _ := runtime.Caller(1)
		t.t.Fatalf("%s\n%s:%d", err.Error(), path.Base(fileName), lineNumber)
	}
}
