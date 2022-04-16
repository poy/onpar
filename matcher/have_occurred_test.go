package matcher_test

import (
	"fmt"
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestHaveOccurred(t *testing.T) {
	m := matcher.HaveOccurred()

	if err := m.Match(fmt.Errorf("some-err")); err != nil {
		t.Fatal("expected err to be nil")
	}

	var e error
	if err := m.Match(e); err == nil {
		t.Fatal("expected err to not be nil")
	}

	if err := m.Match(nil); err == nil {
		t.Fatal("expected err to not be nil")
	}
}
