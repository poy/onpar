package matchers_test

import (
	"fmt"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestHaveOccurred(t *testing.T) {
	m := matchers.HaveOccurred()

	_, err := m.Match(fmt.Errorf("some-err"))
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	var e error
	_, err = m.Match(e)
	if err == nil {
		t.Fatal("expected err to not be nil")
	}

	_, err = m.Match(nil)
	if err == nil {
		t.Fatal("expected err to not be nil")
	}
}
