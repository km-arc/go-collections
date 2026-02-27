package main

import (
	"fmt"
	"github.com/km-arc/go-collections/collection"
)


// ─── Package-level ────────────────────────────────────────────────────────────

func Example() {
	// A realistic pipeline: paginate a list of user IDs, keep the active
	// ones, extract their names, then join for display.
	type User struct {
		ID     int
		Name   string
		Active bool
	}

	users := collection.New([]User{
		{1, "Alice", true},
		{2, "Bob", false},
		{3, "Carol", true},
		{4, "Dave", true},
		{5, "Eve", false},
	})

	page := users.
		Filter(func(u User) bool { return u.Active }).
		ForPage(1, 2)

	names := collection.Pluck(page, func(u User) string { return u.Name })

	fmt.Println(names.Implode(", "))
	// Output: Alice, Carol
}

// ─── Construction ─────────────────────────────────────────────────────────────

func ExampleNew() {
	c := collection.New([]string{"taylor", "abigail", "james"})
	fmt.Println(c.Count())
	// Output: 3
}

func ExampleCollect() {
	c := collection.Collect([]int{1, 2, 3})
	fmt.Println(c.All())
	// Output: [1 2 3]
}

func ExampleTimes() {
	c := collection.Times(5, func(i int) string {
		return fmt.Sprintf("item %d", i)
	})
	fmt.Println(c.All())
	// Output: [item 1 item 2 item 3 item 4 item 5]
}

// ─── Transformation ───────────────────────────────────────────────────────────

func ExampleCollection_Map() {
	result := collection.New([]int{1, 2, 3, 4}).
		Map(func(v int) int { return v * v }).
		All()
	fmt.Println(result)
	// Output: [1 4 9 16]
}

func ExampleCollection_Filter() {
	result := collection.New([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(v int) bool { return v%2 == 0 }).
		All()
	fmt.Println(result)
	// Output: [2 4 6]
}

func ExampleCollection_Reject() {
	result := collection.New([]int{1, 2, 3, 4, 5}).
		Reject(func(v int) bool { return v%2 == 0 }).
		All()
	fmt.Println(result)
	// Output: [1 3 5]
}

func ExampleCollection_FlatMap() {
	result := collection.New([]int{1, 2, 3}).
		FlatMap(func(v int) []int { return []int{v, v * v} }).
		All()
	fmt.Println(result)
	// Output: [1 1 2 4 3 9]
}

func ExampleFlatten() {
	nested := collection.New([][]int{{1, 2}, {3, 4}, {5}})
	fmt.Println(collection.Flatten(nested).All())
	// Output: [1 2 3 4 5]
}

// ─── Slicing ──────────────────────────────────────────────────────────────────

func ExampleCollection_Take() {
	c := collection.New([]int{1, 2, 3, 4, 5})
	fmt.Println(c.Take(3).All())
	fmt.Println(c.Take(-2).All())
	// Output:
	// [1 2 3]
	// [4 5]
}

func ExampleCollection_Skip() {
	fmt.Println(collection.New([]int{1, 2, 3, 4, 5}).Skip(2).All())
	// Output: [3 4 5]
}

func ExampleCollection_Chunk() {
	chunks := collection.New([]int{1, 2, 3, 4, 5}).Chunk(2)
	for _, ch := range chunks {
		fmt.Println(ch.All())
	}
	// Output:
	// [1 2]
	// [3 4]
	// [5]
}

func ExampleCollection_ForPage() {
	c := collection.New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(c.ForPage(1, 3).All())
	fmt.Println(c.ForPage(2, 3).All())
	fmt.Println(c.ForPage(3, 3).All())
	// Output:
	// [1 2 3]
	// [4 5 6]
	// [7 8 9]
}

// ─── Searching ────────────────────────────────────────────────────────────────

func ExampleCollection_First() {
	c := collection.New([]int{1, 2, 3, 4, 5})

	// No predicate — returns first element.
	v, _ := c.First(nil)
	fmt.Println(v)

	// With predicate.
	v, _ = c.First(func(n int) bool { return n > 3 })
	fmt.Println(v)
	// Output:
	// 1
	// 4
}

func ExampleCollection_Last() {
	c := collection.New([]int{1, 2, 3, 4, 5})
	v, _ := c.Last(func(n int) bool { return n < 4 })
	fmt.Println(v)
	// Output: 3
}

func ExampleCollection_Search() {
	idx := collection.New([]string{"apple", "banana", "cherry"}).
		Search(func(s string) bool { return s == "banana" })
	fmt.Println(idx)
	// Output: 1
}

func ExampleCollection_Contains() {
	c := collection.New([]int{1, 2, 3})
	fmt.Println(c.Contains(func(v int) bool { return v == 2 }))
	fmt.Println(c.Contains(func(v int) bool { return v == 9 }))
	// Output:
	// true
	// false
}

// ─── Aggregates ───────────────────────────────────────────────────────────────

func ExampleCollection_Sum() {
	total := collection.New([]int{1, 2, 3, 4, 5}).
		Sum(func(v int) float64 { return float64(v) })
	fmt.Println(total)
	// Output: 15
}

func ExampleCollection_Avg() {
	avg := collection.New([]int{1, 2, 3, 4}).
		Avg(func(v int) float64 { return float64(v) })
	fmt.Println(avg)
	// Output: 2.5
}

func ExampleCollection_Reduce() {
	product := collection.New([]int{1, 2, 3, 4, 5}).
		Reduce(1, func(carry, v int) int { return carry * v })
	fmt.Println(product)
	// Output: 120
}

func ExampleCollection_Min() {
	min, _ := collection.New([]int{5, 3, 8, 1, 9}).
		Min(func(v int) float64 { return float64(v) })
	fmt.Println(min)
	// Output: 1
}

func ExampleCollection_Max() {
	max, _ := collection.New([]int{5, 3, 8, 1, 9}).
		Max(func(v int) float64 { return float64(v) })
	fmt.Println(max)
	// Output: 9
}

func ExampleCollection_Median() {
	// Odd count
	m, _ := collection.New([]int{1, 2, 3, 4, 5}).
		Median(func(v int) float64 { return float64(v) })
	fmt.Println(m)

	// Even count — average of two middles
	m, _ = collection.New([]int{1, 2, 3, 4}).
		Median(func(v int) float64 { return float64(v) })
	fmt.Println(m)
	// Output:
	// 3
	// 2.5
}

// ─── Grouping / Partitioning ──────────────────────────────────────────────────

func ExampleCollection_GroupBy() {
	groups := collection.New([]int{1, 2, 3, 4, 5, 6}).GroupBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Println("even:", groups["even"].Sort(func(a, b int) bool { return a < b }).All())
	fmt.Println("odd:", groups["odd"].Sort(func(a, b int) bool { return a < b }).All())
	// Output:
	// even: [2 4 6]
	// odd: [1 3 5]
}

func ExampleCollection_Partition() {
	pass, fail := collection.New([]int{1, 2, 3, 4, 5, 6}).
		Partition(func(v int) bool { return v%2 == 0 })
	fmt.Println("even:", pass.All())
	fmt.Println("odd:", fail.All())
	// Output:
	// even: [2 4 6]
	// odd: [1 3 5]
}

// ─── Sorting ──────────────────────────────────────────────────────────────────

func ExampleCollection_Sort() {
	result := collection.New([]int{3, 1, 4, 1, 5, 9, 2, 6}).
		Sort(func(a, b int) bool { return a < b }).
		All()
	fmt.Println(result)
	// Output: [1 1 2 3 4 5 6 9]
}

func ExampleCollection_SortDesc() {
	result := collection.New([]int{3, 1, 4, 1, 5, 9, 2, 6}).
		SortDesc(func(a, b int) bool { return a < b }).
		All()
	fmt.Println(result)
	// Output: [9 6 5 4 3 2 1 1]
}

// ─── Set operations ───────────────────────────────────────────────────────────

func ExampleCollection_Unique() {
	result := collection.New([]int{1, 2, 2, 3, 3, 3}).
		Unique().
		Sort(func(a, b int) bool { return a < b }).
		All()
	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleCollection_Diff() {
	result := collection.New([]int{1, 2, 3, 4, 5}).
		Diff(collection.New([]int{2, 4})).
		All()
	fmt.Println(result)
	// Output: [1 3 5]
}

func ExampleCollection_Intersect() {
	result := collection.New([]int{1, 2, 3, 4, 5}).
		Intersect(collection.New([]int{2, 4, 6})).
		All()
	fmt.Println(result)
	// Output: [2 4]
}

// ─── Conditionals ────────────────────────────────────────────────────────────

func ExampleCollection_When() {
	c := collection.New([]int{1, 2, 3})
	result := c.When(true, func(c *collection.Collection[int]) *collection.Collection[int] {
		return c.Map(func(v int) int { return v * 10 })
	})
	fmt.Println(result.All())
	// Output: [10 20 30]
}

func ExampleCollection_WhenEmpty() {
	empty := collection.New([]int{})
	result := empty.WhenEmpty(func(c *collection.Collection[int]) *collection.Collection[int] {
		return collection.New([]int{0})
	})
	fmt.Println(result.All())
	// Output: [0]
}

// ─── Pad / Reverse / Nth ──────────────────────────────────────────────────────

func ExampleCollection_Pad() {
	fmt.Println(collection.New([]int{1, 2, 3}).Pad(5, 0).All())
	fmt.Println(collection.New([]int{1, 2, 3}).Pad(-5, 0).All())
	// Output:
	// [1 2 3 0 0]
	// [0 0 1 2 3]
}

func ExampleCollection_Reverse() {
	fmt.Println(collection.New([]int{1, 2, 3, 4, 5}).Reverse().All())
	// Output: [5 4 3 2 1]
}

func ExampleCollection_Nth() {
	fmt.Println(collection.New([]int{1, 2, 3, 4, 5, 6}).Nth(2, 0).All())
	// Output: [1 3 5]
}

// ─── String ───────────────────────────────────────────────────────────────────

func ExampleCollection_Implode() {
	fmt.Println(collection.New([]int{1, 2, 3}).Implode(", "))
	// Output: 1, 2, 3
}

func ExampleCollection_ImplodeWith() {
	result := collection.New([]string{"hello", "world"}).
		ImplodeWith(" | ", strings.ToUpper)
	fmt.Println(result)
	// Output: HELLO | WORLD
}

// ─── ToJSON ───────────────────────────────────────────────────────────────────

func ExampleCollection_ToJSON() {
	j, _ := collection.New([]int{1, 2, 3}).ToJSON()
	fmt.Println(j)
	// Output: [1,2,3]
}

// ─── Zip ──────────────────────────────────────────────────────────────────────

func ExampleZip() {
	pairs := collection.Zip(
		collection.New([]int{1, 2, 3}),
		collection.New([]string{"a", "b", "c"}),
	)
	for _, p := range pairs {
		fmt.Printf("%v:%v ", p[0], p[1])
	}
	fmt.Println()
	// Output: 1:a 2:b 3:c
}

// ─── MapCollection ────────────────────────────────────────────────────────────

func ExampleNewMap() {
	mc := collection.NewMap(map[string]int{"apples": 5, "bananas": 3})
	fmt.Println(mc.Count())
	v, _ := mc.Get("apples")
	fmt.Println(v)
	// Output:
	// 2
	// 5
}

func ExampleKeyBy() {
	type User struct {
		ID   int
		Name string
	}
	users := collection.New([]User{{1, "Alice"}, {2, "Bob"}, {3, "Carol"}})
	mc := collection.KeyBy(users, func(u User) int { return u.ID })

	u, ok := mc.Get(2)
	fmt.Println(ok, u.Name)
	// Output: true Bob
}

func ExamplePluck() {
	type Product struct {
		Name  string
		Price float64
	}
	products := collection.New([]Product{
		{"Widget", 9.99},
		{"Gadget", 24.99},
		{"Doohickey", 4.99},
	})
	names := collection.Pluck(products, func(p Product) string { return p.Name })
	fmt.Println(names.All())
	// Output: [Widget Gadget Doohickey]
}

func ExampleCombine() {
	keys := collection.New([]string{"name", "city", "lang"})
	vals := collection.New([]string{"Alice", "Detroit", "Go"})
	mc := collection.Combine(keys, vals)

	v, _ := mc.Get("lang")
	fmt.Println(v)
	// Output: Go
}

func ExampleFlip() {
	mc := collection.NewMap(map[string]string{
		"en": "English",
		"fr": "French",
	})
	flipped := collection.Flip(mc)
	v, _ := flipped.Get("English")
	fmt.Println(v)
	// Output: en
}

func ExampleMapCollection_Only() {
	mc := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
	only := mc.Only("a", "c")
	fmt.Println(only.Count())
	fmt.Println(only.Has("b"))
	// Output:
	// 2
	// false
}

func ExampleMapCollection_Union() {
	a := collection.NewMap(map[string]int{"a": 1, "b": 2})
	b := collection.NewMap(map[string]int{"b": 99, "c": 3})
	u := a.Union(b) // a's keys win
	v, _ := u.Get("b")
	fmt.Println(v)
	// Output: 2
}