package evalostic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportElasticSearchQueryMap(t *testing.T) {
	t.Parallel()
	matchPhrase := func(s string) map[string]interface{} {
		return map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"raw": s,
			},
		}
	}
	wildcard := func(s string) map[string]interface{} {
		return map[string]interface{}{
			"wildcard": map[string]interface{}{
				"raw": map[string]interface{}{
					"case_insensitive": true,
					"value":            "*" + s + "*",
				},
			},
		}
	}
	and := func(must ...map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		}
	}
	or := func(should ...map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
			},
		}
	}
	not := func(m map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []map[string]interface{}{m},
			},
		}
	}
	testCases := []struct {
		name           string
		useMatchPhrase bool
		conditions     []string
		expectedResult map[string]interface{}
	}{
		{
			name:           "empty",
			conditions:     nil,
			expectedResult: nil,
		},
		{
			name:           "simple",
			conditions:     []string{`"a"`},
			expectedResult: wildcard("a"),
		},
		{
			name:           "simple match phrase",
			useMatchPhrase: true,
			conditions:     []string{`"a"`},
			expectedResult: matchPhrase("a"),
		},
		{
			name:           "and",
			conditions:     []string{`"a" AND "b"`},
			expectedResult: and(wildcard("a"), wildcard("b")),
		},
		{
			name: "flatten and",
			conditions: []string{
				`"a" AND ("b" AND "c")`,
			},
			expectedResult: and(wildcard("a"), wildcard("b"), wildcard("c")),
		},
		{
			name: "flatten and 2",
			conditions: []string{
				`("a" AND "b") AND "c"`,
			},
			expectedResult: and(wildcard("a"), wildcard("b"), wildcard("c")),
		},
		{
			name: "flatten and 3",
			conditions: []string{
				`("a" AND "b") AND ("c" AND "d")`,
			},
			expectedResult: and(wildcard("a"), wildcard("b"), wildcard("c"), wildcard("d")),
		},
		{
			name: "flatten and 4",
			conditions: []string{
				`"a" AND ("b" AND ("c" AND "d"))`,
			},
			expectedResult: and(wildcard("a"), wildcard("b"), wildcard("c"), wildcard("d")),
		},
		{
			name:           "or",
			conditions:     []string{`"a" OR "b"`},
			expectedResult: or(wildcard("a"), wildcard("b")),
		},
		{
			name: "flatten or",
			conditions: []string{
				`"a" OR ("b" OR "c")`,
			},
			expectedResult: or(wildcard("a"), wildcard("b"), wildcard("c")),
		},
		{
			name: "flatten or 2",
			conditions: []string{
				`("a" OR "b") OR "c"`,
			},
			expectedResult: or(wildcard("a"), wildcard("b"), wildcard("c")),
		},
		{
			name: "flatten or 3",
			conditions: []string{
				`("a" OR "b") OR ("c" OR "d")`,
			},
			expectedResult: or(wildcard("a"), wildcard("b"), wildcard("c"), wildcard("d")),
		},
		{
			name: "flatten or 4",
			conditions: []string{
				`"a" OR ("b" OR ("c" OR "d"))`,
			},
			expectedResult: or(wildcard("a"), wildcard("b"), wildcard("c"), wildcard("d")),
		},
		{
			name:           "not",
			conditions:     []string{`NOT "a"`},
			expectedResult: not(wildcard("a")),
		},
		{
			name: "double not",
			conditions: []string{
				`NOT (NOT "a")`,
			},
			expectedResult: wildcard("a"),
		},
		{
			name: "tripple not",
			conditions: []string{
				`NOT (NOT (NOT "a"))`,
			},
			expectedResult: not(wildcard("a")),
		},
		{
			name: "quadruple not",
			conditions: []string{
				`NOT (NOT (NOT (NOT "a")))`,
			},
			expectedResult: wildcard("a"),
		},
		{
			name: "multiple",
			conditions: []string{
				`"a" AND "b"`,
				`"c" OR "d"`,
				`NOT "e"`,
			},
			expectedResult: or(
				and(wildcard("a"), wildcard("b")),
				wildcard("c"), wildcard("d"), // OR has been flattened
				not(wildcard("e")),
			),
		},
		{
			name: "complex",
			conditions: []string{
				`"a" AND NOT ("b" OR ("c" OR (NOT "d" AND "e")))`,
				`"f" OR NOT ("g" AND NOT "h" AND NOT "i")`,
			},
			expectedResult: or(
				and(
					wildcard("a"),
					not(
						or(
							wildcard("b"),
							wildcard("c"),
							and(
								not(wildcard("d")),
								wildcard("e"),
							),
						),
					),
				),
				// OR has been flattened
				wildcard("f"),
				not(
					and(
						wildcard("g"),
						not(wildcard("h")),
						not(wildcard("i")),
					),
				),
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := New(tc.conditions)
			require.NoError(t, err)
			result := e.ExportElasticSearchQueryMap("raw", tc.useMatchPhrase)
			assert.EqualValues(t, tc.expectedResult, result)
		})
	}
}
