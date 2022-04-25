package onpar

import (
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/poy/onpar/v2/diff"
)

type prefs struct {
}

// Opt is an option type to pass to onpar's constructor.
type Opt func(prefs) prefs

type suite[T any] interface {
	addRunner(runner[T])
	child() child
}

type child interface {
	addSpecs()
}

// Onpar stores the state of the specs and groups
type Onpar[T, U any] struct {
	path []string

	parent suite[T]

	// level is handled by (*Onpar[T]).Group(), which will adjust this field
	// each time it is called. This is how specs get the level name added to
	// their names and how BeforeEach knows which level it is creating a new
	// sub-suite for.
	level *level[T, U]

	// childSuite is assigned by BeforeEach and removed at the end of Group. If
	// BeforeEach is called twice in the same Group (or twice at the top level),
	// this is how it knows to panic.
	//
	// At the end of Group calls, childSuite.addSpecs is called, which will sync the
	// childSuite's specs to the parent.
	childSuite child
	childPath  []string

	// TODO: why are these here, again?
	diffOpts []diff.Opt
}

// New creates a new Onpar suite. The top-level onpar suite must be constructed
// with this. Think `context.Background()`.
//
// It's normal to construct the top-level suite with a BeforeEach by doing the
// following:
//
//     o := BeforeEach(New(t), setupFn)
func New(t *testing.T, opts ...Opt) *Onpar[*testing.T, *testing.T] {
	p := prefs{}
	for _, opt := range opts {
		p = opt(p)
	}
	o := Onpar[*testing.T, *testing.T]{
		level: &level[*testing.T, *testing.T]{
			before: func(t *testing.T) *testing.T {
				return t
			},
		},
	}
	t.Cleanup(func() {
		o.run(t)
	})
	return &o
}

// BeforeEach creates a new child Onpar suite with the requested function as the
// setup function for all specs. It requires a parent Onpar.
//
// The top level Onpar *must* have been constructed with New, otherwise the
// suite will not run.
//
// BeforeEach should be called only once for each level (i.e. each group). It
// will panic if it detects that it is overwriting another BeforeEach call for a
// given level.
func BeforeEach[T, U, V any](parent *Onpar[T, U], setup func(U) V) *Onpar[U, V] {
	if !parent.correctGroup() {
		panic(fmt.Errorf("onpar: BeforeEach called on child suite outside of its group (%v)", path.Join(parent.path...)))
	}
	if parent.child() != nil {
		if len(parent.childPath) == 0 {
			panic(errors.New("onpar: BeforeEach was called more than once at the top level"))
		}
		panic(fmt.Errorf("onpar: BeforeEach was called more than once for group '%s'", path.Join(parent.childPath...)))
	}
	path := parent.path
	if parent.level.name() != "" {
		path = append(parent.path, parent.level.name())
	}
	child := &Onpar[U, V]{
		path:   path,
		parent: parent,
		level: &level[U, V]{
			before: setup,
		},
	}
	parent.childSuite = child
	parent.childPath = child.path
	return child
}

// Spec is a test that runs in parallel with other specs.
func (o *Onpar[T, U]) Spec(name string, f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: Spec called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	spec := concurrentSpec[U]{
		serialSpec: serialSpec[U]{
			specName: name,
			f:        f,
		},
	}
	o.addRunner(spec)
}

// SerialSpec is a test that runs synchronously (i.e. onpar will not call
// `t.Parallel`). While onpar is primarily a parallel testing suite, we
// recognize that sometimes a test just can't be run in parallel. When that is
// the case, use SerialSpec.
func (o *Onpar[T, U]) SerialSpec(name string, f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: SerialSpec called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	spec := serialSpec[U]{
		specName: name,
		f:        f,
	}
	o.addRunner(spec)
}

func (o *Onpar[T, U]) addRunner(r runner[U]) {
	o.level.runners = append(o.level.runners, r)
}

// Group is used to gather and categorize specs. Inside of each group, a new
// child *Onpar may be constructed using BeforeEach.
func (o *Onpar[T, U]) Group(name string, f func()) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: Group called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	oldLevel := o.level
	o.level = &level[T, U]{
		levelName: name,
	}
	defer func() {
		if o.child() != nil {
			o.child().addSpecs()
			o.childSuite = nil
		}
		oldLevel.runners = append(oldLevel.runners,
			&level[U, U]{
				levelName: o.level.name(),
				before: func(v U) U {
					return v
				},
				runners: o.level.runners,
			})
		o.level = oldLevel
	}()

	f()
}

// AfterEach is used to cleanup anything from the specs or BeforeEaches.
// AfterEach may only be called once for each *Onpar value constructed.
func (o *Onpar[T, U]) AfterEach(f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: AfterEach called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	if o.level.after != nil {
		if len(o.childPath) == 0 {
			panic(errors.New("onpar: AfterEach was called more than once at top level"))
		}
		panic(fmt.Errorf("onpar: AfterEach was called more than once for group '%s'", path.Join(o.path...)))
	}
	o.level.after = f
}

func (o *Onpar[T, U]) run(t *testing.T) {
	if o.child() != nil {
		// This happens when New is called before BeforeEach, e.g.:
		//
		//     o := onpar.New()
		//     defer o.Run(t)
		//
		//     b := onpar.BeforeEach(o, setup)
		//
		// Since there's no call to o.Group, the child won't be synced, so we
		// need to do that here.
		o.child().addSpecs()
		o.childSuite = nil
	}
	top, ok := interface{}(o.level).(runner[*testing.T])
	if !ok {
		var empty T
		panic(fmt.Errorf("onpar: Run was called on a child level (type '%T' is not *testing.T)", empty))
	}
	top.runSpecs(t, func(t *testing.T) *testing.T {
		return t
	}, nil)
}

func (o *Onpar[T, U]) child() child {
	return o.childSuite
}

func (o *Onpar[T, U]) correctGroup() bool {
	if o.parent == nil {
		return true
	}
	if o.parent.child() == o {
		return true
	}
	return false
}

// addSpecs is called by parent Group() calls to tell o to add its specs to its
// parent.
func (o *Onpar[T, U]) addSpecs() {
	o.parent.addRunner(o.level)
}

type runner[T any] interface {
	name() string
	runSpecs(t *testing.T, before func(*testing.T) T, after func(T))
}

type concurrentSpec[T any] struct {
	serialSpec[T]
}

func (s concurrentSpec[T]) runSpecs(t *testing.T, before func(*testing.T) T, after func(T)) {
	t.Parallel()

	s.serialSpec.runSpecs(t, before, after)
}

type serialSpec[T any] struct {
	specName string
	f        func(T)
}

func (s serialSpec[T]) name() string {
	return s.specName
}

func (s serialSpec[T]) runSpecs(t *testing.T, before func(*testing.T) T, after func(T)) {
	v := before(t)
	s.f(v)
	if after != nil {
		after(v)
	}
}

type level[T, U any] struct {
	levelName string
	before    func(T) U
	after     func(U)
	runners   []runner[U]
}

func (l *level[T, U]) name() string {
	return l.levelName
}

func (l *level[T, U]) runSpecs(t *testing.T, before func(*testing.T) T, after func(T)) {
	for _, r := range l.runners {
		testFn := func(t *testing.T) {
			var v T
			childBefore := func(t *testing.T) U {
				v = before(t)
				return l.before(v)
			}
			childAfter := func(childV U) {
				if l.after != nil {
					l.after(childV)
				}
				if after != nil {
					after(v)
				}
			}
			r.runSpecs(t, childBefore, childAfter)
		}
		if r.name() == "" {
			// If the name is empty, running the group as a sub-group would
			// result in ugly output. Just run the test function at this level
			// instead.
			testFn(t)
			continue
		}
		t.Run(r.name(), testFn)
	}
}
