# collection

> Laravel-style fluent collections for Go — type-safe, generic, zero-dependency.

[![CI](https://github.com/km-arc/go-collections/actions/workflows/ci.yml/badge.svg)](https://github.com/km-arc/go-collections/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/km-arc/go-collections/branch/main/graph/badge.svg)](https://codecov.io/gh/km-arc/go-collections)
[![Go Reference](https://pkg.go.dev/badge/github.com/km-arc/go-collections.svg)](https://pkg.go.dev/github.com/km-arc/go-collections)
[![Go Report Card](https://goreportcard.com/badge/github.com/km-arc/go-collections)](https://goreportcard.com/report/github.com/km-arc/go-collections)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/km-arc/go-collections)](https://github.com/km-arc/go-collections/releases)

---

## Why?

Go's standard library gives you slices and maps. Laravel gives you 60+ fluent methods that let you express complex data transformations as readable pipelines.

This package brings the same expressive power to Go — with full type safety via generics, zero reflection, and zero external dependencies.

```go
result := collection.New(users).
    Filter(func(u User) bool { return u.Active }).
    Sort(func(a, b User) bool { return a.Name < b.Name }).
    ForPage(1, 20).
    All()
```

---

## Install

```bash
go get github.com/km-arc/go-collections@v1.0.0
```

**Requires Go 1.21+**

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/km-arc/go-collections"
)

func main() {
    result := collection.New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
        Filter(func(n int) bool { return n%2 == 0 }).
        Map(func(n int) int    { return n * n }).
        Take(3).
        All()

    fmt.Println(result) // [4 16 36]
}
```

---

## Two types

| Type | Analogy | Use for |
|---|---|---|
| `Collection[T]` | `collect([1,2,3])` | ordered sequences of any type |
| `MapCollection[K,V]` | `collect(['key' => val])` | key/value stores |

---

## API Reference

### Constructors

```go
collection.New([]T{...})
collection.Collect([]T{...})              // alias for New
collection.Times(n, func(i int) T {...})  // generate n items (1-based)
collection.NewMap(map[K]V{...})
```

### Transformation

| Method | Description |
|---|---|
| `Map(fn)` | Transform each item |
| `MapWithIndex(fn)` | Transform with 0-based index |
| `Filter(fn)` | Keep matching items |
| `Reject(fn)` | Discard matching items |
| `FlatMap(fn)` | Map then flatten one level |
| `Flatten(c)` *(func)* | Collapse `Collection[[]T]` to flat |
| `Reverse()` | Reverse order |
| `Unique()` / `UniqueBy(fn)` | Remove duplicates |
| `Duplicate()` | Items that appear more than once |
| `Shuffle()` | Random order |
| `Sort(less)` / `SortDesc(less)` | Stable sort |
| `Pad(size, val)` | Grow to size filling with val |
| `Transform(fn)` | **Mutable** in-place map |

### Slicing & Pagination

| Method | Description |
|---|---|
| `Slice(offset, len)` | Sub-collection |
| `Take(n)` | First n (negative = last n) |
| `TakeUntil(fn)` / `TakeWhile(fn)` | Conditional take |
| `Skip(n)` | Skip first n |
| `SkipUntil(fn)` / `SkipWhile(fn)` | Conditional skip |
| `Chunk(n)` | Split into n-sized groups |
| `Split(n)` | Divide into n equal groups |
| `ForPage(page, perPage)` | Paginate (1-based) |
| `Nth(n, offset)` | Every n-th item |

### Aggregates

```go
c.Count()
c.Sum(fn) / c.Avg(fn) / c.Average(fn)
c.Min(fn) → (float64, bool)
c.Max(fn) → (float64, bool)
c.Median(fn) → (float64, bool)
c.Mode(fn) → []float64
c.Reduce(initial, fn)
```

### Searching & Testing

```go
c.First(fn) → (T, bool)     // nil fn = first item
c.FirstOrFail(fn) → (T, error)
c.Last(fn) → (T, bool)
c.Search(fn) → int          // index or -1
c.Contains(fn) / c.DoesntContain(fn)
c.Every(fn)
c.Has(index)
c.IsEmpty() / c.IsNotEmpty()
```

### Grouping & Set Operations

```go
c.GroupBy(fn)   → map[string]*Collection[T]
c.Partition(fn) → (pass, fail *Collection[T])
c.Diff(other)
c.Intersect(other)
c.Merge(other)
c.Concat(others...)
collection.Zip(a, b)
```

### Mutable Operations

```go
c.Push(item) / c.Prepend(item)
c.Put(index, item) / c.Forget(index)
c.Pop() / c.Shift() / c.Pull(index)
c.Splice(offset, len, replacements...)
c.Transform(fn)
```

### Conditionals

```go
c.When(bool, fn) / c.Unless(bool, fn)
c.WhenEmpty(fn) / c.WhenNotEmpty(fn)
c.UnlessEmpty(fn) / c.UnlessNotEmpty(fn)
```

### Chaining Utilities

```go
c.Each(fn) / c.EachWithIndex(fn)
c.Tap(fn)   // side-effect, returns c unchanged
c.Pipe(fn)  // transform mid-chain
c.Clone()
c.Keys() / c.Values()
c.All() / c.ToSlice()
c.ToJSON()
c.Implode(glue) / c.ImplodeWith(glue, fn)
c.Random(n)
```

### MapCollection

```go
mc.Get(key) / mc.Put(key, val) / mc.Forget(key) / mc.Has(key)
mc.Keys() → *Collection[K]
mc.Values() → *Collection[V]
mc.Filter(fn) / mc.Reject(fn) / mc.Map(fn)
mc.Only(keys...) / mc.Except(keys...)
mc.Merge(other) / mc.Union(other)
mc.Diff(other) / mc.DiffKeys(other)
mc.Intersect(other) / mc.IntersectByKeys(other)
mc.ToJSON()

// Package-level:
collection.KeyBy(c, fn)         // Collection  → MapCollection
collection.Pluck(c, fn)         // extract field → Collection
collection.Combine(keys, vals)  // two Collections → MapCollection
collection.Flip(mc)             // swap keys & values
collection.MapKeys(mc, fn)      // rekey a MapCollection
```

---

## Benchmarks

Run on Apple M2, Go 1.22, `-benchmem -count=3`:

```
BenchmarkNew_10k-8              384,210 ns/op    81,920 B/op    1 allocs/op
BenchmarkMap_10k-8              274,831 ns/op    81,920 B/op    1 allocs/op
BenchmarkMap_100k-8           2,861,792 ns/op   802,816 B/op    1 allocs/op
BenchmarkFilter_10k-8           163,402 ns/op    45,232 B/op    8 allocs/op
BenchmarkReduce_10k-8            48,961 ns/op         0 B/op    0 allocs/op
BenchmarkSum_10k-8               48,714 ns/op         0 B/op    0 allocs/op
BenchmarkSort_10k-8             901,248 ns/op    81,920 B/op    1 allocs/op
BenchmarkGroupBy_10k-8          712,399 ns/op   453,101 B/op  304 allocs/op
BenchmarkPipeline_10k-8       1,012,741 ns/op   435,720 B/op   17 allocs/op
BenchmarkKeyBy_10k-8            278,102 ns/op   478,432 B/op   13 allocs/op
BenchmarkPluck_10k-8            261,874 ns/op   163,840 B/op    1 allocs/op
```

```bash
go test -bench=. -benchmem ./...
```

---

## Testing

```bash
go test -race ./...                                        # all tests
go test -race -coverprofile=coverage.out ./...            # with coverage
go tool cover -html=coverage.out                          # view in browser
go test -bench=. -benchmem ./...                          # benchmarks
go test -run TestChaining -v ./...                        # single test
```

---

## Differences from Laravel

| PHP/Laravel | This package |
|---|---|
| `collect([...])` | `collection.New([]T{...})` |
| `->first()` → `null` | Returns `(T, bool)` — no nil surprises |
| `->avg('field')` | `c.Avg(func(v T) float64 { return v.Field })` |
| Dynamic `$key => $value` | `MapCollection[K, V]` — fully typed |
| `Collection::macro(...)` | Embed `Collection[T]` in your own struct |

---

## Contributing

1. Add tests for new behaviour (coverage must stay ≥ 80%).
2. Add GoDoc comments on all exported symbols.
3. Run `go vet ./...` and `golangci-lint run` before opening a PR.
4. CI must be green (test, lint, tidy, bench).

---

## License

MIT © 2024
