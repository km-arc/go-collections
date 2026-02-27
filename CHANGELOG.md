# Changelog

All notable changes to this project will be documented in this file.
The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [v1.0.0] — 2024-01-01

### Added

**Collection[T]** — full Laravel list-collection API:
- Constructors: New, Collect, Times
- Transformation: Map, MapWithIndex, Filter, Reject, FlatMap, Flatten, Collapse, Reverse, Shuffle, Sort, SortDesc, Unique, UniqueBy, Duplicate, Pad, Transform (mutable)
- Slicing: Slice, Take, TakeUntil, TakeWhile, Skip, SkipUntil, SkipWhile, Chunk, Split, ForPage, Nth
- Searching: First, FirstOrFail, Last, Search, Contains, DoesntContain, Every, Has, IsEmpty, IsNotEmpty
- Aggregates: Count, Sum, Avg/Average, Min, Max, Median, Mode, Reduce
- Grouping: GroupBy, Partition
- Set ops: Diff, Intersect, Merge, Concat, Zip
- Mutable: Push, Prepend, Put, Forget, Pop, Shift, Pull, Splice
- Conditionals: When, Unless, WhenEmpty, WhenNotEmpty, UnlessEmpty, UnlessNotEmpty
- Utilities: Each, EachWithIndex, Tap, Pipe, Clone, Keys, Values, All, ToSlice, ToJSON, Implode, ImplodeWith, Random

**MapCollection[K, V]** — full Laravel associative-collection API:
- NewMap, Get, Put, Forget, Has, Keys, Values
- Filter, Reject, Map, Each, Every, Contains, Only, Except
- Merge, Union, Diff, DiffKeys, Intersect, IntersectByKeys, ToJSON
- Package functions: KeyBy, Pluck, Combine, Flip, MapKeys

**Infrastructure:**
- GitHub Actions CI: Go 1.21/1.22/1.23 matrix, race detector, coverage, benchmarks, golangci-lint, tidy check
- Release workflow: auto-creates GitHub Release on vX.Y.Z tags
- Codecov with 80% project / 70% patch thresholds
- Benchmark suite at 100 / 10k / 100k item scales
- Runnable GoDoc examples (all verified by go test)
- Package-level doc.go

[v1.0.0]: https://github.com/km-arc/go-collections/releases/tag/v1.0.0
