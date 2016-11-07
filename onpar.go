package onpar

import (
	"fmt"
	"reflect"
	"testing"
)

// Spec is a test that runs in parallel with other specs. The provided function
// takes the `testing.T` for test assertions and any arguments the `BeforeEach()`
// returns.
func Spec(name string, f interface{}) {
	v := reflect.ValueOf(f)
	spec := specInfo{
		name: name,
		f:    &v,
	}
	current.specs = append(current.specs, spec)
}

// Group is used to gather and categorize specs. Each group can have a single
// `BeforeEach()` and `AfterEach()`.
func Group(name string, f func()) {
	newLevel := &level{
		name:   name,
		parent: current,
	}

	current.children = append(current.children, newLevel)

	oldLevel := current
	current = newLevel
	f()
	current = oldLevel
}

// BeforeEach is used for any setup that may be required for the specs.
// Each argument returned will be required to be received by following specs.
// Outer BeforeEaches are invoked before inner ones.
func BeforeEach(f interface{}) {
	if current.before != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered BeforeEach", current.name))
	}

	v := reflect.ValueOf(f)
	current.before = &v
}

// AfterEach is used to cleanup anything from the specs or BeforeEaches.
// The function takes arguments the same as specs. Inner AfterEaches are invoked
// before outer ones.
func AfterEach(f interface{}) {
	if current.after != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered AfterEach", current.name))
	}

	v := reflect.ValueOf(f)
	current.after = &v
}

// Run is used to initiate the tests.
func Run(t *testing.T) {

	traverse(current, func(l *level) bool {
		for _, spec := range l.specs {
			desc := buildDesc(l, spec)
			t.Run(desc, func(tt *testing.T) {
				tt.Parallel()

				args, levelArgs := invokeBeforeEach(tt, l)
				spec.f.Call(args)

				invokeAfterEach(tt, l, levelArgs)
			})
		}
		return true
	})
}

var (
	current *level = new(level)
)

type level struct {
	before, after *reflect.Value
	name          string
	specs         []specInfo

	children []*level
	parent   *level

	beforeEachArgs []reflect.Value
}

type specInfo struct {
	name string
	f    *reflect.Value
}

func invokeBeforeEach(tt *testing.T, l *level) ([]reflect.Value, map[*level][]reflect.Value) {
	args := []reflect.Value{
		reflect.ValueOf(tt),
	}
	levelArgs := make(map[*level][]reflect.Value)

	traverse(current, func(ll *level) bool {
		if ll.before != nil {
			args = append(args, ll.before.Call(args)...)
			levelArgs[ll] = args
		}

		return ll != l
	})
	return args, levelArgs
}

func invokeAfterEach(tt *testing.T, l *level, levelArgs map[*level][]reflect.Value) {
	rTraverse(l, func(ll *level) {
		beforeEachArgs := levelArgs[ll]
		if beforeEachArgs == nil {
			beforeEachArgs = []reflect.Value{
				reflect.ValueOf(tt),
			}
		}

		if ll.after != nil {
			ll.after.Call(beforeEachArgs)
		}
	})
}

func buildDesc(l *level, i specInfo) string {
	desc := i.name
	rTraverse(l, func(ll *level) {
		desc = fmt.Sprintf("%s/%s", ll.name, desc)
	})

	return desc
}

func traverse(l *level, f func(*level) bool) {
	if l == nil {
		return
	}

	if !f(l) {
		return
	}

	for _, child := range l.children {
		traverse(child, f)
	}
}

func rTraverse(l *level, f func(*level)) {
	if l == nil {
		return
	}

	f(l)

	rTraverse(l.parent, f)
}
