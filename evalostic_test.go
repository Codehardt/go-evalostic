package evalostic

import (
	"fmt"
	"runtime"
	"testing"
)

func assert(t *testing.T, b bool) {
	_, f, l, _ := runtime.Caller(1)
	if !b {
		t.Fatalf("assertion failed in %s:%d", f, l)
	}
}

func sameIntegers(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestEvalostic(t *testing.T) {
	e, err := New([]string{
		`"foo" OR "bar"`,
		`"baz" AND "qux"`,
		`("a" OR "b") AND ("c" OR "d")`,
		`"1" AND NOT "2"`,
		`"1" AND NOT "2"`,
	})
	assert(t, err == nil)
	assert(t, sameIntegers(e.Match("foo"), []int{0}))
	assert(t, sameIntegers(e.Match("bar"), []int{0}))
	assert(t, sameIntegers(e.Match("foo bar"), []int{0}))
	assert(t, sameIntegers(e.Match("baz"), []int{}))
	assert(t, sameIntegers(e.Match("baz qux"), []int{1}))
	assert(t, sameIntegers(e.Match("qux baz"), []int{1}))
	assert(t, sameIntegers(e.Match("ab"), []int{}))
	assert(t, sameIntegers(e.Match("ac"), []int{2}))
	assert(t, sameIntegers(e.Match("ad"), []int{2}))
	assert(t, sameIntegers(e.Match("bc"), []int{2}))
	assert(t, sameIntegers(e.Match("bd"), []int{2}))
	assert(t, sameIntegers(e.Match("cd"), []int{}))
	assert(t, sameIntegers(e.Match("abcd"), []int{2}))
	assert(t, sameIntegers(e.Match("1"), []int{3, 4}))
	assert(t, sameIntegers(e.Match("2"), []int{}))
	assert(t, sameIntegers(e.Match("12"), []int{}))
}

func ExampleMatch() {
	e, err := New([]string{
		`"foo" OR "bar"`,
		`NOT "foo" AND ("bar" OR "baz")`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(e.Match("foo"))
	fmt.Println(e.Match("bar"))
	fmt.Println(e.Match("foobar"))
	fmt.Println(e.Match("baz"))
	fmt.Println(e.Match("qux"))
	// Output:
	// [0]
	// [0 1]
	// [0]
	// [1]
	// []
}

func ExampleMatch_Negative() {
	e, err := New([]string{
		`NOT "foo" AND NOT "bar"`,
		`NOT ("foo" AND "bar" AND "baz")`,
		`"foo" OR NOT "baz"`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(e.Match("foo"))
	fmt.Println(e.Match("bar"))
	fmt.Println(e.Match("baz"))
	fmt.Println(e.Match("foo bar"))
	fmt.Println(e.Match("foo bar baz"))
	fmt.Println(e.Match("qux"))
	// Output:
	// [1 2]
	// [1 2]
	// [0 1]
	// [1 2]
	// [2]
	// [0 1 2]
}

func ExampleNegatives() {
	e, err := New([]string{
		`"foo"`,
		`NOT "foo"`,
		`"foo" AND NOT "bar"`,
		`"foo" OR NOT "bar"`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(e.Negatives())
	// Output:
	// [1 3]
}

func ExampleMatch_CaseInsensitive() {
	e, err := New([]string{
		`"FOO"i AND "bar"`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(e.Match("foo bar"))
	fmt.Println(e.Match("FoO BaR"))
	fmt.Println(e.Match("FOO BAR"))
	fmt.Println(e.Match("FoO bar"))
	// Output:
	// [0]
	// []
	// []
	// [0]
}
