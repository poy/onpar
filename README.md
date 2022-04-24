# onpar
[![GoDoc][go-doc-badge]][go-doc] [![travis][travis-badge]][travis]

Parallel testing framework for Go

## Goals

- Provide structured testing, with per-spec setup and teardown.
- Discourage using closure state to share memory between setup/spec/teardown
  functions.
  - Sharing memory between the steps of a spec by using closure state means that
    you're also sharing memory with _other tests_. This often results in test
    pollution.
- Run tests in parallel by default.
  - Most of the time, well-written unit tests are perfectly capable of running
    in parallel, and sometimes running tests in parallel can uncover extra bugs.
    This should be the default.
- Work within standard go test functions, simply wrapping standard `t.Run`
  semantics.
  - `onpar` should not feel utterly alien to people used to standard go testing.
    It does some extra work to allow structured tests, but for the most part it
    isn't hiding any complicated logic - it mostly just calls `t.Run`.

Onpar provides a BDD style of testing, similar to what you might find with
something like ginkgo or goconvey. The biggest difference between onpar and its
peers is that a `BeforeEach` function in `onpar` may return a value, and that
value will become the parameter required in child calls to `Spec`, `AfterEach`,
and `BeforeEach`.

This allows you to write tests that share memory between `BeforeEach`, `Spec`,
and `AfterEach` functions _without sharing memory with other tests_. When used
properly, this makes test pollution nearly impossible and makes it harder to
write flaky tests.

## Running

In order for specs to actually run, `o.Run(t)` must be called at the end of a
test function. Typically, we do this with `defer o.Run(t)` immediately after the
top level suite is constructed.

`o.Run(t)` should never be called on child suites constructed with
`BeforeEach()`.

### Assertions
OnPar provides built-in assertion mechanisms using `Expect` statement to make
expectations. Here is some more information about `Expect` and some of the
matchers that come with OnPar.

- [Expect](expect/README.md)
- [Matchers](matchers/README.md)

### Specs
Test assertions are done within a `Spec()` function. Each `Spec` has a name and
a function with a single argument. The type of the argument is determined by how
the suite was constructed: `New()` returns a suite that takes a `*testing.T`,
while `BeforeEach` constructs a suite that takes the return type of the setup
function.

Each `Spec` is run in parallel (`t.Parallel()` is invoked for each spec before
calling the given function).

```go
func TestSpecs(t *testing.T) {
    type testContext struct {
        t *testing.T
        a int
        b float64
    }

    o := onpar.BeforeEach(onpar.New(), func(t *testing.T) testContext {
        return testContext{t: t, a: 99, b: 101.0}
    })
    defer o.Run(t)

    o.AfterEach(func(tt testContext) {
            // ...
    })

    o.Spec("something informative", func(tt testContext) {
        if tt.a != 99 {
            tt.t.Errorf("%d != 99", tt.a)
        }
    })
}
```

### Grouping
`Group`s are used to keep `Spec`s in logical place. The intention is to gather
each `Spec` in a reasonable place. Each `Group` may construct a new child suite
using `BeforeEach`.


```go
func TestGrouping(t *testing.T) {
    type topContext struct {
        t *testing.T
        a int
        b float64
    }

    o := onpar.BeforeEach(onpar.New(), func(t *testing.T) topContext {
        return topContext{t: t, a: 99, b: 101}
    }
    defer o.Run(t)

    o.Group("some-group", func() {
        type groupContext struct {
            t *testing.T
            s string
        }
        o := onpar.BeforeEach(o, func(tt topContext) groupContext {
            return groupContext{t: t, s: "foo"}
        })

        o.AfterEach(func(tt groupContext) {
            // ...
        })

        o.Spec("something informative", func(tt groupContext) {
            // ...
        })
    })
}
```

### Run Order
Each `BeforeEach()` runs before any `Spec` in the same `Group`. It will also run
before any sub-group `Spec`s and their `BeforeEach`es. Any `AfterEach()` will
run after the `Spec` and before parent `AfterEach`es.

``` go
func TestRunOrder(t *testing.T) {
    type topContext struct {
        t *testing.T
        i int
        s string
    }
    o := onpar.BeforeEach(onpar.New(), func(t *testing.T) topContext {
        // Spec "A": Order = 1
        // Spec "B": Order = 1
        // Spec "C": Order = 1
        return t, 99, "foo"
    })
    defer o.Run(t)

    o.AfterEach(func(tt topContext) {
        // Spec "A": Order = 4
        // Spec "B": Order = 6
        // Spec "C": Order = 6
    })

    o.Group("DA", func() {
        o.AfterEach(func(tt topContext) {
            // Spec "A": Order = 3
            // Spec "B": Order = 5
            // Spec "C": Order = 5
        })

        o.Spec("A", func(tt topContext) {
            // Spec "A": Order = 2
        })

        o.Group("DB", func() {
            type dbContext struct {
                t *testing.T
                f float64
            }
            o.BeforeEach(func(tt topContext) dbContext {
                // Spec "B": Order = 2
                // Spec "C": Order = 2
                return dbContext{t: t, f: 101}
            })

            o.AfterEach(func(tt dbContext) {
                // Spec "B": Order = 4
                // Spec "C": Order = 4
            })

            o.Spec("B", func(tt dbContext) {
                // Spec "B": Order = 3
            })

            o.Spec("C", func(tt dbContext) {
                // Spec "C": Order = 3
            })
        })

        o.Group("DC", func() {
            o.BeforeEach(func(tt topContext) *testing.T {
                // Will not be invoked (there are no specs)
            })

            o.AfterEach(func(t *testing.T) {
                // Will not be invoked (there are no specs)
            })
        })
    })
}
```

## Avoiding Closure
Why bother with returning values from a `BeforeEach`? To avoid closure of
course! When running `Spec`s in parallel (which they always do), each variable
needs a new instance to avoid race conditions. If you use closure, then this
gets tough. So onpar will pass the arguments to the given function returned by
the `BeforeEach`.

The `BeforeEach` is a gatekeeper for arguments. The returned values from
`BeforeEach` are required for the following `Spec`s. Child `Group`s are also
passed what their direct parent `BeforeEach` returns.

[go-doc-badge]:             https://pkg.go.dev/github.com/poy/onpar?status.svg
[go-doc]:                   https://pkg.go.dev/github.com/poy/onpar
[travis-badge]:             https://travis-ci.org/poy/onpar.svg?branch=master
[travis]:                   https://travis-ci.org/poy/onpar?branch=master
