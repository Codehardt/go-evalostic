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
	//
}

func ExampleMatchStrings() {
	ms := func(cond string) {
		n, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		n = n.NormalForm()
		matchStrings := MatchStrings(n)
		for _, matchPath := range matchStrings {
			fmt.Println(matchPath.String())
		}
	}
	ms(`("a" AND NOT "b") AND NOT (NOT "c" OR "d") AND ("f" OR "g")`)
	// Output:
	//
}
