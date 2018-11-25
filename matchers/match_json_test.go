package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestMatchJSON(t *testing.T) {
	t.Parallel()

	t.Run("object", func(t *testing.T) {
		obj := `{"a": 99}`
		m := matchers.MatchJSON([]byte(obj))
		v, err := m.Match(obj)
		if err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if !reflect.DeepEqual(v, obj) {
			t.Errorf("expected %v to equal %v", v, obj)
		}

		if _, err := m.Match(`{"different": 99}`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("list", func(t *testing.T) {
		list := `[1, 2, 3]`
		m := matchers.MatchJSON(list)
		v, err := m.Match(list)

		if err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if !reflect.DeepEqual(v, list) {
			t.Errorf("expected %v to equal %v", v, list)
		}

		if _, err := m.Match(`[3, 2, 1]`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("string", func(t *testing.T) {
		str := `"foo"`
		m := matchers.MatchJSON(str)
		v, err := m.Match(str)

		if err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if !reflect.DeepEqual(v, str) {
			t.Errorf("expected %v to equal %v", v, str)
		}

		if _, err := m.Match(`"bar"`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("numeric", func(t *testing.T) {
		num := `42`
		m := matchers.MatchJSON(num)
		v, err := m.Match(num)

		if err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if !reflect.DeepEqual(v, num) {
			t.Errorf("expected %v to equal %v", v, num)
		}

		if _, err := m.Match(`3.7`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		bool := `true`
		m := matchers.MatchJSON(bool)
		v, err := m.Match(bool)

		if err != nil {
			t.Errorf("expected %v to be nil", err)
		}

		if !reflect.DeepEqual(v, bool) {
			t.Errorf("expected %v to equal %v", v, bool)
		}

		if _, err := m.Match(`false`); err == nil {
			t.Fatalf("expected %v to not be nil", err)
		}
	})
}
