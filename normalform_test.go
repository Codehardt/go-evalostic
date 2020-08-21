package evalostic

import "fmt"

func ExampleNormalForm() {
	nf := func(cond string) {
		n, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		fmt.Println("before:", n.Condition())
		n = n.NormalForm()
		fmt.Println("after:", n.Condition())
	}
	nf(`"a"`)
	nf(`NOT "a"`)
	nf(`"a" AND "b"`)
	nf(`"a" OR "b"`)
	nf(`"a" AND ("b" OR "c")`)
	nf(`"a" OR ("b" AND "c")`)
	nf(`"a" AND NOT ("b" OR "c")`)
	nf(`"a" OR NOT ("b" AND "c")`)
	nf(`"a" AND ("b" OR NOT "c")`)
	nf(`"a" OR ("b" AND NOT "c")`)
	nf(`"a" AND NOT ("b" OR NOT "c")`)
	nf(`"a" OR NOT ("b" AND NOT "c")`)
	nf(`"a" OR ("b" OR ("c" OR "d"))`)
	nf(`("a" OR "b") OR ("c" OR "d")`)
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

func ExampleMatchStrings() {
	ms := func(cond string) {
		n, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		n = n.NormalForm()
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
