package onpar_test

import (
	"reflect"
	"testing"

	"github.com/poy/onpar"
)

func TestSingleNestedSpec(t *testing.T) {
	t.Parallel()
	o, c := createScaffolding()

	t.Run("FakeSpecs", func(t *testing.T) {
		o.Run(t)
	})
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

	t.Run("FakeSpecs", func(t *testing.T) {
		o.Run(t)
	})
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

	t.Run("FakeSpecs", func(t *testing.T) {
		o.Run(t)
	})
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

func TestDoNotInvokeStrandedBeforeEach(t *testing.T) {
	t.Parallel()
	o, c := createScaffolding()

	t.Run("FakeSpecs", func(t *testing.T) {
		o.Run(t)
	})
	objs := chanToSlice(c)

	objCBeforeEach := findSpec(objs, "DC-BeforeEach")
	if objCBeforeEach != nil {
		t.Fatal("should not have invoked BeforeEach")
	}
}

func TestDoNotInvokeStrandedAfterEach(t *testing.T) {
	t.Parallel()
	o, c := createScaffolding()

	t.Run("FakeSpecs", func(t *testing.T) {
		o.Run(t)
	})
	objs := chanToSlice(c)

	objCBeforeEach := findSpec(objs, "DC-AfterEach")
	if objCBeforeEach != nil {
		t.Fatal("should not have invoked AfterEach")
	}
}

func createScaffolding() (*onpar.Onpar, <-chan *testObject) {
	o := onpar.New()
	objs := make(chan *testObject, 100)

	o.BeforeEach(func(t *testing.T) (*testing.T, int, string, *testObject) {
		obj := NewTestObject()
		obj.Use("-BeforeEach")

		objs <- obj

		return t, 99, "something", obj
	})

	o.AfterEach(func(t *testing.T, i int, s string, o *testObject) {
		o.Use("-AfterEach")
	})

	o.Group("DA", func() {
		o.AfterEach(func(t *testing.T, i int, s string, o *testObject) {
			if i != 99 {
				t.Fatalf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Fatalf("expected %s = %s", s, "something")
			}

			o.Use("DA-AfterEach")
			o.Close()
		})

		o.Spec("A", func(t *testing.T, i int, s string, o *testObject) {
			if i != 99 {
				t.Fatalf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Fatalf("expected %s = %s", s, "something")
			}

			o.Use("DA-A")
		})

		o.Group("DB", func() {
			o.BeforeEach(func(t *testing.T, i int, s string, o *testObject) (*testing.T, int, string, *testObject, float64) {
				o.Use("DB-BeforeEach")
				return t, i, s, o, 101
			})

			o.AfterEach(func(t *testing.T, i int, s string, o *testObject, f float64) {
				o.Use("DB-AfterEach")
			})

			o.Spec("B", func(t *testing.T, i int, s string, o *testObject, f float64) {
				o.Use("DB-B")
				if i != 99 {
					t.Fatalf("expected %d = %d", i, 99)
				}

				if s != "something" {
					t.Fatalf("expected %s = %s", s, "something")
				}

				if f != 101 {
					t.Fatalf("expected %f = %f", f, 101.0)
				}
			})

			o.Spec("C", func(t *testing.T, i int, s string, o *testObject, f float64) {
				o.Use("DB-C")
				if i != 99 {
					t.Fatalf("expected %d = %d", i, 99)
				}

				if s != "something" {
					t.Fatalf("expected %s = %s", s, "something")
				}

				if f != 101 {
					t.Fatalf("expected %f = %f", f, 101.0)
				}
			})

			o.Group("DDD", func() {
				o.BeforeEach(func(t *testing.T, i int, s string, o *testObject, f float64) (*testObject, int, *testing.T) {
					o.Use("DDD-BeforeEach")
					return o, i, t
				})

				o.AfterEach(func(o *testObject, i int, t *testing.T) {
					o.Use("DDD-AfterEach")
				})

				o.Spec("E", func(o *testObject, i int, t *testing.T) {
					o.Use("DDD-E")
					if i != 99 {
						t.Fatalf("expected %d = %d", i, 99)
					}
				})
			})
		})

		o.Group("DC", func() {
			o.BeforeEach(func(t *testing.T, i int, s string, o *testObject) {
				o.Use("DC-BeforeEach")
				t.Fatalf("should not have been invoked")
			})

			o.AfterEach(func(t *testing.T, i int, s string, o *testObject) {
				o.Use("DC-AfterEach")
				t.Fatalf("should not have been invoked")
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
