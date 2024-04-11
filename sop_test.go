package evalostic

import (
	"fmt"
	"strconv"
	"strings"
)

var sop = func(cond string) {
	n, err := parseCondition(cond)
	if err != nil {
		panic(err)
	}
	fmt.Printf("----- %s -----\n", cond)
	fmt.Println("before:", n.Condition())
	n = n.SOP()
	fmt.Println("after:", n.Condition())
}

func dnf(cond string) {
	n, err := parseCondition(cond)
	if err != nil {
		panic(err)
	}
	fmt.Println("before:", cond)
	var res []string
	for _, mp := range getAndPaths(n.SOP()) {
		var part []string
		for _, str := range mp {
			val := strconv.Quote(str.str)
			if str.not {
				val = "NOT " + val
			}
			part = append(part, val)
		}
		res = append(res, strings.Join(part, " AND "))
	}
	if len(res) == 0 {
		fmt.Println("after: -")
	} else if len(res) == 1 {
		fmt.Println("after:", res[0])
	} else {
		for i, s := range res {
			if strings.Contains(s, " AND ") {
				res[i] = "(" + s + ")"
			}
		}
		fmt.Println("after:", strings.Join(res, " OR "))
	}
}

func ExampleDNF() {
	dnf(`"a"`)
	dnf(`NOT "a"`)
	dnf(`"a" AND "b"`)
	dnf(`"a" OR "b"`)
	dnf(`"a" AND ("b" OR "c")`)
	dnf(`"a" OR ("b" AND "c")`)
	dnf(`"a" AND NOT ("b" OR "c")`)
	dnf(`"a" OR NOT ("b" AND "c")`)
	dnf(`"a" AND ("b" OR NOT "c")`)
	dnf(`"a" OR ("b" AND NOT "c")`)
	dnf(`"a" AND NOT ("b" OR NOT "c")`)
	dnf(`"a" OR NOT ("b" AND NOT "c")`)
	dnf(`"a" OR ("b" OR ("c" OR "d"))`)
	dnf(`("a" OR "b") OR ("c" OR "d")`)
	// Output:
	// before: "a"
	// after: "a"
	// before: NOT "a"
	// after: NOT "a"
	// before: "a" AND "b"
	// after: "a" AND "b"
	// before: "a" OR "b"
	// after: "a" OR "b"
	// before: "a" AND ("b" OR "c")
	// after: ("a" AND "b") OR ("a" AND "c")
	// before: "a" OR ("b" AND "c")
	// after: "a" OR ("b" AND "c")
	// before: "a" AND NOT ("b" OR "c")
	// after: "a" AND NOT "b" AND NOT "c"
	// before: "a" OR NOT ("b" AND "c")
	// after: "a" OR NOT "b" OR NOT "c"
	// before: "a" AND ("b" OR NOT "c")
	// after: ("a" AND "b") OR ("a" AND NOT "c")
	// before: "a" OR ("b" AND NOT "c")
	// after: "a" OR ("b" AND NOT "c")
	// before: "a" AND NOT ("b" OR NOT "c")
	// after: "a" AND "c" AND NOT "b"
	// before: "a" OR NOT ("b" AND NOT "c")
	// after: "a" OR NOT "b" OR "c"
	// before: "a" OR ("b" OR ("c" OR "d"))
	// after: "a" OR "b" OR "c" OR "d"
	// before: ("a" OR "b") OR ("c" OR "d")
	// after: "a" OR "b" OR "c" OR "d"
}

func ExampleSOP() {
	sop(`"a"`)
	sop(`NOT "a"`)
	sop(`"a" AND "b"`)
	sop(`"a" OR "b"`)
	sop(`"a" AND ("b" OR "c")`)
	sop(`"a" AND ("b" OR "c") AND ("d" OR "e")`)
	sop(`"a" OR ("b" AND "c")`)
	sop(`"a" AND NOT ("b" OR "c")`)
	sop(`"a" OR NOT ("b" AND "c")`)
	sop(`"a" AND ("b" OR NOT "c")`)
	sop(`"a" OR ("b" AND NOT "c")`)
	sop(`"a" AND NOT ("b" OR NOT "c")`)
	sop(`"a" OR NOT ("b" AND NOT "c")`)
	sop(`"a" OR ("b" OR ("c" OR "d"))`)
	sop(`("a" OR "b") OR ("c" OR "d")`)
	// Output:
	// ----- "a" -----
	// before: "a"
	// after: "a"
	// ----- NOT "a" -----
	// before: NOT "a"
	// after: NOT "a"
	// ----- "a" AND "b" -----
	// before: ("a" AND "b")
	// after: ("a" AND "b")
	// ----- "a" OR "b" -----
	// before: ("a" OR "b")
	// after: ("a" OR "b")
	// ----- "a" AND ("b" OR "c") -----
	// before: ("a" AND ("b" OR "c"))
	// after: (("a" AND "b") OR ("a" AND "c"))
	// ----- "a" AND ("b" OR "c") AND ("d" OR "e") -----
	// before: (("a" AND ("b" OR "c")) AND ("d" OR "e"))
	//after: (((("a" AND "b") OR ("a" AND "c")) AND "d") OR ((("a" AND "b") OR ("a" AND "c")) AND "e"))
	// ----- "a" OR ("b" AND "c") -----
	// before: ("a" OR ("b" AND "c"))
	// after: ("a" OR ("b" AND "c"))
	// ----- "a" AND NOT ("b" OR "c") -----
	// before: ("a" AND NOT ("b" OR "c"))
	// after: ("a" AND (NOT "b" AND NOT "c"))
	// ----- "a" OR NOT ("b" AND "c") -----
	// before: ("a" OR NOT ("b" AND "c"))
	// after: ("a" OR (NOT "b" OR NOT "c"))
	// ----- "a" AND ("b" OR NOT "c") -----
	// before: ("a" AND ("b" OR NOT "c"))
	// after: (("a" AND "b") OR ("a" AND NOT "c"))
	// ----- "a" OR ("b" AND NOT "c") -----
	// before: ("a" OR ("b" AND NOT "c"))
	// after: ("a" OR ("b" AND NOT "c"))
	// ----- "a" AND NOT ("b" OR NOT "c") -----
	// before: ("a" AND NOT ("b" OR NOT "c"))
	// after: ("a" AND (NOT "b" AND "c"))
	// ----- "a" OR NOT ("b" AND NOT "c") -----
	// before: ("a" OR NOT ("b" AND NOT "c"))
	// after: ("a" OR (NOT "b" OR "c"))
	// ----- "a" OR ("b" OR ("c" OR "d")) -----
	// before: ("a" OR ("b" OR ("c" OR "d")))
	// after: ("a" OR ("b" OR ("c" OR "d")))
	// ----- ("a" OR "b") OR ("c" OR "d") -----
	// before: (("a" OR "b") OR ("c" OR "d"))
	// after: (("a" OR "b") OR ("c" OR "d"))
}

func ExampleSOP_2() {
	sop(`("a" OR "b" OR "c") AND NOT ("d" OR "e" OR "f") AND ("g" OR "h" OR "i")`)
	sop(`NOT (("a" OR "b" OR "c") AND NOT ("d" OR "e" OR "f") AND ("g" OR "h" OR "i"))`)
	// Output:
	// ----- ("a" OR "b" OR "c") AND NOT ("d" OR "e" OR "f") AND ("g" OR "h" OR "i") -----
	// before: (((("a" OR "b") OR "c") AND NOT (("d" OR "e") OR "f")) AND (("g" OR "h") OR "i"))
	// after: (((((("a" AND ((NOT "d" AND NOT "e") AND NOT "f")) OR ("b" AND ((NOT "d" AND NOT "e") AND NOT "f"))) OR ("c" AND ((NOT "d" AND NOT "e") AND NOT "f"))) AND "g") OR (((("a" AND ((NOT "d" AND NOT "e") AND NOT "f")) OR ("b" AND ((NOT "d" AND NOT "e") AND NOT "f"))) OR ("c" AND ((NOT "d" AND NOT "e") AND NOT "f"))) AND "h")) OR (((("a" AND ((NOT "d" AND NOT "e") AND NOT "f")) OR ("b" AND ((NOT "d" AND NOT "e") AND NOT "f"))) OR ("c" AND ((NOT "d" AND NOT "e") AND NOT "f"))) AND "i"))
	// ----- NOT (("a" OR "b" OR "c") AND NOT ("d" OR "e" OR "f") AND ("g" OR "h" OR "i")) -----
	// before: NOT (((("a" OR "b") OR "c") AND NOT (("d" OR "e") OR "f")) AND (("g" OR "h") OR "i"))
	// after: ((((NOT "a" AND NOT "b") AND NOT "c") OR (("d" OR "e") OR "f")) OR ((NOT "g" AND NOT "h") AND NOT "i"))
}

func ExampleMatchStrings() {
	ms := func(cond string) {
		n, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		n = n.SOP()
		matchStrings := getAndPaths(n)
		for _, matchPath := range matchStrings {
			fmt.Println(matchPath.String())
		}
	}
	ms(`("a" OR "b") AND ("c" OR "d") AND ("e" OR "f")`)
	// Output:
	// ----- ("a" OR "b") AND ("c" OR "d") AND ("e" OR "f") -----
	// "a", "c", "e"
	// "a", "d", "e"
	// "b", "c", "e"
	// "b", "d", "e"
	// "a", "c", "f"
	// "a", "d", "f"
	// "b", "c", "f"
	// "b", "d", "f"
}
