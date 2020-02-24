package matchers_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/poy/onpar/diff"
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

func TestEqualDiff(t *testing.T) {
	t.Parallel()

	m := matchers.Equal(101)
	m.UseDiffOpts(diff.Actual(diff.WithFormat("%s!=")))
	_, err := m.Match(103)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
	format := fmt.Sprintf("103 to equal 101\ndiff: 103!=101")
	if err.Error() != format {
		t.Fatalf("expected '%v' to match '%v'", err.Error(), format)
	}
}
