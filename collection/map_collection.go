package collection

import (
	"encoding/json"
	"fmt"
)

// MapCollection is a generic key/value collection that mirrors the associative
// (map-keyed) operations in Laravel's Collection API: KeyBy, GroupBy results,
// Combine, Pluck, Flip, Only, Except, etc.
//
// The key type K must be comparable; the value type V can be anything.
//
//	mc := NewMap(map[string]int{"a": 1, "b": 2})
//	mc.Keys().All()   // ["a", "b"]  (order not guaranteed)
//	mc.Values().All() // [1, 2]
type MapCollection[K comparable, V any] struct {
	items map[K]V
}

// NewMap creates a MapCollection from an existing map.
// A nil map produces an empty MapCollection.
func NewMap[K comparable, V any](m map[K]V) *MapCollection[K, V] {
	cp := make(map[K]V, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return &MapCollection[K, V]{items: cp}
}

// All returns a copy of the underlying map.
func (mc *MapCollection[K, V]) All() map[K]V {
	cp := make(map[K]V, len(mc.items))
	for k, v := range mc.items {
		cp[k] = v
	}
	return cp
}

// Count returns the number of key/value pairs.
func (mc *MapCollection[K, V]) Count() int {
	return len(mc.items)
}

// IsEmpty reports whether the MapCollection has no entries.
func (mc *MapCollection[K, V]) IsEmpty() bool {
	return len(mc.items) == 0
}

// IsNotEmpty reports whether the MapCollection has at least one entry.
func (mc *MapCollection[K, V]) IsNotEmpty() bool {
	return len(mc.items) > 0
}

// Has reports whether key k exists.
func (mc *MapCollection[K, V]) Has(k K) bool {
	_, ok := mc.items[k]
	return ok
}

// Get retrieves the value for key k. The second return value is false if the
// key does not exist.
func (mc *MapCollection[K, V]) Get(k K) (V, bool) {
	v, ok := mc.items[k]
	return v, ok
}

// Put sets the value for key k in place.
func (mc *MapCollection[K, V]) Put(k K, v V) *MapCollection[K, V] {
	mc.items[k] = v
	return mc
}

// Forget deletes key k in place.
func (mc *MapCollection[K, V]) Forget(k K) *MapCollection[K, V] {
	delete(mc.items, k)
	return mc
}

// Keys returns a Collection containing all keys.
func (mc *MapCollection[K, V]) Keys() *Collection[K] {
	result := make([]K, 0, len(mc.items))
	for k := range mc.items {
		result = append(result, k)
	}
	return New(result)
}

// Values returns a Collection containing all values.
func (mc *MapCollection[K, V]) Values() *Collection[V] {
	result := make([]V, 0, len(mc.items))
	for _, v := range mc.items {
		result = append(result, v)
	}
	return New(result)
}

// Map transforms each value using fn and returns a new MapCollection.
func (mc *MapCollection[K, V]) Map(fn func(K, V) V) *MapCollection[K, V] {
	result := make(map[K]V, len(mc.items))
	for k, v := range mc.items {
		result[k] = fn(k, v)
	}
	return NewMap(result)
}

// Filter keeps only entries for which fn returns true.
func (mc *MapCollection[K, V]) Filter(fn func(K, V) bool) *MapCollection[K, V] {
	result := make(map[K]V)
	for k, v := range mc.items {
		if fn(k, v) {
			result[k] = v
		}
	}
	return NewMap(result)
}

// Reject is the inverse of Filter.
func (mc *MapCollection[K, V]) Reject(fn func(K, V) bool) *MapCollection[K, V] {
	return mc.Filter(func(k K, v V) bool { return !fn(k, v) })
}

// Each calls fn for every key/value pair. Returning false stops iteration.
func (mc *MapCollection[K, V]) Each(fn func(K, V) bool) *MapCollection[K, V] {
	for k, v := range mc.items {
		if !fn(k, v) {
			break
		}
	}
	return mc
}

// Every reports whether all entries satisfy fn.
func (mc *MapCollection[K, V]) Every(fn func(K, V) bool) bool {
	for k, v := range mc.items {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// Contains reports whether any entry satisfies fn.
func (mc *MapCollection[K, V]) Contains(fn func(K, V) bool) bool {
	for k, v := range mc.items {
		if fn(k, v) {
			return true
		}
	}
	return false
}

// Only returns a new MapCollection with only the specified keys.
// Matches Laravel's ->only($keys).
func (mc *MapCollection[K, V]) Only(keys ...K) *MapCollection[K, V] {
	result := make(map[K]V)
	for _, k := range keys {
		if v, ok := mc.items[k]; ok {
			result[k] = v
		}
	}
	return NewMap(result)
}

// Except returns a new MapCollection excluding the specified keys.
func (mc *MapCollection[K, V]) Except(keys ...K) *MapCollection[K, V] {
	exclude := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		exclude[k] = struct{}{}
	}
	result := make(map[K]V)
	for k, v := range mc.items {
		if _, skip := exclude[k]; !skip {
			result[k] = v
		}
	}
	return NewMap(result)
}

// Merge returns a new MapCollection with entries from both. Entries in other
// overwrite those in the receiver for matching keys.
func (mc *MapCollection[K, V]) Merge(other *MapCollection[K, V]) *MapCollection[K, V] {
	result := mc.All()
	for k, v := range other.items {
		result[k] = v
	}
	return NewMap(result)
}

// Diff returns a new MapCollection containing key/value pairs whose values are
// not present in other (value comparison via fmt.Sprintf).
func (mc *MapCollection[K, V]) Diff(other *MapCollection[K, V]) *MapCollection[K, V] {
	exclude := make(map[string]struct{}, len(other.items))
	for _, v := range other.items {
		exclude[fmt.Sprintf("%v", v)] = struct{}{}
	}
	result := make(map[K]V)
	for k, v := range mc.items {
		if _, ok := exclude[fmt.Sprintf("%v", v)]; !ok {
			result[k] = v
		}
	}
	return NewMap(result)
}

// DiffKeys returns a new MapCollection whose keys are not present in other.
func (mc *MapCollection[K, V]) DiffKeys(other *MapCollection[K, V]) *MapCollection[K, V] {
	result := make(map[K]V)
	for k, v := range mc.items {
		if !other.Has(k) {
			result[k] = v
		}
	}
	return NewMap(result)
}

// Intersect returns only key/value pairs whose values appear in other.
func (mc *MapCollection[K, V]) Intersect(other *MapCollection[K, V]) *MapCollection[K, V] {
	include := make(map[string]struct{}, len(other.items))
	for _, v := range other.items {
		include[fmt.Sprintf("%v", v)] = struct{}{}
	}
	result := make(map[K]V)
	for k, v := range mc.items {
		if _, ok := include[fmt.Sprintf("%v", v)]; ok {
			result[k] = v
		}
	}
	return NewMap(result)
}

// IntersectByKeys returns only entries whose keys appear in other.
func (mc *MapCollection[K, V]) IntersectByKeys(other *MapCollection[K, V]) *MapCollection[K, V] {
	result := make(map[K]V)
	for k, v := range mc.items {
		if other.Has(k) {
			result[k] = v
		}
	}
	return NewMap(result)
}

// Flip returns a new MapCollection with keys and values swapped (only works when
// V is comparable). Provided as a package-level function due to Go constraints.
//
//	Flip(NewMap(map[string]string{"a":"x","b":"y"})) // map[string]string{"x":"a","y":"b"}
func Flip[K comparable, V comparable](mc *MapCollection[K, V]) *MapCollection[V, K] {
	result := make(map[V]K, len(mc.items))
	for k, v := range mc.items {
		result[v] = k
	}
	return NewMap(result)
}

// Combine creates a MapCollection from a keys collection and a values collection.
// Panics if the two collections have different lengths.
//
//	Combine(New([]string{"a","b"}), New([]int{1,2}))
//	// MapCollection{"a":1, "b":2}
func Combine[K comparable, V any](keys *Collection[K], values *Collection[V]) *MapCollection[K, V] {
	if keys.Count() != values.Count() {
		panic("collection: Combine requires equal-length collections")
	}
	result := make(map[K]V, keys.Count())
	for i, k := range keys.items {
		result[k] = values.items[i]
	}
	return NewMap(result)
}

// KeyBy creates a MapCollection by keying a Collection using fn.
//
//	type User struct{ ID int; Name string }
//	users := New([]User{{1,"Alice"},{2,"Bob"}})
//	m := KeyBy(users, func(u User) int { return u.ID })
//	m.Get(1) // User{1,"Alice"}, true
func KeyBy[K comparable, V any](c *Collection[V], fn func(V) K) *MapCollection[K, V] {
	result := make(map[K]V, c.Count())
	for _, v := range c.items {
		result[fn(v)] = v
	}
	return NewMap(result)
}

// Pluck extracts a field from each item and returns a Collection of those values.
//
//	type P struct{ Name string }
//	Pluck(New([]P{{"Alice"},{"Bob"}}), func(p P) string { return p.Name }).All()
//	// ["Alice","Bob"]
func Pluck[T any, V any](c *Collection[T], fn func(T) V) *Collection[V] {
	result := make([]V, len(c.items))
	for i, v := range c.items {
		result[i] = fn(v)
	}
	return New(result)
}

// MapKeys returns a new MapCollection whose keys have been transformed by fn.
func MapKeys[K1 comparable, K2 comparable, V any](mc *MapCollection[K1, V], fn func(K1, V) K2) *MapCollection[K2, V] {
	result := make(map[K2]V, len(mc.items))
	for k, v := range mc.items {
		result[fn(k, v)] = v
	}
	return NewMap(result)
}

// ToJSON serialises the MapCollection to a JSON object.
func (mc *MapCollection[K, V]) ToJSON() (string, error) {
	b, err := json.Marshal(mc.items)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Union returns a new MapCollection with all entries from other whose keys are
// NOT already present in the receiver.
func (mc *MapCollection[K, V]) Union(other *MapCollection[K, V]) *MapCollection[K, V] {
	result := mc.All()
	for k, v := range other.items {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	return NewMap(result)
}
