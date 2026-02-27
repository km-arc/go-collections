package collection_test

// Run with:
//   go test -bench=. -benchmem -run='^$' ./...
//
// Compare two revisions:
//   go test -bench=. -benchmem -run='^$' ./... > old.txt
//   # make changes
//   go test -bench=. -benchmem -run='^$' ./... > new.txt
//   benchstat old.txt new.txt

import (
	"fmt"
	"testing"

	"github.com/km-arc/go-collections"
)

// ─── Sizes under test ─────────────────────────────────────────────────────────

var benchSizes = []int{10, 100, 1_000, 10_000}

func makeInts(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i + 1
	}
	return s
}

// ─── Construction ─────────────────────────────────────────────────────────────

func BenchmarkNew(b *testing.B) {
	for _, size := range benchSizes {
		src := makeInts(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = collection.New(src)
			}
		})
	}
}

func BenchmarkTimes(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = collection.Times(size, func(n int) int { return n * n })
			}
		})
	}
}

// ─── Core transforms ──────────────────────────────────────────────────────────

func BenchmarkMap(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Map(func(v int) int { return v * 2 })
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Filter(func(v int) bool { return v%2 == 0 })
			}
		})
	}
}

func BenchmarkReject(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Reject(func(v int) bool { return v%2 == 0 })
			}
		})
	}
}

func BenchmarkFlatMap(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.FlatMap(func(v int) []int { return []int{v, v * v} })
			}
		})
	}
}

func BenchmarkEach(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sum := 0
				c.Each(func(v int) bool { sum += v; return true })
				_ = sum
			}
		})
	}
}

// ─── Aggregates ───────────────────────────────────────────────────────────────

func BenchmarkReduce(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Reduce(0, func(carry, v int) int { return carry + v })
			}
		})
	}
}

func BenchmarkSum(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Sum(func(v int) float64 { return float64(v) })
			}
		})
	}
}

func BenchmarkAvg(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Avg(func(v int) float64 { return float64(v) })
			}
		})
	}
}

func BenchmarkMin(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.Min(func(v int) float64 { return float64(v) })
			}
		})
	}
}

func BenchmarkMax(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.Max(func(v int) float64 { return float64(v) })
			}
		})
	}
}

func BenchmarkMedian(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.Median(func(v int) float64 { return float64(v) })
			}
		})
	}
}

// ─── Searching ────────────────────────────────────────────────────────────────

func BenchmarkFirst(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d/hit", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.First(func(v int) bool { return v == size/2 })
			}
		})
		b.Run(fmt.Sprintf("size=%d/miss", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.First(func(v int) bool { return v == size+9999 })
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Contains(func(v int) bool { return v == size })
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Search(func(v int) bool { return v == size/2 })
			}
		})
	}
}

// ─── Sorting ──────────────────────────────────────────────────────────────────

func BenchmarkSort(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Sort(func(a, b int) bool { return a < b })
			}
		})
	}
}

func BenchmarkShuffle(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Shuffle()
			}
		})
	}
}

// ─── Grouping ─────────────────────────────────────────────────────────────────

func BenchmarkGroupBy(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.GroupBy(func(v int) string {
					if v%2 == 0 {
						return "even"
					}
					return "odd"
				})
			}
		})
	}
}

func BenchmarkPartition(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.Partition(func(v int) bool { return v%2 == 0 })
			}
		})
	}
}

// ─── Chunk / Split ────────────────────────────────────────────────────────────

func BenchmarkChunk(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d/chunk=10", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Chunk(10)
			}
		})
	}
}

// ─── Unique ───────────────────────────────────────────────────────────────────

func BenchmarkUnique(b *testing.B) {
	for _, size := range benchSizes {
		// half duplicates
		src := makeInts(size / 2)
		src = append(src, src...)
		c := collection.New(src)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.Unique()
			}
		})
	}
}

// ─── Set operations ───────────────────────────────────────────────────────────

func BenchmarkDiff(b *testing.B) {
	for _, size := range benchSizes {
		a := collection.New(makeInts(size))
		bSlice := makeInts(size / 2)
		bc := collection.New(bSlice)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = a.Diff(bc)
			}
		})
	}
}

func BenchmarkIntersect(b *testing.B) {
	for _, size := range benchSizes {
		a := collection.New(makeInts(size))
		bSlice := makeInts(size / 2)
		bc := collection.New(bSlice)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = a.Intersect(bc)
			}
		})
	}
}

// ─── Serialisation ────────────────────────────────────────────────────────────

func BenchmarkToJSON(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = c.ToJSON()
			}
		})
	}
}

// ─── Chaining (pipeline) ──────────────────────────────────────────────────────

// BenchmarkPipeline measures the cost of a realistic multi-step chain,
// the kind you'd write in production code.
func BenchmarkPipeline(b *testing.B) {
	for _, size := range benchSizes {
		c := collection.New(makeInts(size))
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.
					Filter(func(v int) bool { return v%2 == 0 }).
					Map(func(v int) int { return v * 3 }).
					Reject(func(v int) bool { return v > size }).
					Take(10).
					All()
			}
		})
	}
}

// ─── MapCollection ────────────────────────────────────────────────────────────

func BenchmarkNewMap(b *testing.B) {
	for _, size := range benchSizes {
		src := make(map[int]int, size)
		for i := 0; i < size; i++ {
			src[i] = i * i
		}
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = collection.NewMap(src)
			}
		})
	}
}

func BenchmarkMapCollectionFilter(b *testing.B) {
	for _, size := range benchSizes {
		src := make(map[int]int, size)
		for i := 0; i < size; i++ {
			src[i] = i
		}
		mc := collection.NewMap(src)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = mc.Filter(func(_ int, v int) bool { return v%2 == 0 })
			}
		})
	}
}

func BenchmarkKeyBy(b *testing.B) {
	type Item struct{ ID, Val int }
	for _, size := range benchSizes {
		items := make([]Item, size)
		for i := range items {
			items[i] = Item{i, i * 2}
		}
		c := collection.New(items)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = collection.KeyBy(c, func(item Item) int { return item.ID })
			}
		})
	}
}

func BenchmarkPluck(b *testing.B) {
	type Item struct{ ID, Val int }
	for _, size := range benchSizes {
		items := make([]Item, size)
		for i := range items {
			items[i] = Item{i, i * 2}
		}
		c := collection.New(items)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = collection.Pluck(c, func(item Item) int { return item.Val })
			}
		})
	}
}