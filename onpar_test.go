package onpar_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/apoydence/onpar"
)

func TestOrder(t *testing.T) {
	var lock sync.Mutex
	var objs []*testObject

	onpar.AfterEach(func(t *testing.T) {
	})

	onpar.Group("DA", func() {
		onpar.BeforeEach(func(t *testing.T) (int, string, *testObject) {
			obj := NewTestObject()
			obj.Use("DA-BeforeEach")

			lock.Lock()
			objs = append(objs, obj)
			lock.Unlock()

			return 99, "something", obj
		})

		onpar.AfterEach(func(t *testing.T, i int, s string, o *testObject) {
			if i != 99 {
				t.Fatalf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Fatalf("expected %s = %s", s, "something")
			}

			o.Use("DA-AfterEach")
			o.Close()
		})

		onpar.Spec("A", func(t *testing.T, i int, s string, o *testObject) {
			if i != 99 {
				t.Fatalf("expected %d = %d", i, 99)
			}

			if s != "something" {
				t.Fatalf("expected %s = %s", s, "something")
			}

			o.Use("DA-A")
		})

		onpar.Group("DB", func() {
			onpar.BeforeEach(func(t *testing.T, i int, s string, o *testObject) float64 {
				o.Use("DB-BeforeEach")
				return 101
			})

			onpar.AfterEach(func(t *testing.T, i int, s string, o *testObject, f float64) {
				o.Use("DB-AfterEach")
			})

			onpar.Spec("B", func(t *testing.T, i int, s string, o *testObject, f float64) {
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

			onpar.Spec("C", func(t *testing.T, i int, s string, o *testObject, f float64) {
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
		})

		onpar.Group("DC", func() {
			onpar.BeforeEach(func(t *testing.T, i int, s string, o *testObject) {
				o.Use("DC-BeforeEach")
			})
		})
	})

	t.Run("", func(tt *testing.T) {
		onpar.Run(tt)
	})

	if len(objs) != 3 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objs), 3)
	}

	objA := findSpec(objs, "DA-A")
	if objA == nil {
		t.Fatal("unable to find spec A")
	}

	if len(objA.c) != 3 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objA.c), 3)
	}

	if !reflect.DeepEqual(objA.c, []string{"DA-BeforeEach", "DA-A", "DA-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objA.c)
	}

	objB := findSpec(objs, "DB-B")
	if objB == nil {
		t.Fatal("unable to find spec B")
	}

	if len(objB.c) != 5 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objB.c), 5)
	}

	if !reflect.DeepEqual(objB.c, []string{"DA-BeforeEach", "DB-BeforeEach", "DB-B", "DB-AfterEach", "DA-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objB.c)
	}

	objC := findSpec(objs, "DB-C")
	if objC == nil {
		t.Fatal("unable to find spec C")
	}

	if len(objC.c) != 5 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objC.c), 5)
	}

	if !reflect.DeepEqual(objC.c, []string{"DA-BeforeEach", "DB-BeforeEach", "DB-C", "DB-AfterEach", "DA-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objC.c)
	}

	objCBeforeEach := findSpec(objs, "DC-BeforeEach")
	if objCBeforeEach != nil {
		t.Fatal("should not have invoked before each")
	}
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
