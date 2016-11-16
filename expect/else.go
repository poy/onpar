package expect

import "github.com/apoydence/onpar/matchers"

type Else struct {
	t  T
	to *To
}

func (e *Else) FailNow() {
	if e.to.parentErr == nil {
		return
	}

	e.t.Fatal(e.to.parentErr.Error())
}

func (e *Else) To(matcher matchers.Matcher) {
	if e.to.parentErr == nil {
		return
	}

	e.to.parentErr = nil
	e.to.To(matcher)
}

//func (t *To) To(matcher matchers.Matcher) *ChainedTo {
//	err := t.parentErr
//	var resultValue interface{}
//	if t.parentErr == nil {
//		resultValue, err = matcher.Match(t.actual)
//		if err != nil {
//			t.t.Errorf(err.Error())
//		}
//	}
//
//	newToCurrent := &To{
//		actual:    t.actual,
//		parentErr: err,
//		t:         t.t,
//	}
//	newToNext := &To{
//		actual:    resultValue,
//		parentErr: err,
//		t:         t.t,
//	}
//
//	return &ChainedTo{
//		t:          t.t,
//		And:        newToCurrent,
//		AndForThat: newToNext,
//		Else: &Else{
//			err: err,
//			t:   t.t,
//		},
//	}
//}
