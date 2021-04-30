package evalostic

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type Condition interface {
	Match(str string) bool
	GetStrings(into map[string]struct{})
	ToCondition() string
}

type NOT struct{ Sub Condition }

func (n NOT) ToCondition() string {
	if isANDorOR(n.Sub) {
		return fmt.Sprintf("NOT (%s)", n.Sub.ToCondition())
	}
	return fmt.Sprintf("NOT %s", n.Sub.ToCondition())
}

func (n NOT) GetStrings(into map[string]struct{}) {
	n.Sub.GetStrings(into)
}

func (n NOT) Match(str string) bool {
	return !n.Sub.Match(str)
}

type AND struct {
	Sub []Condition
}

func isANDorOR(sub Condition) bool {
	_, isAND := sub.(AND)
	_, isOR := sub.(OR)
	return isAND || isOR
}

func (a AND) GetStrings(into map[string]struct{}) {
	for _, sub := range a.Sub {
		sub.GetStrings(into)
	}
}

func (a AND) ToCondition() string {
	res := make([]string, len(a.Sub))
	for i, sub := range a.Sub {
		if isANDorOR(sub) {
			res[i] = fmt.Sprintf("(%s)", sub.ToCondition())
		} else {
			res[i] = sub.ToCondition()
		}
	}
	return fmt.Sprintf("%s", strings.Join(res, " AND "))
}

func (a AND) Match(str string) bool {
	for _, sub := range a.Sub {
		if !sub.Match(str) {
			return false
		}
	}
	return true
}

type OR struct {
	Sub []Condition
}

func (a OR) ToCondition() string {
	res := make([]string, len(a.Sub))
	for i, sub := range a.Sub {
		if isANDorOR(sub) {
			res[i] = fmt.Sprintf("(%s)", sub.ToCondition())
		} else {
			res[i] = sub.ToCondition()
		}
	}
	return fmt.Sprintf("%s", strings.Join(res, " OR "))
}
func (a OR) GetStrings(into map[string]struct{}) {
	for _, sub := range a.Sub {
		sub.GetStrings(into)
	}
}
func or(sub ...Condition) Condition  { return OR{Sub: sub} }
func and(sub ...Condition) Condition { return AND{Sub: sub} }
func not(sub Condition) Condition    { return NOT{Sub: sub} }
func val(str string) Condition       { return VALUE(str) }
func vali(str string) Condition      { return VALUEI(str) }

func (a OR) Match(str string) bool {
	for _, sub := range a.Sub {
		if sub.Match(str) {
			return true
		}
	}
	return len(a.Sub) == 0
}

type VALUE string // Case Sensitive

func (v VALUE) Match(str string) bool {
	return strings.Contains(str, string(v))
}

func (v VALUE) ToCondition() string {
	return strconv.Quote(string(v))
}

func (v VALUE) GetStrings(into map[string]struct{}) {
	into[string(v)] = struct{}{}
}

type VALUEI string // Case Insensitive

func (v VALUEI) Match(str string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(string(v)))
}

func (v VALUEI) ToCondition() string {
	return strconv.Quote(string(v)) + "i"
}

func (v VALUEI) GetStrings(into map[string]struct{}) {
	into[string(v)] = struct{}{}
}

// TheTestConditions is used to test conditions in a tree structure against conditions that were
// parsed with evalostic. It will create a big Condition Comparison Table.
// Rules for the condition below:
// 1.) Do not use upper case strings, they will be automatically added in the tests
// 2.) Do not use more than 10 different strings in the complete conditions
var TheTestConditions = []Condition{
	// START

	val("a"),

	vali("a"),

	or(
		val("a"),
		val("b"),
	),

	or(
		vali("a"),
		vali("b"),
	),

	and(
		val("a"),
		val("b"),
	),

	and(
		vali("a"),
		vali("b"),
	),

	not(
		val("a"),
	),

	not(
		vali("a"),
	),

	not(
		and(
			val("a"),
			val("b"),
		),
	),

	not(
		or(
			val("a"),
			val("b"),
		),
	),

	and(
		or(
			val("a"),
			val("b"),
		),
		or(
			val("c"),
			val("d"),
		),
	),

	or(
		and(
			val("a"),
			val("b"),
		),
		and(
			val("c"),
			val("d"),
		),
	),

	and(
		val("a"),
		val("b"),
		or(
			val("c"),
			val("d"),
			val("e"),
		),
	),

	and(
		not(
			val("a"),
		),
		not(
			val("b"),
		),
		or(
			val("c"),
			val("d"),
			val("e"),
		),
	),

	and(
		val("a"),
		val("b"),
		not(
			or(
				val("c"),
				val("d"),
				val("e"),
			),
		),
	),

	and(
		not(
			val("a"),
		),
		not(
			val("b"),
		),
		not(
			or(
				val("c"),
				val("d"),
				val("e"),
			),
		),
	),

	// END
}

func TestEvalosticAgainstStringContains(t *testing.T) {
	fmt.Printf("%100s | %25s | %20s | %20s\n", "Condition", "Test String", "String Contains", "Evalostic Matcher")
	for _, cond := range TheTestConditions {
		// collect all strings of the condition that will be used to check all combinations against the strings
		strMap := make(map[string]struct{})
		cond.GetStrings(strMap)
		var allStrings []string
		for str := range strMap {
			allStrings = append(allStrings, str, strings.ToUpper(str))
		}
		if len(allStrings) > 20 {
			t.Fatalf("too many string exported, this would lead to memory issues")
		}
		sort.Strings(allStrings)
		// build the condition string that will be compared to the string contains method
		condition := cond.ToCondition()
		ev, err := New([]string{condition})
		if err != nil {
			t.Fatalf("parse %s: %s", condition, err)
		}
		// prepare the two matchers
		matcher1 := func(str string) bool {
			return len(ev.Match(str)) > 0
		}
		matcher2 := cond.Match
		compareStrings := []string{""}
		for _, str := range allStrings {
			deepCopy := make([]string, len(compareStrings), len(compareStrings)*2)
			copy(deepCopy, compareStrings)
			for _, compareString := range compareStrings {
				deepCopy = append(deepCopy, compareString+str)
			}
			compareStrings = deepCopy
		}
		for _, str := range compareStrings {
			matches1 := matcher1(str)
			matches2 := matcher2(str)
			fmt.Printf("%100s | %25s | %20s | %20s\n", condition, strconv.Quote(str), strconv.FormatBool(matches1), strconv.FormatBool(matches2))
			if matches1 != matches2 {
				t.Errorf("wrong behaviour for condition %s string %q string contains %t evalostic %t", condition, str, matches1, matches2)
			}
		}
	}
}

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
