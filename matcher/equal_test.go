package matcher_test

import (
	"fmt"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/v2/expect"
	"github.com/poy/onpar/v2/matcher"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	m := matcher.Equal(101)
	if err := m.Match(101); err != nil {
		t.Errorf("expected %v to be nil", err)
	}

	if err := m.Match(103); err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
}

func TestEqualDiff(t *testing.T) {
	t.Parallel()

	m := matcher.Equal(101)
	mDiffer := newMockDiffer(t, time.Second)
	pers.Return(mDiffer.DiffOutput, "this is a valid diff")
	m.UseDiffer(mDiffer)
	err := m.Match(103)
	if err == nil {
		t.Fatalf("expected %v to not be nil", err)
	}
	expect.Expect(t, mDiffer).To(matcher.Reflect[*mockDiffer](pers.HaveMethodExecuted("Diff", pers.WithArgs(103, 101))))
	format := fmt.Sprintf("expected 103 to equal 101\ndiff: this is a valid diff")
	if err.Error() != format {
		t.Fatalf("expected '%v' to match '%v'", err.Error(), format)
	}
}
