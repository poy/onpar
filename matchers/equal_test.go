package matchers_test

import (
	"testing"

	"github.com/apoydence/onpar/matchers"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	m := matchers.Equal(101)
	_, err := m.Match(101)
	if err != nil {
		t.Errorf("expected %v to be nil", err)
	}

	_, err = m.Match(103)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
}
