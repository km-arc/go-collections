package collection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/km-arc/go-collections"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

func ints(ns ...int) *collection.Collection[int] {
	return collection.New(ns)
}

func strs(ss ...string) *collection.Collection[string] {
	return collection.New(ss)
}

func assertEq[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertSliceEq[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("length: got %d, want %d  (got=%v, want=%v)", len(got), len(want), got, want)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

// ─── New / Collect / Times ────────────────────────────────────────────────────

func TestNew(t *testing.T) {
	c := ints(1, 2, 3)
	assertSliceEq(t, c.All(), []int{1, 2, 3})
}

func TestCollect(t *testing.T) {
	c := collection.Collect([]string{"a", "b"})
	assertEq(t, c.Count(), 2)
}

func TestTimes(t *testing.T) {
	c := collection.Times(5, func(i int) int { return i * i })
	assertSliceEq(t, c.All(), []int{1, 4, 9, 16, 25})
}

// ─── Basic accessors ──────────────────────────────────────────────────────────

func TestCount(t *testing.T) {
	assertEq(t, ints(1, 2, 3).Count(), 3)
	assertEq(t, ints().Count(), 0)
}

func TestIsEmpty(t *testing.T) {
	assertEq(t, ints().IsEmpty(), true)
	assertEq(t, ints(1).IsEmpty(), false)
}

func TestIsNotEmpty(t *testing.T) {
	assertEq(t, ints(1).IsNotEmpty(), true)
	assertEq(t, ints().IsNotEmpty(), false)
}

func TestFirst(t *testing.T) {
	v, ok := ints(1, 2, 3).First(nil)
	assertEq(t, v, 1)
	assertEq(t, ok, true)

	v, ok = ints(1, 2, 3).First(func(x int) bool { return x > 1 })
	assertEq(t, v, 2)
	assertEq(t, ok, true)

	_, ok = ints().First(nil)
	assertEq(t, ok, false)
}

func TestFirstOrFail(t *testing.T) {
	_, err := ints().FirstOrFail(nil)
	if err == nil {
		t.Error("expected error on empty collection")
	}
	v, err := ints(5).FirstOrFail(nil)
	assertEq(t, err, nil)
	assertEq(t, v, 5)
}

func TestLast(t *testing.T) {
	v, ok := ints(1, 2, 3).Last(nil)
	assertEq(t, v, 3)
	assertEq(t, ok, true)

	v, ok = ints(1, 2, 3).Last(func(x int) bool { return x < 3 })
	assertEq(t, v, 2)
	assertEq(t, ok, true)
}

func TestNth(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3, 4, 5, 6).Nth(2, 0).All(), []int{1, 3, 5})
	assertSliceEq(t, ints(1, 2, 3, 4, 5, 6).Nth(2, 1).All(), []int{2, 4, 6})
}

func TestGet(t *testing.T) {
	v, ok := ints(10, 20, 30).Get(1)
	assertEq(t, v, 20)
	assertEq(t, ok, true)

	_, ok = ints(10, 20, 30).Get(99)
	assertEq(t, ok, false)
}

// ─── Map / Filter / Reject ────────────────────────────────────────────────────

func TestMap(t *testing.T) {
	got := ints(1, 2, 3).Map(func(v int) int { return v * 2 }).All()
	assertSliceEq(t, got, []int{2, 4, 6})
}

func TestMapWithIndex(t *testing.T) {
	got := ints(10, 20, 30).MapWithIndex(func(v, i int) int { return v + i }).All()
	assertSliceEq(t, got, []int{10, 21, 32})
}

func TestFilter(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).Filter(func(v int) bool { return v%2 == 0 }).All()
	assertSliceEq(t, got, []int{2, 4})
}

func TestReject(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).Reject(func(v int) bool { return v%2 == 0 }).All()
	assertSliceEq(t, got, []int{1, 3, 5})
}

// ─── Each / Tap / Pipe ────────────────────────────────────────────────────────

func TestEach(t *testing.T) {
	var sum int
	ints(1, 2, 3).Each(func(v int) bool {
		sum += v
		return true
	})
	assertEq(t, sum, 6)
}

func TestEachEarlyStop(t *testing.T) {
	var collected []int
	ints(1, 2, 3, 4, 5).Each(func(v int) bool {
		if v == 3 {
			return false
		}
		collected = append(collected, v)
		return true
	})
	assertSliceEq(t, collected, []int{1, 2})
}

func TestTap(t *testing.T) {
	var sideEffect int
	result := ints(1, 2, 3).
		Map(func(v int) int { return v * 2 }).
		Tap(func(c *collection.Collection[int]) { sideEffect = c.Count() }).
		All()
	assertEq(t, sideEffect, 3)
	assertSliceEq(t, result, []int{2, 4, 6})
}

func TestPipe(t *testing.T) {
	result := ints(1, 2, 3).
		Pipe(func(c *collection.Collection[int]) *collection.Collection[int] {
			return c.Map(func(v int) int { return v + 10 })
		}).All()
	assertSliceEq(t, result, []int{11, 12, 13})
}

// ─── Slice / Take / Skip ──────────────────────────────────────────────────────

func TestSlice(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3, 4, 5).Slice(1, 3).All(), []int{2, 3, 4})
	assertSliceEq(t, ints(1, 2, 3, 4, 5).Slice(3, -1).All(), []int{4, 5})
}

func TestTake(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3, 4, 5).Take(3).All(), []int{1, 2, 3})
	assertSliceEq(t, ints(1, 2, 3, 4, 5).Take(-2).All(), []int{4, 5})
}

func TestTakeUntil(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).TakeUntil(func(v int) bool { return v >= 3 }).All()
	assertSliceEq(t, got, []int{1, 2})
}

func TestTakeWhile(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).TakeWhile(func(v int) bool { return v < 3 }).All()
	assertSliceEq(t, got, []int{1, 2})
}

func TestSkip(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3, 4, 5).Skip(2).All(), []int{3, 4, 5})
}

func TestSkipUntil(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).SkipUntil(func(v int) bool { return v >= 3 }).All()
	assertSliceEq(t, got, []int{3, 4, 5})
}

func TestSkipWhile(t *testing.T) {
	got := ints(1, 2, 3, 4, 5).SkipWhile(func(v int) bool { return v < 3 }).All()
	assertSliceEq(t, got, []int{3, 4, 5})
}

// ─── Reverse / Shuffle ────────────────────────────────────────────────────────

func TestReverse(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3).Reverse().All(), []int{3, 2, 1})
}

func TestShuffle(t *testing.T) {
	c := ints(1, 2, 3, 4, 5)
	shuffled := c.Shuffle()
	// Count and content should be the same even if order differs.
	assertEq(t, shuffled.Count(), 5)
	assertEq(t, shuffled.Sum(func(v int) float64 { return float64(v) }), 15)
}

// ─── Sort ─────────────────────────────────────────────────────────────────────

func TestSort(t *testing.T) {
	got := ints(3, 1, 4, 1, 5).Sort(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{1, 1, 3, 4, 5})
}

func TestSortDesc(t *testing.T) {
	got := ints(3, 1, 4, 1, 5).SortDesc(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{5, 4, 3, 1, 1})
}

// ─── Chunk / Split ────────────────────────────────────────────────────────────

func TestChunk(t *testing.T) {
	chunks := ints(1, 2, 3, 4, 5).Chunk(2)
	if len(chunks) != 3 {
		t.Fatalf("want 3 chunks, got %d", len(chunks))
	}
	assertSliceEq(t, chunks[0].All(), []int{1, 2})
	assertSliceEq(t, chunks[1].All(), []int{3, 4})
	assertSliceEq(t, chunks[2].All(), []int{5})
}

func TestSplit(t *testing.T) {
	groups := ints(1, 2, 3, 4, 5).Split(3)
	if len(groups) != 3 {
		t.Fatalf("want 3 groups, got %d", len(groups))
	}
}

// ─── Flatten / FlatMap ────────────────────────────────────────────────────────

func TestFlatten(t *testing.T) {
	nested := collection.New([][]int{{1, 2}, {3, 4}, {5}})
	got := collection.Flatten(nested).All()
	assertSliceEq(t, got, []int{1, 2, 3, 4, 5})
}

func TestFlatMap(t *testing.T) {
	got := ints(1, 2, 3).FlatMap(func(v int) []int { return []int{v, v * v} }).All()
	assertSliceEq(t, got, []int{1, 1, 2, 4, 3, 9})
}

// ─── Unique ───────────────────────────────────────────────────────────────────

func TestUnique(t *testing.T) {
	got := ints(1, 2, 2, 3, 3, 3).Unique().Sort(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{1, 2, 3})
}

func TestUniqueBy(t *testing.T) {
	type P struct{ Name string }
	items := collection.New([]P{{"Alice"}, {"Bob"}, {"Alice"}})
	got := items.UniqueBy(func(p P) string { return p.Name }).Count()
	assertEq(t, got, 2)
}

func TestDuplicate(t *testing.T) {
	got := ints(1, 2, 2, 3, 3).Duplicate().Sort(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{2, 3})
}

// ─── Merge / Concat / Diff / Intersect ───────────────────────────────────────

func TestMerge(t *testing.T) {
	got := ints(1, 2).Merge(ints(3, 4)).All()
	assertSliceEq(t, got, []int{1, 2, 3, 4})
}

func TestConcat(t *testing.T) {
	got := ints(1).Concat(ints(2, 3), ints(4)).All()
	assertSliceEq(t, got, []int{1, 2, 3, 4})
}

func TestDiff(t *testing.T) {
	got := ints(1, 2, 3, 4).Diff(ints(2, 4)).Sort(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{1, 3})
}

func TestIntersect(t *testing.T) {
	got := ints(1, 2, 3, 4).Intersect(ints(2, 4, 6)).Sort(func(a, b int) bool { return a < b }).All()
	assertSliceEq(t, got, []int{2, 4})
}

// ─── Zip / Pad / ForPage ─────────────────────────────────────────────────────

func TestZip(t *testing.T) {
	pairs := collection.Zip(ints(1, 2, 3), strs("a", "b", "c"))
	if len(pairs) != 3 {
		t.Fatalf("want 3 pairs, got %d", len(pairs))
	}
	assertEq(t, pairs[0][0], 1)
	assertEq(t, pairs[0][1], "a")
}

func TestPad(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3).Pad(5, 0).All(), []int{1, 2, 3, 0, 0})
	assertSliceEq(t, ints(1, 2, 3).Pad(-5, 0).All(), []int{0, 0, 1, 2, 3})
}

func TestForPage(t *testing.T) {
	c := ints(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assertSliceEq(t, c.ForPage(1, 3).All(), []int{1, 2, 3})
	assertSliceEq(t, c.ForPage(2, 3).All(), []int{4, 5, 6})
	assertSliceEq(t, c.ForPage(4, 3).All(), []int{10})
}

// ─── Contains / Every / Search / Has ─────────────────────────────────────────

func TestContains(t *testing.T) {
	assertEq(t, ints(1, 2, 3).Contains(func(v int) bool { return v == 2 }), true)
	assertEq(t, ints(1, 2, 3).Contains(func(v int) bool { return v == 9 }), false)
}

func TestDoesntContain(t *testing.T) {
	assertEq(t, ints(1, 2, 3).DoesntContain(func(v int) bool { return v == 9 }), true)
}

func TestEvery(t *testing.T) {
	assertEq(t, ints(2, 4, 6).Every(func(v int) bool { return v%2 == 0 }), true)
	assertEq(t, ints(2, 3, 6).Every(func(v int) bool { return v%2 == 0 }), false)
}

func TestSearch(t *testing.T) {
	assertEq(t, strs("a", "b", "c").Search(func(v string) bool { return v == "b" }), 1)
	assertEq(t, strs("a", "b", "c").Search(func(v string) bool { return v == "z" }), -1)
}

func TestHas(t *testing.T) {
	assertEq(t, ints(1, 2, 3).Has(2), true)
	assertEq(t, ints(1, 2, 3).Has(5), false)
}

// ─── Mutable operations ───────────────────────────────────────────────────────

func TestPush(t *testing.T) {
	c := ints(1, 2)
	c.Push(3)
	assertSliceEq(t, c.All(), []int{1, 2, 3})
}

func TestPrepend(t *testing.T) {
	c := ints(2, 3)
	c.Prepend(1)
	assertSliceEq(t, c.All(), []int{1, 2, 3})
}

func TestPut(t *testing.T) {
	c := ints(1, 2, 3)
	_ = c.Put(1, 99)
	assertEq(t, c.All()[1], 99)
}

func TestForget(t *testing.T) {
	c := ints(1, 2, 3)
	_ = c.Forget(1)
	assertSliceEq(t, c.All(), []int{1, 3})
}

func TestPop(t *testing.T) {
	c := ints(1, 2, 3)
	v, ok := c.Pop()
	assertEq(t, v, 3)
	assertEq(t, ok, true)
	assertEq(t, c.Count(), 2)
}

func TestShift(t *testing.T) {
	c := ints(1, 2, 3)
	v, ok := c.Shift()
	assertEq(t, v, 1)
	assertEq(t, ok, true)
	assertEq(t, c.Count(), 2)
}

func TestPull(t *testing.T) {
	c := ints(10, 20, 30)
	v, ok := c.Pull(1)
	assertEq(t, v, 20)
	assertEq(t, ok, true)
	assertSliceEq(t, c.All(), []int{10, 30})
}

func TestSplice(t *testing.T) {
	c := ints(1, 2, 3, 4, 5)
	removed := c.Splice(2, 2)
	assertSliceEq(t, removed.All(), []int{3, 4})
	assertSliceEq(t, c.All(), []int{1, 2, 5})
}

func TestSpliceWithReplacement(t *testing.T) {
	c := ints(1, 2, 3, 4, 5)
	removed := c.Splice(1, 2, 10, 11)
	assertSliceEq(t, removed.All(), []int{2, 3})
	assertSliceEq(t, c.All(), []int{1, 10, 11, 4, 5})
}

func TestTransform(t *testing.T) {
	c := ints(1, 2, 3)
	c.Transform(func(v int) int { return v * 2 })
	assertSliceEq(t, c.All(), []int{2, 4, 6})
}

// ─── Reduce / Sum / Avg / Min / Max / Median / Mode ──────────────────────────

func TestReduce(t *testing.T) {
	sum := ints(1, 2, 3, 4).Reduce(0, func(carry, v int) int { return carry + v })
	assertEq(t, sum, 10)
}

func TestSum(t *testing.T) {
	assertEq(t, ints(1, 2, 3, 4).Sum(func(v int) float64 { return float64(v) }), 10.0)
}

func TestAvg(t *testing.T) {
	assertEq(t, ints(1, 2, 3, 4).Avg(func(v int) float64 { return float64(v) }), 2.5)
	assertEq(t, ints().Avg(func(v int) float64 { return float64(v) }), 0.0)
}

func TestMin(t *testing.T) {
	v, ok := ints(3, 1, 4, 1, 5).Min(func(v int) float64 { return float64(v) })
	assertEq(t, ok, true)
	assertEq(t, v, 1.0)

	_, ok = ints().Min(func(v int) float64 { return float64(v) })
	assertEq(t, ok, false)
}

func TestMax(t *testing.T) {
	v, ok := ints(3, 1, 4, 1, 5).Max(func(v int) float64 { return float64(v) })
	assertEq(t, ok, true)
	assertEq(t, v, 5.0)
}

func TestMedian(t *testing.T) {
	v, _ := ints(1, 2, 3, 4, 5).Median(func(v int) float64 { return float64(v) })
	assertEq(t, v, 3.0)

	v, _ = ints(1, 2, 3, 4).Median(func(v int) float64 { return float64(v) })
	assertEq(t, v, 2.5)
}

func TestMode(t *testing.T) {
	mode := ints(1, 2, 2, 3, 3, 3).Mode(func(v int) float64 { return float64(v) })
	assertSliceEq(t, mode, []float64{3})
}

// ─── GroupBy / Partition ──────────────────────────────────────────────────────

func TestGroupBy(t *testing.T) {
	groups := ints(1, 2, 3, 4, 5, 6).GroupBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	assertEq(t, groups["even"].Count(), 3)
	assertEq(t, groups["odd"].Count(), 3)
}

func TestPartition(t *testing.T) {
	even, odd := ints(1, 2, 3, 4).Partition(func(v int) bool { return v%2 == 0 })
	assertSliceEq(t, even.All(), []int{2, 4})
	assertSliceEq(t, odd.All(), []int{1, 3})
}

// ─── When / Unless ────────────────────────────────────────────────────────────

func TestWhen(t *testing.T) {
	got := ints(1, 2, 3).When(true, func(c *collection.Collection[int]) *collection.Collection[int] {
		return c.Map(func(v int) int { return v * 10 })
	}).All()
	assertSliceEq(t, got, []int{10, 20, 30})

	got = ints(1, 2, 3).When(false, func(c *collection.Collection[int]) *collection.Collection[int] {
		return c.Map(func(v int) int { return v * 10 })
	}).All()
	assertSliceEq(t, got, []int{1, 2, 3})
}

func TestUnless(t *testing.T) {
	got := ints(1, 2, 3).Unless(false, func(c *collection.Collection[int]) *collection.Collection[int] {
		return c.Map(func(v int) int { return v * 10 })
	}).All()
	assertSliceEq(t, got, []int{10, 20, 30})
}

func TestWhenEmpty(t *testing.T) {
	called := false
	ints(1).WhenEmpty(func(c *collection.Collection[int]) *collection.Collection[int] {
		called = true
		return c
	})
	assertEq(t, called, false)

	ints().WhenEmpty(func(c *collection.Collection[int]) *collection.Collection[int] {
		called = true
		return c
	})
	assertEq(t, called, true)
}

// ─── Implode ──────────────────────────────────────────────────────────────────

func TestImplode(t *testing.T) {
	assertEq(t, ints(1, 2, 3).Implode(", "), "1, 2, 3")
}

func TestImplodeWith(t *testing.T) {
	got := strs("a", "b", "c").ImplodeWith("-", strings.ToUpper)
	assertEq(t, got, "A-B-C")
}

// ─── Random ───────────────────────────────────────────────────────────────────

func TestRandom(t *testing.T) {
	r := ints(1, 2, 3, 4, 5).Random(3)
	assertEq(t, r.Count(), 3)
}

// ─── Conversion ───────────────────────────────────────────────────────────────

func TestToJSON(t *testing.T) {
	j, err := ints(1, 2, 3).ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	assertEq(t, j, "[1,2,3]")
}

func TestClone(t *testing.T) {
	orig := ints(1, 2, 3)
	clone := orig.Clone()
	clone.Push(4)
	assertEq(t, orig.Count(), 3) // original unaffected
	assertEq(t, clone.Count(), 4)
}

func TestValues(t *testing.T) {
	assertSliceEq(t, ints(1, 2, 3).Values().All(), []int{1, 2, 3})
}

func TestKeys(t *testing.T) {
	assertSliceEq(t, ints(10, 20, 30).Keys().All(), []int{0, 1, 2})
}

// ─── Chaining (integration) ───────────────────────────────────────────────────

func TestChaining(t *testing.T) {
	// mirrors the Laravel README example:
	// collect([1..10]).filter(even).map(*2).take(3)
	got := ints(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
		Filter(func(v int) bool { return v%2 == 0 }).
		Map(func(v int) int { return v * 2 }).
		Take(3).
		All()
	assertSliceEq(t, got, []int{4, 8, 12})
}

// ─────────────────────────────────────────────────────────────────────────────
// MapCollection tests
// ─────────────────────────────────────────────────────────────────────────────

func TestNewMap(t *testing.T) {
	mc := collection.NewMap(map[string]int{"a": 1, "b": 2})
	assertEq(t, mc.Count(), 2)
	assertEq(t, mc.Has("a"), true)
	assertEq(t, mc.Has("z"), false)
}

func TestMapCollectionGet(t *testing.T) {
	mc := collection.NewMap(map[string]int{"x": 42})
	v, ok := mc.Get("x")
	assertEq(t, v, 42)
	assertEq(t, ok, true)

	_, ok = mc.Get("missing")
	assertEq(t, ok, false)
}

func TestMapCollectionPutForget(t *testing.T) {
	mc := collection.NewMap(map[string]int{})
	mc.Put("k", 99)
	v, _ := mc.Get("k")
	assertEq(t, v, 99)
	mc.Forget("k")
	assertEq(t, mc.Has("k"), false)
}

func TestMapCollectionFilter(t *testing.T) {
	mc := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
	filtered := mc.Filter(func(_ string, v int) bool { return v > 1 })
	assertEq(t, filtered.Count(), 2)
}

func TestMapCollectionOnly(t *testing.T) {
	mc := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
	only := mc.Only("a", "c")
	assertEq(t, only.Count(), 2)
	assertEq(t, only.Has("b"), false)
}

func TestMapCollectionExcept(t *testing.T) {
	mc := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
	except := mc.Except("b")
	assertEq(t, except.Count(), 2)
	assertEq(t, except.Has("b"), false)
}

func TestMapCollectionMerge(t *testing.T) {
	a := collection.NewMap(map[string]int{"a": 1})
	b := collection.NewMap(map[string]int{"b": 2, "a": 99})
	merged := a.Merge(b)
	v, _ := merged.Get("a")
	assertEq(t, v, 99) // b overwrites a
	assertEq(t, merged.Count(), 2)
}

func TestMapCollectionUnion(t *testing.T) {
	a := collection.NewMap(map[string]int{"a": 1})
	b := collection.NewMap(map[string]int{"b": 2, "a": 99})
	union := a.Union(b)
	v, _ := union.Get("a")
	assertEq(t, v, 1) // a wins in union
	assertEq(t, union.Count(), 2)
}

func TestFlip(t *testing.T) {
	mc := collection.NewMap(map[string]string{"hello": "world", "foo": "bar"})
	flipped := collection.Flip(mc)
	assertEq(t, flipped.Has("world"), true)
	assertEq(t, flipped.Has("hello"), false)
}

func TestCombine(t *testing.T) {
	keys := strs("a", "b", "c")
	vals := ints(1, 2, 3)
	mc := collection.Combine(keys, vals)
	v, ok := mc.Get("b")
	assertEq(t, ok, true)
	assertEq(t, v, 2)
}

func TestKeyBy(t *testing.T) {
	type U struct{ ID int }
	users := collection.New([]U{{1}, {2}, {3}})
	mc := collection.KeyBy(users, func(u U) int { return u.ID })
	u, ok := mc.Get(2)
	assertEq(t, ok, true)
	assertEq(t, u.ID, 2)
}

func TestPluck(t *testing.T) {
	type P struct{ Name string }
	items := collection.New([]P{{"Alice"}, {"Bob"}, {"Carol"}})
	names := collection.Pluck(items, func(p P) string { return p.Name }).All()
	assertSliceEq(t, names, []string{"Alice", "Bob", "Carol"})
}

func TestMapCollectionToJSON(t *testing.T) {
	mc := collection.NewMap(map[string]int{"a": 1})
	j, err := mc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	assertEq(t, j, `{"a":1}`)
}

// ─── Example / smoke test ─────────────────────────────────────────────────────

func ExampleNew() {
	c := collection.New([]int{1, 2, 3, 4, 5})
	even := c.Filter(func(v int) bool { return v%2 == 0 })
	fmt.Println(even.All())
	// Output: [2 4]
}

func ExampleCollection_Map() {
	doubled := collection.New([]int{1, 2, 3}).Map(func(v int) int { return v * 2 })
	fmt.Println(doubled.All())
	// Output: [2 4 6]
}

func ExampleCollection_GroupBy() {
	groups := collection.New([]int{1, 2, 3, 4, 5, 6}).GroupBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Println("even:", groups["even"].Count())
	fmt.Println("odd:", groups["odd"].Count())
	// Output:
	// even: 3
	// odd: 3
}
