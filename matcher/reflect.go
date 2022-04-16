package matcher

type reflectionMatcher interface {
	Match(v interface{}) (interface{}, error)
}

// ReflectMatcher is a generic matcher that strips out the generics, allowing
// old reflection-based matchers (from v1) to be used in generic contexts.
type ReflectMatcher[T any] struct {
	matchers []reflectionMatcher
}

// Reflect returns a ReflectionMatcher, which strips the generics from the
// actual value before passing it in to m.
//
// For the sake of brevity, this reflect matcher allows multiple matchers, which
// will in turn act like the old `matchers.Chain`. See the docs for
// `"github.com/poy/onpar/matchers".Chain` for details.
func Reflect[T any](first reflectionMatcher, chain ...reflectionMatcher) Matcher[T] {
	return ReflectMatcher[T]{
		matchers: append([]reflectionMatcher{first}, chain...),
	}
}

// Match calls the child matcher with v.
func (m ReflectMatcher[T]) Match(v T) error {
	var curr interface{} = v
	for _, matcher := range m.matchers {
		next, err := matcher.Match(curr)
		if err != nil {
			return err
		}
		curr = next
	}
	return nil
}
