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
	decisionTree *decisionTreeNode
	ahoCorasick  ahocorasick.AhoCorasick
	strings      map[string]int
	mapping      map[int][]int // which string can be found in which condition
}

// New builds a new Evalostic matcher that compiles all conditions to one big rule set that can be applied to strings.
func New(conditions []string) (*Evalostic, error) {
	e := Evalostic{
		decisionTree: new(decisionTreeNode),
		strings:      make(map[string]int),
		mapping:      make(map[int][]int),
	}
	e.decisionTree.children = make(map[decisionTreeEntry]*decisionTreeNode)
	e.decisionTree.notChildren = make(map[decisionTreeEntry]*decisionTreeNode)
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
		condStrings, _ := extractStrings(root)
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
		for _, mp := range getAndPaths(root.NormalForm()) {
			mpi := make(andPathIndex, len(mp))
			for i, ms := range mp {
				mpi[i] = andStringIndex{ci: ms.ci, not: ms.not, i: e.strings[ms.str]}
			}
			e.decisionTree.add(mpi, i)
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
	decisionTreeEntries := make(map[decisionTreeEntry]struct{})
	for _, si := range stringIndices {
		decisionTreeEntries[decisionTreeEntry{value: si}] = struct{}{}
	}
	for _, si := range stringIndicesCaseInsensitive {
		decisionTreeEntries[decisionTreeEntry{value: si, ci: true}] = struct{}{}
	}
	unique := make(map[int]struct{})
	for _, matchingCondition := range e.decisionTree.find(decisionTreeEntries) {
		unique[matchingCondition] = struct{}{}
	}
	for matchingCondition := range unique {
		matchingConditions = append(matchingConditions, matchingCondition)
	}
	sort.Ints(matchingConditions)
	return
}
