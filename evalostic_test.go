package evalostic

import (
	"fmt"
	"testing"

	testify "github.com/stretchr/testify/assert"
)

func TestEvalostic(t *testing.T) {
	assert := testify.New(t)
	e, err := New([]string{
		`"foo" OR "bar"`,
		`"baz" AND "qux"`,
		`("a" OR "b") AND ("c" OR "d")`,
		`"1" AND NOT "2"`,
		`"1" AND NOT "2"`,
	})
	assert.NoError(err)
	assert.ElementsMatch(e.Match("foo"), []int{0})
	assert.ElementsMatch(e.Match("bar"), []int{0})
	assert.ElementsMatch(e.Match("foo bar"), []int{0})
	assert.ElementsMatch(e.Match("baz"), nil)
	assert.ElementsMatch(e.Match("baz qux"), []int{1})
	assert.ElementsMatch(e.Match("qux baz"), []int{1})
	assert.ElementsMatch(e.Match("ab"), nil)
	assert.ElementsMatch(e.Match("ac"), []int{2})
	assert.ElementsMatch(e.Match("ad"), []int{2})
	assert.ElementsMatch(e.Match("bc"), []int{2})
	assert.ElementsMatch(e.Match("bd"), []int{2})
	assert.ElementsMatch(e.Match("cd"), nil)
	assert.ElementsMatch(e.Match("abcd"), []int{2})
	assert.ElementsMatch(e.Match("1"), []int{3, 4})
	assert.ElementsMatch(e.Match("2"), nil)
	assert.ElementsMatch(e.Match("12"), nil)
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
