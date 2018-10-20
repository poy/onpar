package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestContain(t *testing.T) {
	t.Parallel()
	m := matchers.Contain("a", 1)
	values := []interface{}{"a", "b", "c", 1, 2, 3}
	v, err := m.Match(values)

	if err != nil {
		t.Fatalf("expected err (%s) to be nil", err)
	}

	if !reflect.DeepEqual(v, values) {
		t.Fatalf("expected %v to equal %v", v, values)
	}

	_, err = m.Match([]interface{}{"d", "e"})
	if err == nil {
		t.Fatal("expected err to not be nil")
	}

	_, err = m.Match("invalid")
	if err == nil {
		t.Fatal("expected err to not be nil")
	}
}
