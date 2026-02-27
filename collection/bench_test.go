package collection_test

// Benchmarks for Collection[T] and MapCollection[K, V].
//
// Run with:
//
//	go test -bench=. -benchmem ./...
//	go test -bench=BenchmarkFilter -benchmem -count=5  # single benchmark, 5 runs
//	go test -bench=. -benchtime=3s                      # run each for 3 seconds

import (
	"strconv"
	"testing"

	"github.com/km-arc/go-collections"
)

// ─── seed data helpers ────────────────────────────────────────────────────────

func makeInts(n int) *collection.Collection[int] {
	s := make([]int, n)
	for i := range s {
		s[i] = i + 1
	}
	return collection.New(s)
}

func makeStrings(n int) *collection.Collection[string] {
	s := make([]string, n)
	for i := range s {
		s[i] = "item-" + strconv.Itoa(i)
	}
	return collection.New(s)
}

// ─── Construction ─────────────────────────────────────────────────────────────

func BenchmarkNew_100(b *testing.B) {
	s := make([]int, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection.New(s)
	}
}

func BenchmarkNew_10k(b *testing.B) {
	s := make([]int, 10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection.New(s)
	}
}

func BenchmarkTimes_1k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		collection.Times(1_000, func(i int) int { return i * i })
	}
}

// ─── Map ──────────────────────────────────────────────────────────────────────

func BenchmarkMap_100(b *testing.B) {
	c := makeInts(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Map(func(v int) int { return v * 2 })
	}
}

func BenchmarkMap_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Map(func(v int) int { return v * 2 })
	}
}

func BenchmarkMap_100k(b *testing.B) {
	c := makeInts(100_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Map(func(v int) int { return v * 2 })
	}
}

// ─── Filter ───────────────────────────────────────────────────────────────────

func BenchmarkFilter_100(b *testing.B) {
	c := makeInts(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Filter(func(v int) bool { return v%2 == 0 })
	}
}

func BenchmarkFilter_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Filter(func(v int) bool { return v%2 == 0 })
	}
}

// ─── Reduce / Sum / Avg ───────────────────────────────────────────────────────

func BenchmarkReduce_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Reduce(0, func(carry, v int) int { return carry + v })
	}
}

func BenchmarkSum_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Sum(func(v int) float64 { return float64(v) })
	}
}

func BenchmarkAvg_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Avg(func(v int) float64 { return float64(v) })
	}
}

// ─── Sort ─────────────────────────────────────────────────────────────────────

func BenchmarkSort_100(b *testing.B) {
	// Re-create each time so we sort unsorted data each iteration.
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		c := makeInts(100)
		b.StartTimer()
		c.Sort(func(a, z int) bool { return a < z })
	}
}

func BenchmarkSort_10k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		c := makeInts(10_000)
		b.StartTimer()
		c.Sort(func(a, z int) bool { return a < z })
	}
}

// ─── GroupBy ──────────────────────────────────────────────────────────────────

func BenchmarkGroupBy_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GroupBy(func(v int) string {
			if v%2 == 0 {
				return "even"
			}
			return "odd"
		})
	}
}

// ─── Unique ───────────────────────────────────────────────────────────────────

func BenchmarkUnique_1k_few_dupes(b *testing.B) {
	c := makeInts(1_000) // no duplicates
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Unique()
	}
}

func BenchmarkUnique_1k_many_dupes(b *testing.B) {
	// 1000 items but only 10 distinct values → heavy dedup work
	s := make([]int, 1_000)
	for i := range s {
		s[i] = i % 10
	}
	c := collection.New(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Unique()
	}
}

// ─── Chunk ────────────────────────────────────────────────────────────────────

func BenchmarkChunk_10k_size50(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Chunk(50)
	}
}

// ─── FlatMap ──────────────────────────────────────────────────────────────────

func BenchmarkFlatMap_1k(b *testing.B) {
	c := makeInts(1_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.FlatMap(func(v int) []int { return []int{v, v * v} })
	}
}

// ─── Chained pipeline ─────────────────────────────────────────────────────────

// BenchmarkPipeline measures a realistic filter→map→sort→take chain.
func BenchmarkPipeline_10k(b *testing.B) {
	c := makeInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Filter(func(v int) bool { return v%3 == 0 }).
			Map(func(v int) int { return v * 2 }).
			Sort(func(a, z int) bool { return a < z }).
			Take(100)
	}
}

// ─── String collection ────────────────────────────────────────────────────────

func BenchmarkImplode_1k(b *testing.B) {
	c := makeStrings(1_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Implode(", ")
	}
}

func BenchmarkToJSON_1k(b *testing.B) {
	c := makeInts(1_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.ToJSON()
	}
}

// ─── MapCollection ────────────────────────────────────────────────────────────

func BenchmarkMapCollectionFilter_10k(b *testing.B) {
	m := make(map[int]int, 10_000)
	for i := 0; i < 10_000; i++ {
		m[i] = i
	}
	mc := collection.NewMap(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Filter(func(_ int, v int) bool { return v%2 == 0 })
	}
}

func BenchmarkKeyBy_10k(b *testing.B) {
	type Item struct{ ID int }
	s := make([]Item, 10_000)
	for i := range s {
		s[i] = Item{i}
	}
	c := collection.New(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection.KeyBy(c, func(it Item) int { return it.ID })
	}
}

func BenchmarkPluck_10k(b *testing.B) {
	type Item struct{ Name string }
	s := make([]Item, 10_000)
	for i := range s {
		s[i] = Item{"name-" + strconv.Itoa(i)}
	}
	c := collection.New(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection.Pluck(c, func(it Item) string { return it.Name })
	}
}
