package matchers_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/expect"
	"github.com/poy/onpar/v2/matchers"
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
	mockDiffer := newMockDiffer(t, time.Second)
	pers.Return(mockDiffer.DiffOutput, "this is a valid diff")
	m.UseDiffer(mockDiffer)
	_, err := m.Match(103)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
	expect.Expect(t, mockDiffer).To(pers.HaveMethodExecuted("Diff", pers.WithArgs(103, 101)))
	format := fmt.Sprintf("expected 103 to equal 101\ndiff: this is a valid diff")
	if err.Error() != format {
		t.Fatalf("expected '%v' to match '%v'", err.Error(), format)
	}
}
