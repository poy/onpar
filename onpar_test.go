package onpar_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/poy/onpar"
)

const testTimeout = time.Second

func TestPanicsWithMissingRun(t *testing.T) {
	t.Parallel()

	mockT := newMockTestRunner(t, testTimeout)
	var cleanup func()
	seq := pers.CallSequence(t)
	pers.Expect(seq, mockT, "Cleanup", pers.StoreArgs(&cleanup))

	onpar.New(mockT)

	seq.Check(t)

	if cleanup == nil {
		t.Fatalf("expected Cleanup to be called with a cleanup function")
	}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected onpar to panic if Run was never called")
		}
	}()

	pers.Expect(seq, mockT, "Failed", pers.Returning(false))
	cleanup()
}

func TestMissingRunIsNotMentionedIfTestPanics(t *testing.T) {
	t.Parallel()

	mockT := newMockTestRunner(t, testTimeout)
	seq := pers.CallSequence(t)
	var cleanup func()
	pers.Expect(seq, mockT, "Cleanup", pers.StoreArgs(&cleanup))

	onpar.New(mockT)

	seq.Check(t)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected test panic to surface")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected test panic to be surfaced as a string")
		}
		if strings.Contains(msg, "missing 'defer o.Run()'") {
			t.Fatalf("did not expect onpar to mention missing o.Run when calling context panicked")
		}
	}()

	pers.Expect(seq, mockT, "Failed", pers.Returning(false))
	defer cleanup()
	panic("boom")
}

func TestMissingRunIsNotMentionedIfTestIsFailed(t *testing.T) {
	t.Parallel()

	mockT := newMockTestRunner(t, testTimeout)
	seq := pers.CallSequence(t)
	var cleanup func()
	pers.Expect(seq, mockT, "Cleanup", pers.StoreArgs(&cleanup))

	onpar.New(mockT)

	seq.Check(t)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected onpar not to panic if the test is failed")
		}
	}()

	pers.Expect(seq, mockT, "Failed", pers.Returning(true))
	cleanup()
}

func TestMissingRunIsNotMentionedWithSpecPanic(t *testing.T) {
	t.Parallel()

	mockT := newMockTestRunner(t, testTimeout)
	seq := pers.CallSequence(t)
	var cleanup func()
	pers.Expect(seq, mockT, "Cleanup", pers.StoreArgs(&cleanup))
	pers.Panic(mockT.RunOutput, "boom")

	o := onpar.New(mockT)
	o.Spec("boom", func(t *testing.T) {
		panic("boom")
	})

	seq.Check(t)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected child panic to surface")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected child panic to be surfaced as a string")
		}
		if strings.Contains(msg, "missing 'defer o.Run()'") {
			t.Fatalf("did not expect onpar to mention missing o.Run when o.Run was called")
		}
	}()

	pers.Expect(seq, mockT, "Failed", pers.Returning(false))
	defer cleanup()
	o.Run()
}

func TestSingleNestedSpec(t *testing.T) {
	t.Parallel()
	c := createScaffolding(t)

	objs := chanToSlice(c)

	if len(objs) != 4 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objs), 4)
	}

	objA := findSpec(objs, "DA-A")
	if objA == nil {
		t.Fatal("unable to find spec A")
	}

	if len(objA.c) != 4 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objA.c), 4)
	}

	if !reflect.DeepEqual(objA.c, []string{"-BeforeEach", "DA-A", "DA-AfterEach", "-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objA.c)
	}
}

func TestInvokeFirstChildAndPeerSpec(t *testing.T) {
	t.Parallel()
	c := createScaffolding(t)

	objs := chanToSlice(c)

	objB := findSpec(objs, "DB-B")
	if objB == nil {
		t.Fatal("unable to find spec B")
	}

	if len(objB.c) != 6 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objB.c), 6)
	}

	if !reflect.DeepEqual(objB.c, []string{"-BeforeEach", "DB-BeforeEach", "DB-B", "DB-AfterEach", "DA-AfterEach", "-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objB.c)
	}
}

func TestInvokeSecondChildAndPeerSpec(t *testing.T) {
	t.Parallel()
	c := createScaffolding(t)

	objs := chanToSlice(c)

	objC := findSpec(objs, "DB-C")
	if objC == nil {
		t.Fatal("unable to find spec C")
	}

	if len(objC.c) != 6 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objC.c), 6)
	}

	if !reflect.DeepEqual(objC.c, []string{"-BeforeEach", "DB-BeforeEach", "DB-C", "DB-AfterEach", "DA-AfterEach", "-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objC.c)
	}
}

func TestNewWithBeforeEach(t *testing.T) {
	c := make(chan string, 100)

	t.Run("FakeSpecs", func(t *testing.T) {
		o := onpar.New(t)
		defer o.Run()

		o.SerialSpec("it runs a spec without a beforeeach", func(*testing.T) {
			c <- "A"
		})

		b := onpar.BeforeEach(o, func(*testing.T) string {
			c <- "B-BeforeEach"
			return "foo"
		})

		b.SerialSpec("it runs a spec on a BeforeEach", func(string) {
			c <- "B"
		})
	})

	expected := []string{"A", "B-BeforeEach", "B"}
	for len(expected) > 0 {
		select {
		case v := <-c:
			if v != expected[0] {
				t.Fatalf("expected %v, got %v", expected[0], v)
				return
			}
			expected = expected[1:]
		default:
			t.Fatalf("expected %v to be called but it never was", expected[0])
			return
		}
	}
}

func TestGroupNestsRunCalls(t *testing.T) {
	c := make(chan string, 100)

	t.Run("FakeSpecs", func(t *testing.T) {
		o := onpar.New(t)
		defer o.Run()

		sendName := func(t *testing.T) {
			c <- t.Name()
		}

		o.Spec("A", sendName)

		o.Group("B", func() {
			o.Spec("C", sendName)

			o.Spec("D", sendName)
		})

		o.Spec("E", sendName)

		o.Group("F", func() {
			o.Group("G", func() {
				o.Spec("H", sendName)
			})
		})
	})

	expected := []string{
		"TestGroupNestsRunCalls/FakeSpecs/A",
		"TestGroupNestsRunCalls/FakeSpecs/B/C",
		"TestGroupNestsRunCalls/FakeSpecs/B/D",
		"TestGroupNestsRunCalls/FakeSpecs/E",
		"TestGroupNestsRunCalls/FakeSpecs/F/G/H",
	}

	findMatch := func(v string) {
		// We aren't guaranteed order here since the specs run in parallel.
		for i, e := range expected {
			if v == e {
				expected = append(expected[:i], expected[i+1:]...)
				return
			}
		}
		t.Fatalf("test name %v was not expected (or was run twice)", v)
	}

	for len(expected) > 0 {
		select {
		case v := <-c:
			findMatch(v)
		default:
			t.Fatalf("specs %v were never called", expected)
			return
		}
	}
}

func TestSerialSpecsAreOrdered(t *testing.T) {
	c := make(chan string, 100)

	t.Run("FakeSpecs", func(t *testing.T) {
		o := onpar.New(t)
		defer o.Run()

		sendName := func(t *testing.T) {
			c <- t.Name()
		}

		o.SerialSpec("A", sendName)

		o.Group("B", func() {
			o.SerialSpec("C", sendName)

			o.SerialSpec("D", sendName)
		})

		o.SerialSpec("E", sendName)

		o.Group("F", func() {
			o.Group("G", func() {
				o.SerialSpec("H", sendName)
			})
		})
	})
	close(c)

	expected := []string{
		"TestSerialSpecsAreOrdered/FakeSpecs/A",
		"TestSerialSpecsAreOrdered/FakeSpecs/B/C",
		"TestSerialSpecsAreOrdered/FakeSpecs/B/D",
		"TestSerialSpecsAreOrdered/FakeSpecs/E",
		"TestSerialSpecsAreOrdered/FakeSpecs/F/G/H",
	}

	i := 0
	for v := range c {
		if i >= len(expected) {
			t.Fatalf("only expected %d specs, but there are %d unexpected extra calls", len(expected), len(c))
		}
		thisExp := expected[i]
		if v != thisExp {
			t.Fatalf("expected run %d to be '%v'; got '%v'", i, thisExp, v)
		}
		i++
	}
	if i < len(expected) {
		t.Fatalf("expected %d specs, but there were only %d calls", len(expected), i)
	}
}

func TestUsingSuiteOutsideGroupPanics(t *testing.T) {
	var r any

	t.Run("FakeSpecs", func(t *testing.T) {
		defer func() {
			r = recover()
		}()

		o := onpar.New(t)
		defer o.Run()

		o.Group("foo", func() {
			// The most likely scenario for a suite accidentally being used
			// outside of its group is if it is reassigned to o. This seems
			// unlikely, since the types usually won't match (the setup
			// function's parameter and return types have to exactly match the
			// parent suite's) - but it's worth ensuring that this panics.
			o = onpar.BeforeEach(o, func(t *testing.T) *testing.T {
				return t
			})

			o.Spec("bar", func(*testing.T) {})
		})

		o.Spec("baz", func(*testing.T) {})
	})

	if r == nil {
		t.Fatalf("expected adding a spec to a *OnPar value outside of its group to panic")
	}
}

func TestUsingParentWithoutGroupPanics(t *testing.T) {
	var r any

	t.Run("FakeSpecs", func(t *testing.T) {
		defer func() {
			r = recover()
		}()

		o := onpar.New(t)
		defer o.Run()

		o.Group("foo", func() {
			b := onpar.BeforeEach(o, func(t *testing.T) string {
				return "foo"
			})

			onpar.BeforeEach(b, func(string) int {
				return 1
			})
		})
	})

	if r == nil {
		t.Fatalf("expected creating a child suite on a parent without a group to panic")
	}
}

func createScaffolding(t *testing.T) <-chan *testObject {
	objs := make(chan *testObject, 100)

	t.Run("FakeSpecs", func(t *testing.T) {
		o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) mockTest {
			obj := NewTestObject()
			obj.Use("-BeforeEach")

			objs <- obj

			return mockTest{t, 99, "something", obj}
		})
		defer o.Run()

		o.AfterEach(func(tt mockTest) {
			tt.o.Use("-AfterEach")
		})

		o.Group("DA", func() {
			o := onpar.BeforeEach(o, func(tt mockTest) mockTest {
				return tt
			})

			o.AfterEach(func(tt mockTest) {
				if tt.i != 99 {
					tt.t.Fatalf("expected %d = %d", tt.i, 99)
				}

				if tt.s != "something" {
					tt.t.Fatalf("expected %s = %s", tt.s, "something")
				}

				tt.o.Use("DA-AfterEach")
				tt.o.Close()
			})

			o.Spec("A", func(tt mockTest) {
				if tt.i != 99 {
					tt.t.Fatalf("expected %d = %d", tt.i, 99)
				}

				if tt.s != "something" {
					tt.t.Fatalf("expected %s = %s", tt.s, "something")
				}

				tt.o.Use("DA-A")
			})

			o.Group("DB", func() {
				type subMockTest struct {
					t *testing.T
					i int
					s string
					o *testObject
					f float64
				}

				o := onpar.BeforeEach(o, func(tt mockTest) subMockTest {
					tt.o.Use("DB-BeforeEach")
					return subMockTest{t: tt.t, i: tt.i, s: tt.s, o: tt.o, f: 101}
				})

				o.AfterEach(func(tt subMockTest) {
					tt.o.Use("DB-AfterEach")
				})

				o.Spec("B", func(tt subMockTest) {
					tt.o.Use("DB-B")
					if tt.i != 99 {
						tt.t.Fatalf("expected %d = %d", tt.i, 99)
					}

					if tt.s != "something" {
						tt.t.Fatalf("expected %s = %s", tt.s, "something")
					}

					if tt.f != 101 {
						tt.t.Fatalf("expected %f = %f", tt.f, 101.0)
					}
				})

				o.Spec("C", func(tt subMockTest) {
					tt.o.Use("DB-C")
					if tt.i != 99 {
						tt.t.Fatalf("expected %d = %d", tt.i, 99)
					}

					if tt.s != "something" {
						tt.t.Fatalf("expected %s = %s", tt.s, "something")
					}

					if tt.f != 101 {
						tt.t.Fatalf("expected %f = %f", tt.f, 101.0)
					}
				})

				o.Group("DDD", func() {
					type subSubMockTest struct {
						o *testObject
						i int
						t *testing.T
					}
					o := onpar.BeforeEach(o, func(tt subMockTest) subSubMockTest {
						tt.o.Use("DDD-BeforeEach")
						return subSubMockTest{o: tt.o, i: tt.i, t: tt.t}
					})

					o.AfterEach(func(tt subSubMockTest) {
						tt.o.Use("DDD-AfterEach")
					})

					o.Spec("E", func(tt subSubMockTest) {
						tt.o.Use("DDD-E")
						if tt.i != 99 {
							tt.t.Fatalf("expected %d = %d", tt.i, 99)
						}
					})
				})
			})
		})
	})

	return objs
}

func chanToSlice(c <-chan *testObject) []*testObject {
	var results []*testObject
	l := len(c)
	for i := 0; i < l; i++ {
		results = append(results, <-c)
	}
	return results
}

func findSpec(objs []*testObject, name string) *testObject {
	for _, obj := range objs {
		for _, specName := range obj.c {
			if name == specName {
				return obj
			}
		}
	}
	return nil
}

type testObject struct {
	c    []string
	done bool
}

func NewTestObject() *testObject {
	return &testObject{
		c: make([]string, 0),
	}
}

func (t *testObject) Use(i string) {
	t.c = append(t.c, i)
}

func (t *testObject) Close() {
	if t.done {
		panic("close() called too many times")
	}
	t.done = true
}

type mockTest struct {
	t *testing.T
	i int
	s string
	o *testObject
}
