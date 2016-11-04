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
	description   string
	its           []it

	children []*level
	parent   *level

	beforeEachArgs []reflect.Value
}

type it struct {
	description string
	f           *reflect.Value
}

func BeforeEach(f interface{}) {
	if current.before != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered BeforeEach", current.description))
	}

	v := reflect.ValueOf(f)
	current.before = &v
}

func AfterEach(f interface{}) {
	if current.after != nil {
		panic(fmt.Sprintf("Level '%s' already has a registered AfterEach", current.description))
	}

	v := reflect.ValueOf(f)
	current.after = &v
}

func It(description string, f interface{}) {
	v := reflect.ValueOf(f)
	current.its = append(current.its, it{description: description, f: &v})
}

func Describe(description string, f func()) {
	newLevel := &level{
		description: description,
		parent:      current,
	}

	current.children = append(current.children, newLevel)

	oldLevel := current
	current = newLevel
	f()
	current = oldLevel
}

func Run(t *testing.T) {

	t.Parallel()

	traverse(current, func(l *level) bool {
		for _, i := range l.its {
			t.Run(i.description, func(tt *testing.T) {
				tt.Parallel()
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

				i.f.Call(args)

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
			})
		}
		return true
	})

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
