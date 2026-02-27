// Package collection provides a generic, fluent collection type for Go that mirrors
// the Laravel Illuminate\Support\Collection API as closely as Go's type system allows.
//
// Collections wrap a slice of items and expose a rich set of functional methods
// (Map, Filter, Reduce, GroupBy, …) that can be chained together:
//
//	result := New([]int{1, 2, 3, 4, 5}).
//	    Filter(func(v int) bool { return v%2 == 0 }).
//	    Map(func(v int) int { return v * 10 }).
//	    All()
//	// [20, 40]
//
// # Design choices vs Laravel
//
//   - Go generics replace PHP's dynamic typing. The type parameter T is the
//     element type; for mixed-type collections use Collection[any].
//   - Methods that return a single non-collection value (First, Last, Sum, …)
//     return (T, bool) or (T, error) instead of panicking.
//   - PHP's key/value maps are modelled with the separate MapCollection[K, V] type
//     (see map_collection.go) for methods like KeyBy, GroupBy, Combine, etc.
//   - Methods that do not have a meaningful Go equivalent (e.g. macro, make as a
//     static constructor) are either omitted or provided as package-level helpers.
package collection

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
)

// Collection is a generic, ordered sequence of items of type T.
// It is immutable by convention — every method that transforms the collection
// returns a new Collection rather than modifying the receiver in place.
// The only exceptions are the mutable helpers Push, Prepend, Put, Forget, Pop,
// Shift, and Transform, which mirror Laravel's mutable counterparts.
type Collection[T any] struct {
	items []T
}

// ─────────────────────────────────────────────────────────────────────────────
// Constructors
// ─────────────────────────────────────────────────────────────────────────────

// New creates a new Collection from a slice.
// A nil or empty slice produces an empty collection.
//
//	c := New([]string{"alice", "bob", "carol"})
func New[T any](items []T) *Collection[T] {
	cp := make([]T, len(items))
	copy(cp, items)
	return &Collection[T]{items: cp}
}

// Collect is a convenience alias for New, matching Laravel's collect() helper.
//
//	c := Collect([]int{1, 2, 3})
func Collect[T any](items []T) *Collection[T] {
	return New(items)
}

// Times creates a new collection containing n items produced by calling fn with
// the 1-based index (matching Laravel's 1-based behaviour).
//
//	c := Times(5, func(i int) int { return i * i }) // [1, 4, 9, 16, 25]
func Times[T any](n int, fn func(index int) T) *Collection[T] {
	items := make([]T, n)
	for i := 0; i < n; i++ {
		items[i] = fn(i + 1)
	}
	return New(items)
}

// ─────────────────────────────────────────────────────────────────────────────
// Basic accessors
// ─────────────────────────────────────────────────────────────────────────────

// All returns a copy of the underlying slice.
//
//	New([]int{1, 2, 3}).All() // []int{1, 2, 3}
func (c *Collection[T]) All() []T {
	cp := make([]T, len(c.items))
	copy(cp, c.items)
	return cp
}

// Count returns the number of items in the collection.
// Mirrors Laravel's count() / ->count().
func (c *Collection[T]) Count() int {
	return len(c.items)
}

// IsEmpty reports whether the collection contains no items.
func (c *Collection[T]) IsEmpty() bool {
	return len(c.items) == 0
}

// IsNotEmpty reports whether the collection contains at least one item.
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.items) > 0
}

// First returns the first item that passes the optional predicate.
// If no predicate is given (pass nil) the very first item is returned.
// The second return value is false when the collection is empty or no item
// satisfies the predicate.
//
//	c := New([]int{1, 2, 3, 4})
//	v, ok := c.First(nil)           // 1, true
//	v, ok  = c.First(func(v int) bool { return v > 2 }) // 3, true
func (c *Collection[T]) First(fn func(T) bool) (T, bool) {
	for _, v := range c.items {
		if fn == nil || fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FirstOrFail is like First but returns an error instead of false when nothing
// is found, matching Laravel's ->firstOrFail().
func (c *Collection[T]) FirstOrFail(fn func(T) bool) (T, error) {
	v, ok := c.First(fn)
	if !ok {
		return v, fmt.Errorf("collection: item not found")
	}
	return v, nil
}

// Last returns the last item that passes the optional predicate.
// If no predicate is given (pass nil) the very last item is returned.
//
//	v, ok := New([]int{1, 2, 3}).Last(nil) // 3, true
func (c *Collection[T]) Last(fn func(T) bool) (T, bool) {
	for i := len(c.items) - 1; i >= 0; i-- {
		if fn == nil || fn(c.items[i]) {
			return c.items[i], true
		}
	}
	var zero T
	return zero, false
}

// Nth returns every n-th element, starting from offset (0-based).
// Matches Laravel's ->nth($n, $offset).
//
//	New([]int{1,2,3,4,5,6}).Nth(2, 0).All() // [1, 3, 5]
func (c *Collection[T]) Nth(n, offset int) *Collection[T] {
	var result []T
	for i := offset; i < len(c.items); i += n {
		result = append(result, c.items[i])
	}
	return New(result)
}

// Get returns the item at position (0-based). The second return value is false
// when the index is out of range.
func (c *Collection[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(c.items) {
		var zero T
		return zero, false
	}
	return c.items[index], true
}

// ─────────────────────────────────────────────────────────────────────────────
// Transformation – returns new Collection[T]
// ─────────────────────────────────────────────────────────────────────────────

// Map passes each item through fn and returns a new collection of the results.
// The callback receives the item; use MapWithIndex for (value, index) callbacks.
//
//	New([]int{1,2,3}).Map(func(v int) int { return v*2 }).All() // [2,4,6]
func (c *Collection[T]) Map(fn func(T) T) *Collection[T] {
	result := make([]T, len(c.items))
	for i, v := range c.items {
		result[i] = fn(v)
	}
	return New(result)
}

// MapWithIndex is like Map but the callback also receives the 0-based index.
func (c *Collection[T]) MapWithIndex(fn func(T, int) T) *Collection[T] {
	result := make([]T, len(c.items))
	for i, v := range c.items {
		result[i] = fn(v, i)
	}
	return New(result)
}

// Filter keeps only items for which fn returns true.
// With no predicate (nil) it removes zero-value items (shallow falsy check via
// JSON round-trip is avoided; callers should pass an explicit predicate).
//
//	New([]int{1,2,3,4}).Filter(func(v int) bool { return v%2==0 }).All() // [2,4]
func (c *Collection[T]) Filter(fn func(T) bool) *Collection[T] {
	var result []T
	for _, v := range c.items {
		if fn(v) {
			result = append(result, v)
		}
	}
	return New(result)
}

// Reject is the inverse of Filter — keeps only items for which fn returns false.
//
//	New([]int{1,2,3,4}).Reject(func(v int) bool { return v%2==0 }).All() // [1,3]
func (c *Collection[T]) Reject(fn func(T) bool) *Collection[T] {
	return c.Filter(func(v T) bool { return !fn(v) })
}

// Each calls fn for every item. Returning false from fn stops iteration early
// (like Laravel's callback returning false).
//
//	New([]int{1,2,3}).Each(func(v int) bool {
//	    fmt.Println(v)
//	    return true
//	})
func (c *Collection[T]) Each(fn func(T) bool) *Collection[T] {
	for _, v := range c.items {
		if !fn(v) {
			break
		}
	}
	return c
}

// EachWithIndex is like Each but the callback also receives the 0-based index.
func (c *Collection[T]) EachWithIndex(fn func(T, int) bool) *Collection[T] {
	for i, v := range c.items {
		if !fn(v, i) {
			break
		}
	}
	return c
}

// Tap calls fn with the collection and returns the unchanged collection,
// allowing side-effects mid-chain. Mirrors Laravel's ->tap().
func (c *Collection[T]) Tap(fn func(*Collection[T])) *Collection[T] {
	fn(c)
	return c
}

// Pipe passes the collection to fn and returns the result.
// Useful for inserting arbitrary transformations mid-chain.
func (c *Collection[T]) Pipe(fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return fn(c)
}

// Slice returns a sub-collection starting at offset, limited to length items.
// Passing a negative length returns all remaining items.
// Mirrors Laravel's ->slice($offset, $length).
func (c *Collection[T]) Slice(offset, length int) *Collection[T] {
	if offset < 0 {
		offset = max(0, len(c.items)+offset)
	}
	if offset >= len(c.items) {
		return New([]T{})
	}
	end := len(c.items)
	if length >= 0 {
		end = min(offset+length, end)
	}
	return New(c.items[offset:end])
}

// Take returns the first n items. Negative n takes from the end.
//
//	New([]int{1,2,3,4,5}).Take(3).All()  // [1,2,3]
//	New([]int{1,2,3,4,5}).Take(-2).All() // [4,5]
func (c *Collection[T]) Take(n int) *Collection[T] {
	if n >= 0 {
		return c.Slice(0, n)
	}
	return c.Slice(n, -1)
}

// TakeUntil returns items until fn returns true.
func (c *Collection[T]) TakeUntil(fn func(T) bool) *Collection[T] {
	var result []T
	for _, v := range c.items {
		if fn(v) {
			break
		}
		result = append(result, v)
	}
	return New(result)
}

// TakeWhile returns items while fn returns true.
func (c *Collection[T]) TakeWhile(fn func(T) bool) *Collection[T] {
	var result []T
	for _, v := range c.items {
		if !fn(v) {
			break
		}
		result = append(result, v)
	}
	return New(result)
}

// Skip skips the first n items and returns the rest.
func (c *Collection[T]) Skip(n int) *Collection[T] {
	return c.Slice(n, -1)
}

// SkipUntil skips items until fn returns true, then returns the remaining items.
func (c *Collection[T]) SkipUntil(fn func(T) bool) *Collection[T] {
	for i, v := range c.items {
		if fn(v) {
			return New(c.items[i:])
		}
	}
	return New([]T{})
}

// SkipWhile skips items while fn returns true.
func (c *Collection[T]) SkipWhile(fn func(T) bool) *Collection[T] {
	for i, v := range c.items {
		if !fn(v) {
			return New(c.items[i:])
		}
	}
	return New([]T{})
}

// Reverse returns a new collection with items in reverse order.
func (c *Collection[T]) Reverse() *Collection[T] {
	n := len(c.items)
	result := make([]T, n)
	for i, v := range c.items {
		result[n-1-i] = v
	}
	return New(result)
}

// Shuffle returns a new collection with items in random order.
func (c *Collection[T]) Shuffle() *Collection[T] {
	result := c.All()
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return New(result)
}

// Sort returns a sorted collection using the provided less function.
// The original collection is not modified.
//
//	New([]int{3,1,2}).Sort(func(a, b int) bool { return a < b }).All() // [1,2,3]
func (c *Collection[T]) Sort(less func(a, b T) bool) *Collection[T] {
	result := c.All()
	sort.SliceStable(result, func(i, j int) bool { return less(result[i], result[j]) })
	return New(result)
}

// SortDesc returns a reverse-sorted collection using the provided less function.
func (c *Collection[T]) SortDesc(less func(a, b T) bool) *Collection[T] {
	return c.Sort(func(a, b T) bool { return less(b, a) })
}

// Chunk splits the collection into multiple smaller collections of the given size.
//
//	New([]int{1,2,3,4,5}).Chunk(2) // [[1,2],[3,4],[5]]
func (c *Collection[T]) Chunk(size int) []*Collection[T] {
	if size <= 0 {
		return nil
	}
	var chunks []*Collection[T]
	for i := 0; i < len(c.items); i += size {
		end := min(i+size, len(c.items))
		chunks = append(chunks, New(c.items[i:end]))
	}
	return chunks
}

// Split divides the collection into n groups of roughly equal size.
// Mirrors Laravel's ->split($n).
func (c *Collection[T]) Split(n int) []*Collection[T] {
	if n <= 0 || len(c.items) == 0 {
		return nil
	}
	size := (len(c.items) + n - 1) / n
	return c.Chunk(size)
}

// Flatten collapses one level of nesting. Because Go is statically typed, this
// only applies to Collection[[]T] — use FlatMap for the general case.
// This package-level function is provided instead of a method to work around
// Go's generic constraints.
//
//	Flatten(New([][]int{{1,2},{3,4}})).All() // [1,2,3,4]
func Flatten[T any](c *Collection[[]T]) *Collection[T] {
	var result []T
	for _, sub := range c.items {
		result = append(result, sub...)
	}
	return New(result)
}

// FlatMap maps over the collection then flattens one level.
//
//	New([]int{1,2,3}).FlatMap(func(v int) []int { return []int{v, v*v} }).All()
//	// [1, 1, 2, 4, 3, 9]
func (c *Collection[T]) FlatMap(fn func(T) []T) *Collection[T] {
	var result []T
	for _, v := range c.items {
		result = append(result, fn(v)...)
	}
	return New(result)
}

// Collapse collapses a Collection of slices into a flat Collection.
// This is a package-level function — same reason as Flatten.
func Collapse[T any](c *Collection[[]T]) *Collection[T] {
	return Flatten(c)
}

// Unique returns a new collection with duplicate items removed.
// It uses a fmt.Sprintf key to determine uniqueness — callers that need
// custom equality should use UniqueBy.
func (c *Collection[T]) Unique() *Collection[T] {
	seen := make(map[string]struct{})
	var result []T
	for _, v := range c.items {
		key := fmt.Sprintf("%v", v)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			result = append(result, v)
		}
	}
	return New(result)
}

// UniqueBy returns a new collection with duplicates removed, using fn to
// produce a string key for comparison.
//
//	type User struct{ Name string; Age int }
//	users := New([]User{{"Alice",30},{"Bob",25},{"Alice",31}})
//	users.UniqueBy(func(u User) string { return u.Name }).Count() // 2
func (c *Collection[T]) UniqueBy(fn func(T) string) *Collection[T] {
	seen := make(map[string]struct{})
	var result []T
	for _, v := range c.items {
		k := fn(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, v)
		}
	}
	return New(result)
}

// Duplicate returns items that appear more than once in the collection.
func (c *Collection[T]) Duplicate() *Collection[T] {
	counts := make(map[string]int)
	for _, v := range c.items {
		counts[fmt.Sprintf("%v", v)]++
	}
	var result []T
	seen := make(map[string]bool)
	for _, v := range c.items {
		k := fmt.Sprintf("%v", v)
		if counts[k] > 1 && !seen[k] {
			seen[k] = true
			result = append(result, v)
		}
	}
	return New(result)
}

// Merge appends the items from other to a new collection.
// Mirrors Laravel's ->merge() on list collections.
func (c *Collection[T]) Merge(other *Collection[T]) *Collection[T] {
	result := c.All()
	result = append(result, other.items...)
	return New(result)
}

// Concat appends all items from others to a new collection.
func (c *Collection[T]) Concat(others ...*Collection[T]) *Collection[T] {
	result := c.All()
	for _, o := range others {
		result = append(result, o.items...)
	}
	return New(result)
}

// Diff returns a new collection containing items from c that are not in other.
// Comparison uses fmt.Sprintf("%v") — use DiffBy for custom equality.
func (c *Collection[T]) Diff(other *Collection[T]) *Collection[T] {
	exclude := make(map[string]struct{}, len(other.items))
	for _, v := range other.items {
		exclude[fmt.Sprintf("%v", v)] = struct{}{}
	}
	var result []T
	for _, v := range c.items {
		if _, ok := exclude[fmt.Sprintf("%v", v)]; !ok {
			result = append(result, v)
		}
	}
	return New(result)
}

// Intersect returns only items that are present in both collections.
func (c *Collection[T]) Intersect(other *Collection[T]) *Collection[T] {
	include := make(map[string]struct{}, len(other.items))
	for _, v := range other.items {
		include[fmt.Sprintf("%v", v)] = struct{}{}
	}
	var result []T
	for _, v := range c.items {
		if _, ok := include[fmt.Sprintf("%v", v)]; ok {
			result = append(result, v)
		}
	}
	return New(result)
}

// Zip merges two collections at corresponding indices into pairs.
// If the collections differ in length, the shorter one is padded with zero values.
//
//	Zip(New([]int{1,2,3}), New([]string{"a","b","c"}))
//	// [][2]any{{1,"a"},{2,"b"},{3,"c"}}
func Zip[A, B any](a *Collection[A], b *Collection[B]) [][2]any {
	n := max(len(a.items), len(b.items))
	result := make([][2]any, n)
	for i := 0; i < n; i++ {
		var va A
		var vb B
		if i < len(a.items) {
			va = a.items[i]
		}
		if i < len(b.items) {
			vb = b.items[i]
		}
		result[i] = [2]any{va, vb}
	}
	return result
}

// Pad grows the collection to size, filling new positions with value.
// Positive size pads on the right; negative pads on the left.
//
//	New([]int{1,2,3}).Pad(5, 0).All()  // [1,2,3,0,0]
//	New([]int{1,2,3}).Pad(-5, 0).All() // [0,0,1,2,3]
func (c *Collection[T]) Pad(size int, value T) *Collection[T] {
	abs := size
	if abs < 0 {
		abs = -abs
	}
	if abs <= len(c.items) {
		return New(c.items)
	}
	padding := make([]T, abs-len(c.items))
	for i := range padding {
		padding[i] = value
	}
	if size > 0 {
		return New(append(c.All(), padding...))
	}
	return New(append(padding, c.items...))
}

// ForPage returns the items for a given page number with the given per-page size.
// Pages are 1-based. Mirrors Laravel's ->forPage($page, $perPage).
func (c *Collection[T]) ForPage(page, perPage int) *Collection[T] {
	return c.Slice((page-1)*perPage, perPage)
}

// ─────────────────────────────────────────────────────────────────────────────
// Searching / testing
// ─────────────────────────────────────────────────────────────────────────────

// Contains reports whether any item in the collection satisfies fn.
//
//	New([]int{1,2,3}).Contains(func(v int) bool { return v == 2 }) // true
func (c *Collection[T]) Contains(fn func(T) bool) bool {
	for _, v := range c.items {
		if fn(v) {
			return true
		}
	}
	return false
}

// DoesntContain is the inverse of Contains.
func (c *Collection[T]) DoesntContain(fn func(T) bool) bool {
	return !c.Contains(fn)
}

// Every reports whether all items satisfy fn.
func (c *Collection[T]) Every(fn func(T) bool) bool {
	for _, v := range c.items {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Search returns the 0-based index of the first item satisfying fn, or -1.
//
//	New([]string{"a","b","c"}).Search(func(v string) bool { return v == "b" }) // 1
func (c *Collection[T]) Search(fn func(T) bool) int {
	for i, v := range c.items {
		if fn(v) {
			return i
		}
	}
	return -1
}

// Has reports whether an item exists at the given 0-based index.
func (c *Collection[T]) Has(index int) bool {
	return index >= 0 && index < len(c.items)
}

// ─────────────────────────────────────────────────────────────────────────────
// Mutable operations (match Laravel's mutable methods)
// ─────────────────────────────────────────────────────────────────────────────

// Push appends item to the collection in place and returns the receiver.
// Matches Laravel's mutable ->push().
func (c *Collection[T]) Push(item T) *Collection[T] {
	c.items = append(c.items, item)
	return c
}

// Prepend inserts item at the beginning of the collection in place.
func (c *Collection[T]) Prepend(item T) *Collection[T] {
	c.items = append([]T{item}, c.items...)
	return c
}

// Put sets the item at the given 0-based index in place.
// If index == Count(), it appends. Returns an error if out of range.
func (c *Collection[T]) Put(index int, item T) error {
	if index == len(c.items) {
		c.items = append(c.items, item)
		return nil
	}
	if index < 0 || index >= len(c.items) {
		return fmt.Errorf("collection: index %d out of range [0, %d]", index, len(c.items))
	}
	c.items[index] = item
	return nil
}

// Forget removes the item at the given 0-based index in place.
func (c *Collection[T]) Forget(index int) error {
	if index < 0 || index >= len(c.items) {
		return fmt.Errorf("collection: index %d out of range", index)
	}
	c.items = append(c.items[:index], c.items[index+1:]...)
	return nil
}

// Pop removes and returns the last item in place.
func (c *Collection[T]) Pop() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	last := c.items[len(c.items)-1]
	c.items = c.items[:len(c.items)-1]
	return last, true
}

// Shift removes and returns the first item in place.
func (c *Collection[T]) Shift() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	first := c.items[0]
	c.items = c.items[1:]
	return first, true
}

// Pull removes and returns the item at the given index in place.
func (c *Collection[T]) Pull(index int) (T, bool) {
	v, ok := c.Get(index)
	if !ok {
		return v, false
	}
	_ = c.Forget(index)
	return v, true
}

// Splice removes and returns a sub-collection starting at offset.
// If length >= 0 only that many items are removed; otherwise all remaining items.
// Replacement items (if any) are inserted in place of the removed items.
func (c *Collection[T]) Splice(offset, length int, replacement ...T) *Collection[T] {
	if offset < 0 {
		offset = max(0, len(c.items)+offset)
	}
	if offset > len(c.items) {
		offset = len(c.items)
	}
	end := len(c.items)
	if length >= 0 {
		end = min(offset+length, end)
	}
	removed := New(c.items[offset:end])
	tail := make([]T, len(c.items[end:]))
	copy(tail, c.items[end:])
	c.items = append(c.items[:offset], replacement...)
	c.items = append(c.items, tail...)
	return removed
}

// Transform applies fn to every item in place (mutating).
// Matches Laravel's mutable ->transform().
func (c *Collection[T]) Transform(fn func(T) T) *Collection[T] {
	for i, v := range c.items {
		c.items[i] = fn(v)
	}
	return c
}

// ─────────────────────────────────────────────────────────────────────────────
// Aggregates
// ─────────────────────────────────────────────────────────────────────────────

// Reduce reduces the collection to a single value using fn.
// carry is the initial value.
//
//	New([]int{1,2,3,4}).Reduce(0, func(carry, v int) int { return carry + v }) // 10
func (c *Collection[T]) Reduce(carry T, fn func(carry, item T) T) T {
	for _, v := range c.items {
		carry = fn(carry, v)
	}
	return carry
}

// Sum returns the sum of all items using the provided extractor.
// Use SumFloat for float64 precision.
//
//	New([]int{1,2,3}).Sum(func(v int) float64 { return float64(v) }) // 6.0
func (c *Collection[T]) Sum(fn func(T) float64) float64 {
	var total float64
	for _, v := range c.items {
		total += fn(v)
	}
	return total
}

// Avg returns the arithmetic mean. Returns 0 for an empty collection.
func (c *Collection[T]) Avg(fn func(T) float64) float64 {
	if len(c.items) == 0 {
		return 0
	}
	return c.Sum(fn) / float64(len(c.items))
}

// Average is an alias for Avg, matching Laravel's ->average().
func (c *Collection[T]) Average(fn func(T) float64) float64 {
	return c.Avg(fn)
}

// Min returns the minimum value using the provided extractor.
func (c *Collection[T]) Min(fn func(T) float64) (float64, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	m := fn(c.items[0])
	for _, v := range c.items[1:] {
		if x := fn(v); x < m {
			m = x
		}
	}
	return m, true
}

// Max returns the maximum value using the provided extractor.
func (c *Collection[T]) Max(fn func(T) float64) (float64, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	m := fn(c.items[0])
	for _, v := range c.items[1:] {
		if x := fn(v); x > m {
			m = x
		}
	}
	return m, true
}

// Median returns the median value. For even-length collections the average of
// the two middle values is returned.
func (c *Collection[T]) Median(fn func(T) float64) (float64, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	vals := make([]float64, len(c.items))
	for i, v := range c.items {
		vals[i] = fn(v)
	}
	sort.Float64s(vals)
	n := len(vals)
	if n%2 == 1 {
		return vals[n/2], true
	}
	return (vals[n/2-1] + vals[n/2]) / 2, true
}

// Mode returns all values that appear most frequently.
func (c *Collection[T]) Mode(fn func(T) float64) []float64 {
	if len(c.items) == 0 {
		return nil
	}
	counts := make(map[float64]int)
	for _, v := range c.items {
		counts[fn(v)]++
	}
	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}
	var result []float64
	for val, count := range counts {
		if count == maxCount {
			result = append(result, val)
		}
	}
	sort.Float64s(result)
	return result
}

// ─────────────────────────────────────────────────────────────────────────────
// Grouping / partitioning
// ─────────────────────────────────────────────────────────────────────────────

// GroupBy groups the collection's items into a map of collections, keyed by the
// string returned from fn.
//
//	type P struct{ Cat string; Name string }
//	items := New([]P{{"fruit","apple"},{"veg","carrot"},{"fruit","banana"}})
//	groups := items.GroupBy(func(p P) string { return p.Cat })
//	groups["fruit"].Count() // 2
func (c *Collection[T]) GroupBy(fn func(T) string) map[string]*Collection[T] {
	result := make(map[string]*Collection[T])
	for _, v := range c.items {
		k := fn(v)
		if _, ok := result[k]; !ok {
			result[k] = New([]T{})
		}
		result[k].Push(v)
	}
	return result
}

// Partition splits the collection into two: items that pass the predicate and
// those that don't. Returns (passing, failing).
//
//	even, odd := New([]int{1,2,3,4}).Partition(func(v int) bool { return v%2==0 })
func (c *Collection[T]) Partition(fn func(T) bool) (*Collection[T], *Collection[T]) {
	var pass, fail []T
	for _, v := range c.items {
		if fn(v) {
			pass = append(pass, v)
		} else {
			fail = append(fail, v)
		}
	}
	return New(pass), New(fail)
}

// ─────────────────────────────────────────────────────────────────────────────
// Conditional execution
// ─────────────────────────────────────────────────────────────────────────────

// When calls fn with the collection when condition is true, then returns the
// result; otherwise returns the receiver unchanged.
// Mirrors Laravel's ->when($value, $callback, $default).
func (c *Collection[T]) When(condition bool, fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	if condition {
		return fn(c)
	}
	return c
}

// WhenEmpty calls fn when the collection is empty.
func (c *Collection[T]) WhenEmpty(fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsEmpty(), fn)
}

// WhenNotEmpty calls fn when the collection is not empty.
func (c *Collection[T]) WhenNotEmpty(fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsNotEmpty(), fn)
}

// Unless calls fn when condition is false.
func (c *Collection[T]) Unless(condition bool, fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(!condition, fn)
}

// UnlessEmpty calls fn when the collection is not empty.
func (c *Collection[T]) UnlessEmpty(fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.Unless(c.IsEmpty(), fn)
}

// UnlessNotEmpty calls fn when the collection is empty.
func (c *Collection[T]) UnlessNotEmpty(fn func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.Unless(c.IsNotEmpty(), fn)
}

// ─────────────────────────────────────────────────────────────────────────────
// String / implode
// ─────────────────────────────────────────────────────────────────────────────

// Implode joins the string representation of each item using glue.
// For custom formatting pass a mapper via ImplodeWith.
//
//	New([]int{1,2,3}).Implode(", ") // "1, 2, 3"
func (c *Collection[T]) Implode(glue string) string {
	result := ""
	for i, v := range c.items {
		if i > 0 {
			result += glue
		}
		result += fmt.Sprintf("%v", v)
	}
	return result
}

// ImplodeWith maps each item to a string using fn, then joins with glue.
func (c *Collection[T]) ImplodeWith(glue string, fn func(T) string) string {
	result := ""
	for i, v := range c.items {
		if i > 0 {
			result += glue
		}
		result += fn(v)
	}
	return result
}

// ─────────────────────────────────────────────────────────────────────────────
// Random
// ─────────────────────────────────────────────────────────────────────────────

// Random returns n randomly chosen items as a new collection.
// If n >= Count() the whole (shuffled) collection is returned.
func (c *Collection[T]) Random(n int) *Collection[T] {
	return c.Shuffle().Take(n)
}

// ─────────────────────────────────────────────────────────────────────────────
// Conversion
// ─────────────────────────────────────────────────────────────────────────────

// ToSlice is an alias for All().
func (c *Collection[T]) ToSlice() []T {
	return c.All()
}

// ToJSON serialises the collection to a JSON array.
// Returns an error if any item is not JSON-serialisable.
func (c *Collection[T]) ToJSON() (string, error) {
	b, err := json.Marshal(c.items)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// Clone returns a deep copy of the collection.
func (c *Collection[T]) Clone() *Collection[T] {
	return New(c.items)
}

// Values resets numeric keys (i.e. returns a clean copy) — useful after Sort
// or Diff to get a contiguous slice. In Go this is always true; this method
// exists purely for API parity with Laravel.
func (c *Collection[T]) Values() *Collection[T] {
	return c.Clone()
}

// Keys returns the 0-based indices of the collection as a Collection[int].
func (c *Collection[T]) Keys() *Collection[int] {
	result := make([]int, len(c.items))
	for i := range c.items {
		result[i] = i
	}
	return New(result)
}

// min/max helpers for Go < 1.21 compatibility within this file.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
