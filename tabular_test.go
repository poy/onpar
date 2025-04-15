package onpar_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/poy/onpar"
)

func TestTableSpec_Entry(t *testing.T) {
	t.Parallel()
	c := createTableScaffolding(t)

	objs := chanToSlice(c)

	objA := findSpec(objs, "DA-A")
	if objA == nil {
		t.Fatal("unable to find spec A")
	}

	if len(objA.c) != 3 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objA.c), 3)
	}

	if !reflect.DeepEqual(objA.c, []string{"-BeforeEach", "DA-A", "-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objA.c)
	}
}

func ExampleSpecTable_Entry() {
	var t *testing.T
	o := onpar.New(t)
	defer o.Run()

	type table struct {
		input          string
		expectedOutput string
	}
	f := func(in string) string {
		return in + "world"
	}
	onpar.TableSpec(o, func(t *testing.T, tt table) {
		output := f(tt.input)
		if output != tt.expectedOutput {
			t.Fatalf("expected %v to produce %v; got %v", tt.input, tt.expectedOutput, output)
		}
	}).
		Entry("simple output", table{"hello", "helloworld"}).
		Entry("with a space", table{"hello ", "hello world"}).
		Entry("and a comma", table{"hello, ", "hello, world"})
}

func TestTableSpec_FnEntry(t *testing.T) {
	t.Parallel()
	c := createTableScaffolding(t)

	objs := chanToSlice(c)

	objB := findSpec(objs, "DA-B")
	if objB == nil {
		t.Fatal("unable to find spec B")
	}

	if len(objB.c) != 4 {
		t.Fatalf("expected objs (len=%d) to have len %d", len(objB.c), 4)
	}

	if !reflect.DeepEqual(objB.c, []string{"-BeforeEach", "-TableSpecEntrySetup", "DA-B", "-AfterEach"}) {
		t.Fatalf("invalid call order for spec A: %v", objB.c)
	}
}

func ExampleSpecTable_FnEntry() {
	var t *testing.T
	o := onpar.New(t)
	defer o.Run()

	type table struct {
		input          string
		expectedOutput string
	}
	f := func(in string) string {
		return in + "world"
	}
	onpar.TableSpec(o, func(t *testing.T, tt table) {
		output := f(tt.input)
		if output != tt.expectedOutput {
			t.Fatalf("expected %v to produce %v; got %v", tt.input, tt.expectedOutput, output)
		}
	}).
		FnEntry("simple output", func(t *testing.T) table {
			var buf bytes.Buffer
			if _, err := buf.WriteString("hello"); err != nil {
				t.Fatalf("expected buffer write to succeed; got %v", err)
			}
			return table{input: buf.String(), expectedOutput: "helloworld"}
		})
}

func TestTableGroup_Entry(t *testing.T) {
	t.Parallel()
	c := createTableScaffolding(t)

	objs := chanToSlice(c)

	t.Run("DB-A-A", func(t *testing.T) {
		objA := findSpec(objs, "DB-A-A")
		if objA == nil {
			t.Fatal("unable to find spec DB-A-A")
		}

		if len(objA.c) != 3 {
			t.Fatalf("expected objs (len=%d) to have len %d", len(objA.c), 3)
		}

		if !reflect.DeepEqual(objA.c, []string{"-BeforeEach", "DB-A-A", "-AfterEach"}) {
			t.Fatalf("invalid call order for group A: %v", objA.c)
		}
	})

	t.Run("DB-A-B", func(t *testing.T) {
		objB := findSpec(objs, "DB-A-B")
		if objB == nil {
			t.Fatal("unable to find spec DB-A-B")
		}

		if len(objB.c) != 3 {
			t.Fatalf("expected objs (len=%d) to have len %d", len(objB.c), 3)
		}

		if !reflect.DeepEqual(objB.c, []string{"-BeforeEach", "DB-A-B", "-AfterEach"}) {
			t.Fatalf("invalid call order for group B: %v", objB.c)
		}
	})
}

func ExampleGroupTable_Entry() {
	var t *testing.T
	o := onpar.New(t)
	defer o.Run()

	type table struct {
		input          string
		expectedOutput string
	}
	f := func(in string) string {
		return in + "world"
	}
	onpar.TableGroup(o, func(tt table) {
		o.Spec("produces expected output", func(t *testing.T) {
			output := f(tt.input)
			if output != tt.expectedOutput {
				t.Fatalf("expected %v to produce %v; got %v", tt.input, tt.expectedOutput, output)
			}
		})

		o.Spec("does not contain fakePersonallyIdentifiableInfo", func(t *testing.T) {
			output := f(tt.input)
			if strings.Contains(output, "fakePersonallyIdentifiableInfo") {
				t.Fatalf("expected %v not to contain fakePersonallyIdentifiableInfo", output)
			}
		})
	}).
		Entry("simple output", table{"hello", "helloworld"}).
		Entry("with a space", table{"hello ", "hello world"}).
		Entry("and a comma", table{"hello, ", "hello, world"})
}

func createTableScaffolding(t *testing.T) <-chan *testObject {
	objs := make(chan *testObject, 100)

	t.Run("FakeSpecs", func(t *testing.T) {
		o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) *mockTest {
			obj := NewTestObject()
			obj.Use("-BeforeEach")

			objs <- obj

			return &mockTest{t, 99, "something", obj}
		})
		defer o.Run()

		o.AfterEach(func(tt *mockTest) {
			tt.o.Use("-AfterEach")
		})

		type table struct {
			name     string
			expected mockTest
		}

		onpar.TableSpec(o, func(tt *mockTest, tab table) {
			if tt.i != tab.expected.i {
				tt.t.Fatalf("expected %d = %d", tt.i, tab.expected.i)
			}

			if tt.s != tab.expected.s {
				tt.t.Fatalf("expected %s = %s", tt.s, tab.expected.s)
			}

			tt.o.Use(tab.name)
		}).
			Entry("DA-A", table{name: "DA-A", expected: mockTest{i: 99, s: "something"}}).
			FnEntry("DA-B", func(tt *mockTest) table {
				tt.i = 21
				tt.s = "foo"
				tt.o.Use("-TableSpecEntrySetup")
				return table{name: "DA-B", expected: mockTest{i: 21, s: "foo"}}
			})

		onpar.TableGroup(o, func(tab table) {
			aName := fmt.Sprintf("%s-A", tab.name)
			o.Spec(aName, func(tt *mockTest) {
				if tt.i != tab.expected.i {
					tt.t.Fatalf("expected %d = %d", tt.i, tab.expected.i)
				}
				tt.o.Use(aName)
			})

			bName := fmt.Sprintf("%s-B", tab.name)
			o.Spec(bName, func(tt *mockTest) {
				if tt.s != tab.expected.s {
					tt.t.Fatalf("expected %s = %s", tt.s, tab.expected.s)
				}
				tt.o.Use(bName)
			})
		}).Entry("DB-A", table{name: "DB-A", expected: mockTest{i: 99, s: "something"}})
	})

	return objs
}
