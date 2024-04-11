package evalostic

import "fmt"

func Example_parse() {
	p := func(cond string) {
		t, err := tokenize(cond)
		if err != nil {
			panic(err)
		}
		root, err := parse(t)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		fmt.Println(root.String())
	}
	p(`"foo"`)
	p(`"foo" AND "bar"`)
	p(`"foo" AND NOT "bar"`)
	p(`"foo" AND NOT ("bar" OR "baz")`)
	p(`("foo" OR "bar") AND ("bar" OR "baz") AND ("baaz" OR "qux")`)
	// Output:
	// ----- "foo" -----
	// nodeVAL{"foo"}
	// ----- "foo" AND "bar" -----
	// nodeAND{nodeVAL{"foo"},nodeVAL{"bar"}}
	// ----- "foo" AND NOT "bar" -----
	// nodeAND{nodeVAL{"foo"},nodeNOT{nodeVAL{"bar"}}}
	// ----- "foo" AND NOT ("bar" OR "baz") -----
	// nodeAND{nodeVAL{"foo"},nodeNOT{nodeOR{nodeVAL{"bar"},nodeVAL{"baz"}}}}
	// ----- ("foo" OR "bar") AND ("bar" OR "baz") AND ("baaz" OR "qux") -----
	// nodeAND{nodeAND{nodeOR{nodeVAL{"foo"},nodeVAL{"bar"}},nodeOR{nodeVAL{"bar"},nodeVAL{"baz"}}},nodeOR{nodeVAL{"baaz"},nodeVAL{"qux"}}}
}

func Example_parse_multi() {
	p := func(cond string) {
		t, err := tokenize(cond)
		if err != nil {
			panic(err)
		}
		root, err := parse(t)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		fmt.Println(root.Condition())
	}
	p(`NOT "foo" OR NOT "bar" OR NOT "baz" OR NOT "qux"`)
	p(`"qux" AND "foo" OR "bar" AND "baz"`)
	p(`"qux" OR "foo" AND "bar" OR "baz"`)
	// Output:
	// ----- NOT "foo" OR NOT "bar" OR NOT "baz" OR NOT "qux" -----
	// ((NOT "foo" OR NOT "bar") OR (NOT "baz" OR NOT "qux"))
	// ----- "qux" AND "foo" OR "bar" AND "baz" -----
	// (("qux" AND "foo") OR ("bar" AND "baz"))
	// ----- "qux" OR "foo" AND "bar" OR "baz" -----
	// (("qux" OR ("foo" AND "bar")) OR "baz")
}
