package evalostic

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
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

var ev *Evalostic

var validCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(chars int) string {
	res := make([]byte, chars)
	for i := 0; i < chars; i++ {
		res[i] = validCharacters[rand.Intn(len(validCharacters))]
	}
	return string(res)
}

func randomCondition(maxConcatenations int) string {
	var cond string
	if rand.Intn(5) == 0 {
		cond += "NOT "
	}
	if cond != "" && !strings.HasSuffix(cond, " ") {
		cond += " "
	}
	concatenations := rand.Intn(maxConcatenations)
	if concatenations == 0 {
		cond += strconv.Quote(randomString(3 + rand.Intn(20)))
		if rand.Intn(2) == 0 {
			cond += "i"
		}
		return cond
	}
	cond += "("
	for i := 0; i < concatenations+1; i++ {
		if i > 0 {
			cond += []string{" AND ", " OR "}[rand.Intn(2)]
		}
		cond += randomCondition(maxConcatenations - 1)
	}
	cond += ")"
	return cond
}

func Example_randomString() {
	rand.Seed(0)
	fmt.Println(randomString(10))
	fmt.Println(randomString(10))
	// Output:
	// mUNERA9rI2
	// cvTK4UHomc
}

func Example_randomCondition() {
	rand.Seed(0)
	fmt.Println(randomCondition(3))
	fmt.Println(randomCondition(3))
	fmt.Println(randomCondition(3))
	fmt.Println(randomCondition(3))
	fmt.Println(randomCondition(3))
	fmt.Println(randomCondition(3))
	// Output:
	// "ERA9rI2cvTK4UHom"i
	// NOT "QvymkzADm"
	// ("HwxmE4tL20SrW"i AND (NOT "U86R7wIBbUt9RwI9U" OR "aWsz0l"i) AND "egogMR6spJPZHaPT0w4"i)
	// "nRPswXn0i6jt3PARUDb0oU"
	// NOT "aT3LnXs1oH7gq"i
	// (NOT ("K4k"i OR NOT "UcTwCQmwFJPivbFQtOWS") OR NOT "8bmfKkaauLI"i)
}

func BenchmarkNew(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		ev, err = New([]string{`"foo" AND ("bar" OR NOT ("baz" AND NOT "qux"))`})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkNewN(b *testing.B, n int) {
	rand.Seed(0)
	conds := make([]string, n)
	for i := 0; i < n; i++ {
		conds[i] = randomCondition(3)
		if i < 10 {
			b.Log(conds[i])
		}
		if i == 10 {
			b.Logf("[skipping remaining %d conditions in outout]", n-10)
		}
	}
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		ev, err = New(conds)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNew_10(b *testing.B)     { benchmarkNewN(b, 10) }
func BenchmarkNew_100(b *testing.B)    { benchmarkNewN(b, 100) }
func BenchmarkNew_1000(b *testing.B)   { benchmarkNewN(b, 1000) }
func BenchmarkNew_10000(b *testing.B)  { benchmarkNewN(b, 10000) }
func BenchmarkNew_100000(b *testing.B) { benchmarkNewN(b, 100000) }

var matches []int

func benchmarkMatchN(b *testing.B, n int) {
	rand.Seed(0)
	conds := make([]string, n)
	for i := 0; i < n; i++ {
		conds[i] = randomCondition(3)
	}
	ev, err := New(conds)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matches = ev.Match("ERA9rI2cvTK4UHomQvymkzADmHwxmE4tL20SrWU86R7wIBbUt9RwI9UaWsz0legogMR6spJPZHaPT0w4n" +
			"RPswXn0i6jt3PARUDb0oUaT3LnXs1oH7gqK4kUcTwCQmwFJPivbFQtOWS8bmfKkaauLIHID5QOEwEYos1y7eNI1Had7lItaey9Y" +
			"WwoRhgmpWQc9DYV9uRIl8ILMGrvME4e8vdKlJUpEMrZqXlqPLn0eqyWgWIyTrJVVYSzaAAuZrNuNtZCN")
	}
	if len(matches) > 20 {
		b.Logf("matches: %+v...+%d matches", matches[:20], len(matches)-20)
	} else {
		b.Logf("matches: %+v", matches)
	}
}

func BenchmarkEvalostic_Match_10(b *testing.B)     { benchmarkMatchN(b, 10) }
func BenchmarkEvalostic_Match_100(b *testing.B)    { benchmarkMatchN(b, 100) }
func BenchmarkEvalostic_Match_1000(b *testing.B)   { benchmarkMatchN(b, 1000) }
func BenchmarkEvalostic_Match_10000(b *testing.B)  { benchmarkMatchN(b, 10000) }
func BenchmarkEvalostic_Match_100000(b *testing.B) { benchmarkMatchN(b, 100000) }
