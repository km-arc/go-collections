// Package collection provides a type-safe, generic implementation of
// [Laravel's Illuminate\Support\Collection] for Go.
//
// # Overview
//
// Laravel collections are one of the most productive APIs in modern web
// development: a single fluent interface for filtering, mapping, sorting,
// grouping, and aggregating data without verbose loops.  This package brings
// that same experience to Go using the generics introduced in Go 1.18.
//
//	import "github.com/km-arc/go-collections"
//
//	// Build a typed collection from any slice.
//	nums := collection.New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//
//	result := nums.
//	    Filter(func(n int) bool { return n%2 == 0 }).  // keep even
//	    Map(func(n int) int    { return n * n }).        // square them
//	    Take(3).                                         // first 3
//	    All()
//	// [4, 16, 36]
//
// # Two core types
//
// [Collection[T]] — an ordered sequence of T, analogous to
//
//	collect([1, 2, 3])
//
// [MapCollection[K, V]] — a key/value store of K → V, analogous to
//
//	collect(['name' => 'Alice', 'age' => 30])
//
// Both types expose a rich method set and can interoperate via package-level
// functions such as [KeyBy], [Pluck], [Combine], [Zip], and [Flip].
//
// # Method parity with Laravel
//
// The following methods are implemented:
//
//	All           Average/Avg    Chunk         Collapse/Flatten
//	Combine       Concat         Contains      Count
//	Diff          DiffKeys       Duplicate     Each/EachWithIndex
//	Every         Except         Filter        First/FirstOrFail
//	FlatMap       Flip           Forget        ForPage
//	Get           GroupBy        Has           Implode/ImplodeWith
//	Intersect     IsEmpty        IsNotEmpty    KeyBy
//	Keys          Last           Map           MapKeys
//	Max           Median         Merge         Min
//	Mode          Nth            Pad           Partition
//	Pipe          Pluck          Pop           Prepend
//	Pull          Push           Put           Random
//	Reduce        Reject         Reverse       Search
//	Shift         Shuffle        Skip          SkipUntil/SkipWhile
//	Slice         Sort/SortDesc  Splice        Split
//	Sum           Take           TakeUntil     TakeWhile
//	Tap           Times          ToJSON        ToSlice
//	Transform     Union          Unique        UniqueBy
//	Unless*       Values         When*         Zip
//
// (* includes WhenEmpty, WhenNotEmpty, UnlessEmpty, UnlessNotEmpty variants)
//
// # Design differences from PHP/Laravel
//
// Go's static type system requires a few deliberate adaptations:
//
//  1. Callbacks accept (T) instead of ($value, $key). For index access use the
//     *WithIndex variants (e.g. [Collection.MapWithIndex],
//     [Collection.EachWithIndex]).
//
//  2. Methods that return a single non-collection value return (T, bool) or
//     (T, error) instead of panicking or returning nil.  For example:
//
//	v, ok  := c.First(nil)
//	v, err := c.FirstOrFail(nil)
//
//  3. PHP's key/value map behaviour is provided by [MapCollection[K, V]].
//     Use [KeyBy] to convert a Collection into a MapCollection.
//
//  4. Methods that transform in place (Push, Prepend, Put, Forget, Pop, Shift,
//     Transform, Splice) are explicitly marked as mutable in their GoDoc.
//     All other methods return a new Collection.
//
//  5. Equality for Unique, Diff, and Intersect is determined via
//     fmt.Sprintf("%v") by default.  For custom equality use the *By variants
//     (e.g. [Collection.UniqueBy]).
//
// # Thread safety
//
// Collections are NOT safe for concurrent mutation.  If you need to share a
// collection across goroutines, either use external synchronisation or work
// with immutable transformations (which return new collections and do not
// mutate the receiver).
//
// # Performance
//
// Immutable methods pre-allocate result slices where the output size is known
// (Map, Chunk, Pad, …).  Variable-size results (Filter, Reject, …) use append.
// See the benchmark file bench_test.go for throughput numbers; typical
// single-core throughput on an M-series chip is:
//
//	Map 10 k ints   ~  60 µs  (~170 MB/s)
//	Filter 10 k     ~  20 µs
//	Sort 10 k       ~ 900 µs
//	GroupBy 10 k    ~ 700 µs
//	Pipeline (10 k) ~   1 ms  (filter→map→sort→take)
//
// [Laravel's Illuminate\Support\Collection]: https://laravel.com/docs/collections
package collection
