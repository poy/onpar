// This file was generated by git.sr.ht/~nelsam/hel/v4.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package diff_test

import (
	"context"
	"time"

	"git.sr.ht/~nelsam/hel/v4/vegr"
	"github.com/poy/onpar/v2/diff/str"
)

type mockSprinter struct {
	t            vegr.T
	timeout      time.Duration
	SprintCalled chan bool
	SprintInput  struct {
		Arg0 chan []any
	}
	SprintOutput struct {
		Ret0 chan string
	}
}

func newMockSprinter(t vegr.T, timeout time.Duration) *mockSprinter {
	m := &mockSprinter{t: t, timeout: timeout}
	m.SprintCalled = make(chan bool, 100)
	m.SprintInput.Arg0 = make(chan []any, 100)
	m.SprintOutput.Ret0 = make(chan string, 100)
	return m
}
func (m *mockSprinter) Sprint(arg0 ...any) (ret0 string) {
	m.t.Helper()
	m.SprintCalled <- true
	m.SprintInput.Arg0 <- arg0
	vegr.PopulateReturns(m.t, "Sprint", m.timeout, m.SprintOutput, &ret0)
	return ret0
}

type mockStringDiffAlgorithm struct {
	t           vegr.T
	timeout     time.Duration
	DiffsCalled chan bool
	DiffsInput  struct {
		Ctx              chan context.Context
		Actual, Expected chan []rune
	}
	DiffsOutput struct {
		Ret0 chan chan str.Diff
	}
}

func newMockStringDiffAlgorithm(t vegr.T, timeout time.Duration) *mockStringDiffAlgorithm {
	m := &mockStringDiffAlgorithm{t: t, timeout: timeout}
	m.DiffsCalled = make(chan bool, 100)
	m.DiffsInput.Ctx = make(chan context.Context, 100)
	m.DiffsInput.Actual = make(chan []rune, 100)
	m.DiffsInput.Expected = make(chan []rune, 100)
	m.DiffsOutput.Ret0 = make(chan chan str.Diff, 100)
	return m
}
func (m *mockStringDiffAlgorithm) Diffs(ctx context.Context, actual, expected []rune) (ret0 chan str.Diff) {
	m.t.Helper()
	m.DiffsCalled <- true
	m.DiffsInput.Ctx <- ctx
	m.DiffsInput.Actual <- actual
	m.DiffsInput.Expected <- expected
	vegr.PopulateReturns(m.t, "Diffs", m.timeout, m.DiffsOutput, &ret0)
	return ret0
}

type mockContext struct {
	t              vegr.T
	timeout        time.Duration
	DeadlineCalled chan bool
	DeadlineOutput struct {
		Deadline chan time.Time
		Ok       chan bool
	}
	DoneCalled chan bool
	DoneOutput struct {
		Ret0 chan (<-chan struct{})
	}
	ErrCalled chan bool
	ErrOutput struct {
		Ret0 chan error
	}
	ValueCalled chan bool
	ValueInput  struct {
		Key chan any
	}
	ValueOutput struct {
		Ret0 chan any
	}
}

func newMockContext(t vegr.T, timeout time.Duration) *mockContext {
	m := &mockContext{t: t, timeout: timeout}
	m.DeadlineCalled = make(chan bool, 100)
	m.DeadlineOutput.Deadline = make(chan time.Time, 100)
	m.DeadlineOutput.Ok = make(chan bool, 100)
	m.DoneCalled = make(chan bool, 100)
	m.DoneOutput.Ret0 = make(chan (<-chan struct{}), 100)
	m.ErrCalled = make(chan bool, 100)
	m.ErrOutput.Ret0 = make(chan error, 100)
	m.ValueCalled = make(chan bool, 100)
	m.ValueInput.Key = make(chan any, 100)
	m.ValueOutput.Ret0 = make(chan any, 100)
	return m
}
func (m *mockContext) Deadline() (deadline time.Time, ok bool) {
	m.t.Helper()
	m.DeadlineCalled <- true
	vegr.PopulateReturns(m.t, "Deadline", m.timeout, m.DeadlineOutput, &deadline, &ok)
	return deadline, ok
}
func (m *mockContext) Done() (ret0 <-chan struct{}) {
	m.t.Helper()
	m.DoneCalled <- true
	vegr.PopulateReturns(m.t, "Done", m.timeout, m.DoneOutput, &ret0)
	return ret0
}
func (m *mockContext) Err() (ret0 error) {
	m.t.Helper()
	m.ErrCalled <- true
	vegr.PopulateReturns(m.t, "Err", m.timeout, m.ErrOutput, &ret0)
	return ret0
}
func (m *mockContext) Value(key any) (ret0 any) {
	m.t.Helper()
	m.ValueCalled <- true
	m.ValueInput.Key <- key
	vegr.PopulateReturns(m.t, "Value", m.timeout, m.ValueOutput, &ret0)
	return ret0
}

type mockDiff struct {
	t          vegr.T
	timeout    time.Duration
	CostCalled chan bool
	CostOutput struct {
		Ret0 chan float64
	}
	SectionsCalled chan bool
	SectionsOutput struct {
		Ret0 chan []str.DiffSection
	}
}

func newMockDiff(t vegr.T, timeout time.Duration) *mockDiff {
	m := &mockDiff{t: t, timeout: timeout}
	m.CostCalled = make(chan bool, 100)
	m.CostOutput.Ret0 = make(chan float64, 100)
	m.SectionsCalled = make(chan bool, 100)
	m.SectionsOutput.Ret0 = make(chan []str.DiffSection, 100)
	return m
}
func (m *mockDiff) Cost() (ret0 float64) {
	m.t.Helper()
	m.CostCalled <- true
	vegr.PopulateReturns(m.t, "Cost", m.timeout, m.CostOutput, &ret0)
	return ret0
}
func (m *mockDiff) Sections() (ret0 []str.DiffSection) {
	m.t.Helper()
	m.SectionsCalled <- true
	vegr.PopulateReturns(m.t, "Sections", m.timeout, m.SectionsOutput, &ret0)
	return ret0
}
