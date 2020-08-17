package evalostic

import "fmt"

func Example_extractStrings() {
	es := func(cond string) {
		root, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		str, positive := extractStrings(root)
		fmt.Printf("----- %s -----\n", cond)
		fmt.Printf("strings: %+v\n", str)
		fmt.Printf("positive: %t\n", positive)
	}
	es(`"foo"`)
	es(`"foo" AND "bar"`)
	es(`"foo" OR "bar"`)
	es(`NOT "foo"`)
	es(`"foo" AND NOT "bar"`)
	es(`"foo" OR NOT "bar"`)
	es(`NOT ("foo" OR "bar")`)
	es(`NOT ("foo" AND "bar")`)
	es(`NOT ("foo" OR NOT "bar")`)
	es(`NOT ("foo" AND NOT "bar")`)
	// Output:
	// ----- "foo" -----
	// strings: [foo]
	// positive: true
	// ----- "foo" AND "bar" -----
	// strings: [foo bar]
	// positive: true
	// ----- "foo" OR "bar" -----
	// strings: [foo bar]
	// positive: true
	// ----- NOT "foo" -----
	// strings: [foo]
	// positive: false
	// ----- "foo" AND NOT "bar" -----
	// strings: [foo bar]
	// positive: true
	// ----- "foo" OR NOT "bar" -----
	// strings: [foo bar]
	// positive: false
	// ----- NOT ("foo" OR "bar") -----
	// strings: [foo bar]
	// positive: false
	// ----- NOT ("foo" AND "bar") -----
	// strings: [foo bar]
	// positive: false
	// ----- NOT ("foo" OR NOT "bar") -----
	// strings: [foo bar]
	// positive: true
	// ----- NOT ("foo" AND NOT "bar") -----
	// strings: [foo bar]
	// positive: false
}
