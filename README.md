# onpar
Parallel testing framework for Go

### Specs
Test assertions are done within a `Spec()` function. Each `Spec` has a name and a function. The function takes a `testing.T` as an argument and any output from a `BeforeEach()`. Each `Spec` is run in parallel (`t.Parallel()` is invoked for each spec before calling the given function).

```golang
o := onpar.New()

o.BeforeEach(func(t *testing.T) (a int, b float64) {
    return 99, 101.0
})

o.AfterEach(func(t *testing.T, a int, b float64) {
        // ...
})

o.Spec("something informative", func(t *testing.T, a int, b float64) {
    if a != 99 {
        t.Errorf("%d != 99", a)
    }
})
```

### Grouping
`Group`s are used to keep `Spec`s in logical place. The intention is to gather each `Spec` in a reasonable place. Each `Group` can have a `BeforeEach()` and a `AfterEach()` but are not required to.


```golang
o := onpar.New()

o.BeforeEach(func(t *testing.T) (a int, b float64) {
    return 99, 101.0
})

o.Group("some-group", func() {
    o.BeforeEach(func(t *teting.T, a int, b float64) (s string) {
        return "foo"
    })

    o.AfterEach(func(t *teting.T, a int, b float64, s string) {
        // ...
    })
    
    o.Spec("something informative", func(t *testing.T, a int, b float64, s string) {
        // ...
    })
})
```

### Run Order
Each `BeforeEach()` runs before any `Spec` in the same `Group`. It will also run before any sub-group `Spec`s and their `BeforeEach`es. Any `AfterEach()` will run after the `Spec` and before parent `AfterEach`es.

### Avoiding Closure
Why bother with returning values from a `BeforeEach`? To avoid closure of course! When running `Spec`s in parallel (which they always do), each variable needs a new instance to avoid race conditions. If you use closure, then this gets tough. So onpar will pass the arguments to the given function.
