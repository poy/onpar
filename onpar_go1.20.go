//go:build go1.20

package onpar

func init() {
	panic("onpar: go 1.20.x introduced a breaking change that broke onpar v2, and we could not solve it without a breaking change of our own. Please upgrade to onpar v3.")
}
