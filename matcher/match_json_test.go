package matcher_test

import (
	"testing"

	"github.com/poy/onpar/v2/matcher"
)

func TestMatchJSON(t *testing.T) {
	t.Parallel()

	t.Run("object", func(t *testing.T) {
		obj := `{"a": 99}`
		m := matcher.MatchJSON([]byte(obj))
		if err := m.Match([]byte(obj)); err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if err := m.Match([]byte(`{"different": 99}`)); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("list", func(t *testing.T) {
		list := `[1, 2, 3]`
		m := matcher.MatchJSON(list)

		if err := m.Match(list); err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if err := m.Match(`[3, 2, 1]`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("string", func(t *testing.T) {
		str := `"foo"`
		m := matcher.MatchJSON(str)

		if err := m.Match(str); err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if err := m.Match(`"bar"`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("numeric", func(t *testing.T) {
		num := `42`
		m := matcher.MatchJSON(num)

		if err := m.Match(num); err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if err := m.Match(`3.7`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		bool := `true`
		m := matcher.MatchJSON(bool)

		if err := m.Match(bool); err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if err := m.Match(`false`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})
}
