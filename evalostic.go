package evalostic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Codehardt/go-ahocorasick"
)

// Evalostic is a matcher that can apply multiple conditions on a string with some performance optimizations.
// The biggest optimization is that only conditions that contain at least one keyword of the string will be checked,
// these strings will be filtered with the Aho-Corasick algorithm. The only exception are negative conditions (see comment of: Negatives() function).
type Evalostic struct {
	conditions  []node
	ahoCorasick ahocorasick.AhoCorasick
	strings     map[string]int
	mapping     map[int][]int // which string can be found in which condition
	negatives   []int         // conditions that contain strings in a NOT path
}

// New builds a new Evalostic matcher that compiles all conditions to one big rule set that can be applied to strings.
func New(conditions []string) (*Evalostic, error) {
	e := Evalostic{
		conditions: make([]node, len(conditions)),
		strings:    make(map[string]int),
		mapping:    make(map[int][]int),
	}
	var stringCounter int
	var allStrings []string
	for i, condition := range conditions {
		if condition == "" {
			continue // allow empty conditions but ignore them
		}
		root, err := parseCondition(condition)
		if err != nil {
			return nil, fmt.Errorf("condition %d: %s", i, err)
		}
		e.conditions[i] = root
		condStrings, positive := extractStrings(root)
		if !positive {
			e.negatives = append(e.negatives, i)
		}
		for _, str := range condStrings {
			strI, ok := e.strings[str]
			if !ok {
				strI = stringCounter
				stringCounter++
				e.strings[str] = strI
				allStrings = append(allStrings, str)
			}
			e.mapping[strI] = append(e.mapping[strI], i)
		}
	}
	if len(allStrings) > 0 {
		e.ahoCorasick = ahocorasick.New(allStrings)
	}
	return &e, nil
}

// Match returns all indices of conditions that match the provided string
func (e *Evalostic) Match(s string) (matchingConditions []int) {
	var (
		stringIndices                []int
		stringIndicesCaseInsensitive []int
	)
	if e.ahoCorasick != nil {
		stringIndices = e.ahoCorasick.Match(s)
		stringIndicesCaseInsensitive = e.ahoCorasick.Match(strings.ToLower(s))
	}
	possibleConditions := make(map[int]struct{})
	for _, cond := range e.negatives {
		possibleConditions[cond] = struct{}{}
	}
	for _, strI := range append(stringIndices, stringIndicesCaseInsensitive...) {
		for _, conditionI := range e.mapping[strI] {
			possibleConditions[conditionI] = struct{}{}
		}
	}
	for possibleCondition := range possibleConditions {
		if conditionMatches(e.conditions[possibleCondition], s) {
			matchingConditions = append(matchingConditions, possibleCondition)
		}
	}
	sort.Ints(matchingConditions)
	return
}

// Negatives returns all conditions with a negative path. Negative paths are paths that were negated by a NOT
// and weren't combined with an AND and an positive path.
// Examples for negative paths:
// - NOT "foo"
// - "foo" OR NOT "bar"
// Examples for positive paths:
// - "foo"
// - "foo" AND NOT "bar"
// The big disadvantage of negative paths are that all conditions have to be visited, even if no keyword of the condition
// were found by Aho-Corasick in the string. You should try to avoid using many negative paths due to runtime O(n),
// where n = number of negative conditions
func (e *Evalostic) Negatives() []int {
	res := make([]int, len(e.negatives))
	copy(res, e.negatives)
	return res
}
