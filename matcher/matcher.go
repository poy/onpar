package matcher

// Matcher is a type that matches expected against actuals.
type Matcher[T any] interface {
	Match(actual T) error
}
