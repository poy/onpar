package matchers_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/matchers"
)

func TestMatchJSON(t *testing.T) {
	t.Parallel()

	m := matchers.MatchJSON([]byte(`{"a": 99}`))
	v, err := m.Match(`{"a": 99}`)
	if err != nil {
		t.Errorf("expected %v to be nil", err)
	}

	if !reflect.DeepEqual(v, `{"a": 99}`) {
		t.Errorf("expected %v to equal %v", v, `{"a": 99}`)
	}

	_, err = m.Match(`{"different": 99}`)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
}
