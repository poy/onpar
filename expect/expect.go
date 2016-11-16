package expect

import "github.com/apoydence/onpar/matchers"

// T is a type that we can perform assertions with.
type T interface {
	Errorf(format string, args ...interface{})
	Fatal(...interface{})
	FailNow()
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

type ChainedTo struct {
	t          T
	And        *To
	AndForThat *To
	Else       *Else
}

func (t *To) To(matcher matchers.Matcher) *ChainedTo {
	err := t.parentErr
	var resultValue interface{}
	if t.parentErr == nil {
		resultValue, err = matcher.Match(t.actual)
		if err != nil {
			t.t.Errorf(err.Error())
		}
	}

	newToCurrent := &To{
		actual:    t.actual,
		parentErr: err,
		t:         t.t,
	}
	newToNext := &To{
		actual:    resultValue,
		parentErr: err,
		t:         t.t,
	}

	return &ChainedTo{
		t:          t.t,
		And:        newToCurrent,
		AndForThat: newToNext,
		Else: &Else{
			t:  t.t,
			to: newToCurrent,
		},
	}
}
