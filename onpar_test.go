package onpar_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar/v2"
)

func TestSingleNestedSpec(t *testing.T) {
	t.Parallel()
	o, c := createScaffolding()

	t.Run("FakeSpecs", o.Run)
	objs := chanToSlice(c)

	if len(objs) != 4 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objs), 3)
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
	o, c := createScaffolding()

	t.Run("FakeSpecs", o.Run)
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
	o, c := createScaffolding()

	t.Run("FakeSpecs", o.Run)
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
	o := onpar.New()

	c := make(chan string, 100)

	o.Spec("it runs a spec without a beforeeach", func(*testing.T) {
		c <- "A"
	})

	b := onpar.BeforeEach(o, func(*testing.T) string {
		c <- "B-BeforeEach"
		return "foo"
	})

	b.Spec("it runs a spec on a BeforeEach", func(string) {
		c <- "B"
	})

	t.Run("FakeSpecs", o.Run)

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
	o := onpar.New()

	c := make(chan string, 100)

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

	t.Run("FakeSpecs", o.Run)

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

func createScaffolding() (*onpar.Onpar[*testing.T, mockTest], <-chan *testObject) {
	objs := make(chan *testObject, 100)

	o := onpar.BeforeEach(onpar.New(), func(t *testing.T) mockTest {
		obj := NewTestObject()
		obj.Use("-BeforeEach")

		objs <- obj

		return mockTest{t, 99, "something", obj}
	})

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

	return o, objs
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
