package evalostic

import "fmt"

func Example_tokenize() {
	tk := func(cond string) {
		tokens, err := tokenize(cond)
		if err != nil {
			panic(err)
		}
		fmt.Printf("----- %s -----\n", cond)
		for _, token := range tokens {
			fmt.Println(token.String())
		}
	}
	tk(`"foo"`)
	tk(`"foo" AND "bar"`)
	tk(`"foo" AND ("bar" OR "baz")`)
	tk(`"foo" AND ("bar" OR NOT "baz")`)
	// Output:
	// ----- "foo" -----
	// nodeVAL ( "foo" ) at pos 1
	// ----- "foo" AND "bar" -----
	// nodeVAL ( "foo" ) at pos 1
	// nodeAND ( AND ) at pos 7
	// nodeVAL ( "bar" ) at pos 11
	// ----- "foo" AND ("bar" OR "baz") -----
	// nodeVAL ( "foo" ) at pos 1
	// nodeAND ( AND ) at pos 7
	// LPAR ( ( ) at pos 11
	// nodeVAL ( "bar" ) at pos 12
	// nodeOR ( OR ) at pos 18
	// nodeVAL ( "baz" ) at pos 21
	// RPAR ( ) ) at pos 26
	// ----- "foo" AND ("bar" OR NOT "baz") -----
	// nodeVAL ( "foo" ) at pos 1
	// nodeAND ( AND ) at pos 7
	// LPAR ( ( ) at pos 11
	// nodeVAL ( "bar" ) at pos 12
	// nodeOR ( OR ) at pos 18
	// nodeNOT ( NOT ) at pos 21
	// nodeVAL ( "baz" ) at pos 25
	// RPAR ( ) ) at pos 30
}
