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
	/*p(`"foo"`)
	p(`"foo" AND "bar"`)
	p(`"foo" AND NOT "bar"`)
	p(`"foo" AND NOT ("bar" OR "baz")`)
	p(`"foo" AND (("bar") AND "baz")`)*/
	p(`"foo" AND ("bar" AND ("baz"))`)
	// Output:
	// ----- "foo" -----
	// nodeVAL{"foo"}
	// ----- "foo" AND "bar" -----
	// nodeAND{nodeVAL{"foo"},nodeVAL{"bar"}}
	// ----- "foo" AND NOT "bar" -----
	// nodeAND{nodeVAL{"foo"},nodeNOT{nodeVAL{"bar"}}}
	// ----- "foo" AND NOT ("bar" OR "baz") -----
	// nodeAND{nodeVAL{"foo"},nodeNOT{nodeOR{nodeVAL{"bar"},nodeVAL{"baz"}}}}
}
