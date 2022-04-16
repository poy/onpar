package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestContain(t *testing.T) {
	t.Parallel()
	m := matcher.Contain[[]string]("a", "b")
	values := []string{"a", "b", "c"}

	if err := m.Match(values); err != nil {
		t.Fatalf("expected err (%s) to be nil", err)
	}

	if err := m.Match([]string{"d", "e"}); err == nil {
		t.Fatal("expected err to not be nil")
	}
}
