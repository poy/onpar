package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestHaveKey(t *testing.T) {
	m := matcher.HaveKey[map[int]string, int, string](99)

	if err := m.Match(map[int]string{1: "a"}); err == nil {
		t.Error("expected err to not be nil")
	}

	if err := m.Match(map[int]string{99: "a"}); err != nil {
		t.Fatal("expected err to be nil")
	}
}
