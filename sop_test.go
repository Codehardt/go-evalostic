package evalostic

import "fmt"

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

func ExampleSOP() {
	sop(`"a"`)
	sop(`NOT "a"`)
	sop(`"a" AND "b"`)
	sop(`"a" OR "b"`)
	sop(`"a" AND ("b" OR "c")`)
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
	ms(`("a" AND NOT "b") AND NOT (NOT "c" OR "d") AND ("f" OR "g")`)
	// Output:
	// ----- ("a" AND NOT "b") AND NOT (NOT "c" OR "d") AND ("f" OR "g") -----
	// "a", "c", "f", NOT "b", NOT "d"
	// "a", "c", "g", NOT "b", NOT "d"
}
