package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	m := matchers.Equal(101)
	v, err := m.Match(101)
	if err != nil {
		t.Errorf("expected %v to be nil", err)
	}

	if !reflect.DeepEqual(v, 101) {
		t.Errorf("expected %v to equal %v", v, 101)
	}

	_, err = m.Match(103)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
}
