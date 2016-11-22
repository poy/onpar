# OnPar Matchers

OnPar provides a set of minimalistic matchers to get you started.
However, the intention is to be able to write your own custom matchers so that
your code is more readable.


## Matchers List
- [String Matchers](#string-matchers)
- [Logical Matchers](#logical-matchers)
- [Error Matchers](#error-matchers)
- [Channel Matchers](#channel-matchers)
- [Collection Matchers](#collection-matchers)
- [Other Matchers](#other-matchers)


## String Matchers
### StartsWith
```go
Expect(t, "foobar").To(StartsWith("foo"))
```
### EndsWith
```go
Expect(t, "foobar").To(EndsWith("bar"))
```
### Contains
```go
Expect(t, "foobar").To(Contains("ooba"))
```

## Logical Matchers
### And
### Or
### Not
```go
Expect(t, false).To(Not(BeTrue()))
```
### BeAbove
```go
Expect(t, 100).To(BeAbove(99))
```
### BeBelow
```go
Expect(t, 100).To(BeBelow(101))
```
### BeFalse
```go
Expect(t, 2 == 3).To(BeFalse())
```
### BeTrue
```go
Expect(t, 2 == 2).To(BeTrue())
```
### Equal
```go
Expect(t, 42).To(BeEqual(42))
```

## Error Matchers
### HaveOccurred
```go
Expect(t, err).To(HaveOccurred())

Expect(t, nil).To(Not(HaveOccurred()))
```

## Channel Matchers
### Always
### Receive

## Collection Matchers
### HaveCap
This matcher works on Slices, Arrays, Maps and Channels.

```go
Expect(t, []string{"foo", "bar"}).To(HaveCap(2))
```
### HaveKey
This matcher works on Maps.

```go
Expect(t, fooMap).To(HaveKey("foo"))
```

### HaveLen
This matcher works on Strings, Slices, Arrays, Maps and Channels.
```go
Expect(t, "12345").To(HaveLen(5))
```

## Other Matchers
### Chain
### ViaPolling
### MatchRegexp
