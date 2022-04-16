# OnPar Matchers

OnPar provides a set of minimalistic matchers to get you started.
However, the intention is to be able to write your own custom matchers so that
your code is more readable.

## Generics

Most of the matchers provided by OnPar use generics to hopefully handle most
cases. This is new as of onpar v2. If you need to use any old matchers (from
v1), you can pass them to `matcher.Reflect`.

You will see in some examples below why we encourage users to build their own
matchers. Trying to make matchers generic ends up requiring a lot of type
parameters which often cannot be inferred. Most of the time, matchers built for
a specific enough use-case will only support one type, so they won't need any
type parameters. `BeTrue` and `HaveOccurred` are good examples of matchers that
have clear and specific use cases.

## Matchers List
- [String Matchers](#string-matchers)
- [Logical Matchers](#logical-matchers)
- [Error Matchers](#error-matchers)
- [Channel Matchers](#channel-matchers)
- [Collection Matchers](#collection-matchers)
- [Other Matchers](#other-matchers)


## String Matchers
### StartWith
StartWithMatcher accepts a string and succeeds if the actual string starts with
the expected string.

```go
Expect(t, "foobar").To(StartWith("foo"))
```

### EndWith
EndWithMatcher accepts a string and succeeds if the actual string ends with
the expected string.

```go
Expect(t, "foobar").To(EndWith("bar"))
```

### ContainSubstring
ContainSubstringMatcher accepts a string and succeeds if the expected string is a
sub-string of the actual.

```go
Expect(t, "foobar").To(ContainSubstring("ooba"))
```
### MatchRegexp

## Logical Matchers
### Not
NotMatcher accepts a matcher and will succeed if the specified matcher fails.

```go
Expect(t, false).To(Not[bool](BeTrue()))
```

### BeAbove
BeAboveMatcher accepts a float64. It succeeds if the actual is greater
than the expected.

```go
Expect(t, 100).To(BeAbove(99))
```

### BeBelow
BeBelowMatcher accepts a float64. It succeeds if the actual is
less than the expected.

```go
Expect(t, 100).To(BeBelow(101))
```

### BeFalse
BeFalseMatcher will succeed if actual is false.

```go
Expect(t, 2 == 3).To(BeFalse())
```

### BeTrue
BeTrueMatcher will succeed if actual is true.

```go
Expect(t, 2 == 2).To(BeTrue())
```

### Equal
EqualMatcher performs a DeepEqual between the actual and expected.

```go
Expect(t, 42).To(Equal(42))
```

## Error Matchers
### HaveOccurred
HaveOccurredMatcher will succeed if the actual value is a non-nil error.

```go
Expect(t, err).To(HaveOccurred())

Expect(t, nil).To(Not[error](HaveOccurred()))
```

## Channel Matchers
### Receive
ReceiveMatcher will attempt to receive from the channel but will not block. It
fails if the channel is closed.

```go
c := make(chan bool, 1)
c <- true
Expect(t, c).To(Receive[chan bool, bool](Anything[bool]()))
```

If your use case requires waiting on the channel, we provide a `ReceiveWait`
option to add a timeout to the select instead of a default:

``` go
c := make(chan bool, 1)
go func() {
  time.Sleep(time.Millisecond)
  c <- true
}
Expect(t, c).To(Receive[chan bool, bool](BeTrue(), ReceiveWait(time.Second)))
```

### BeClosed
BeClosedMatcher succeeds if the channel is closed.

```go
c := make(chan bool)
close(c)
Expect(t, c).To(BeClosed[chan bool, bool]())
```

## Collection Matchers
### HaveCap
This matcher works on Slices and Channels and will succeed if the type has the
specified capacity.

```go
Expect(t, []string{"foo", "bar"}).To(HaveCap[[]string, string](2))
```
### HaveKey
HaveKeyMatcher accepts map types and will succeed if the map contains the
specified key.

```go
fooMap := map[string]int{"foo":69}
Expect(t, fooMap).To(HaveKey[map[string]int, string, int]("foo", WithValue(69)))
```
## Other Matchers

### Always
AlwaysMatcher matches by polling the child matcher until it returns an error. It
will return an error the first time the child matcher returns an error. If the
child matcher never returns an error, then it will return a nil.

By default, the duration is 100ms with an interval of 10ms.

```go
isTrue := func() bool {
  return true
}
Expect(t, isTrue).To(Always[func() bool, bool](BeTrue()))
```

### Eventually
EventuallyMatcher matches by polling, similar to Always; however, Eventually
will poll until the child matcher _stops_ returning an error. The first time
that the child matcher returns success, the EventuallyMatcher will return
success.

``` go
times := 0
isTrue := func() bool {
  times++
  if times == 10 {
    return true
  }
  return false
}
Expect(t, isTrue).To(Eventually[func() bool, bool](BeTrue()))
```
