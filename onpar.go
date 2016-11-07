package onpar

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	current *level
)

func init() {
	current = &level{}
}

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

func BeforeEach(f interface{}) {
	if current.before != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered BeforeEach", current.name))
	}

	v := reflect.ValueOf(f)
	current.before = &v
}

func AfterEach(f interface{}) {
	if current.after != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered AfterEach", current.name))
	}

	v := reflect.ValueOf(f)
	current.after = &v
}

func Spec(name string, f interface{}) {
	v := reflect.ValueOf(f)
	spec := specInfo{
		name: name,
		f:    &v,
	}
	current.specs = append(current.specs, spec)
}

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
	var desc string
	traverse(current, func(ll *level) bool {
		desc = fmt.Sprintf("%s/%s", desc, ll.name)
		return ll != l
	})

	return fmt.Sprintf("%s/%s", desc, i.name)
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
