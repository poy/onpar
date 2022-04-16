package matcher_test

import (
	"testing"

	"github.com/poy/onpar/matchers"
	"github.com/poy/onpar/v2/matcher"
)

func TestReflectV1Matcher(t *testing.T) {
	t.Parallel()

	m := matchers.Equal(10)
	rm := matcher.Reflect[int](m)

	if err := rm.Match(10); err != nil {
		t.Errorf("expected the v1 matchers.Equal to pass; got %v", err)
	}

	if err := rm.Match(11); err == nil {
		t.Errorf("expected  the v1 matchers.Equal to fail; got <nil>")
	}
}

func TestReflectV1Chain(t *testing.T) {
	rm := matcher.Reflect[map[string]int](
		matchers.HaveKey("foo"),
		matchers.Equal(1),
	)
	if err := rm.Match(map[string]int{"foo": 1}); err != nil {
		t.Errorf("expected the v1 HaveKey and Equal matchers to chain and pass with the expected map; got %v", err)
	}

	if err := rm.Match(map[string]int{"foo": 2}); err == nil {
		t.Errorf("expected the v1 HaveKey and Equal matchers to chain and fail with an unexpected map value; got %v", err)
	}

	if err := rm.Match(map[string]int{"bar": 1}); err == nil {
		t.Errorf("expected the v1 HaveKey and Equal matchers to chain and fail with an unexpected map key; got %v", err)
	}
}
