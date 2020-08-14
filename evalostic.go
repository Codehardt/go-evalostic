package evalostic

import (
	"fmt"
	"sort"
	"strings"
)

type Evalostic struct {
	conditions  []Node
	ahoCorasick *ahoCorasick
	strings     map[string]int
	mapping     map[int][]int // which string can be found in which condition
}

func New(conditions []string) (*Evalostic, error) {
	e := Evalostic{
		conditions: make([]Node, len(conditions)),
		strings:    make(map[string]int),
		mapping:    make(map[int][]int),
	}
	var stringCounter int
	var allStrings []string
	for i, condition := range conditions {
		root, err := parseCondition(condition)
		if err != nil {
			return nil, fmt.Errorf("condition %d: %s", i, err)
		}
		e.conditions[i] = root
		condStrings := extractStrings(root)
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
	e.ahoCorasick = newAhoCorasick(allStrings)
	return &e, nil
}

func (e *Evalostic) Match(s string) (matchingConditions []int) {
	stringIndices := e.ahoCorasick.match(s)
	possibleConditions := make(map[int]struct{})
	for _, strI := range stringIndices {
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

func conditionMatches(n Node, s string) bool {
	switch v := n.(type) {
	case AND:
		return conditionMatches(v.node1, s) && conditionMatches(v.node2, s)
	case OR:
		return conditionMatches(v.node1, s) || conditionMatches(v.node2, s)
	case NOT:
		return !conditionMatches(v.node, s)
	case VAL:
		return strings.Contains(s, v.nodeValue)
	default:
		return false
	}
}
