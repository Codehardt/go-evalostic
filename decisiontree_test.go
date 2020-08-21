package evalostic

import (
	"fmt"
	"sort"
	"strings"
)

func ExampleDecisionTree() {
	ct := new(decisionTreeNode)
	ct.children = make(map[decisionTreeEntry]*decisionTreeNode)
	ct.notChildren = make(map[decisionTreeEntry]*decisionTreeNode)
	stringMap := make(map[string]int)
	for i, cond := range []string{
		`"a"`,
		`"a" AND "b"`,
		`"c" AND ("d" OR "e")`,
		`"e" AND NOT ("d" AND ("f"i OR NOT "g"))`,
	} {
		n, err := parseCondition(cond)
		if err != nil {
			panic(err)
		}
		n = n.NormalForm()
		allStrings, _ := extractStrings(n)
		sort.Strings(allStrings)
		for _, str := range allStrings {
			if _, ok := stringMap[str]; !ok {
				stringMap[str] = len(stringMap)
			}
		}
		for _, mp := range getAndPaths(n) {
			mpi := make([]andStringIndex, len(mp))
			for i, ms := range mp {
				mpi[i] = andStringIndex{ci: ms.ci, not: ms.not, i: stringMap[ms.str]}
			}
			ct.add(mpi, i)
		}
	}
	ct.print(0)
	// Output:
	// outputs:
	//     []
	// children:
	//     +++ 0 +++
	//     outputs:
	//         [0]
	//     children:
	//         +++ 1 +++
	//         outputs:
	//             [1]
	//         children:
	//             (none)
	//         not children:
	//             (none)
	//     not children:
	//         (none)
	//     +++ 2 +++
	//     outputs:
	//         []
	//     children:
	//         +++ 3 +++
	//         outputs:
	//             [2]
	//         children:
	//             (none)
	//         not children:
	//             (none)
	//         +++ 4 +++
	//         outputs:
	//             [2]
	//         children:
	//             (none)
	//         not children:
	//             (none)
	//     not children:
	//         (none)
	//     +++ 4 +++
	//     outputs:
	//         []
	//     children:
	//         +++ 6 +++
	//         outputs:
	//             []
	//         children:
	//             (none)
	//         not children:
	//             +++ NOT 3 +++
	//             outputs:
	//                 [3]
	//             children:
	//                 (none)
	//             not children:
	//                 (none)
	//             +++ NOT 5i +++
	//             outputs:
	//                 [3]
	//             children:
	//                 (none)
	//             not children:
	//                 (none)
	//     not children:
	//         +++ NOT 3 +++
	//         outputs:
	//             []
	//         children:
	//             (none)
	//         not children:
	//             +++ NOT 3 +++
	//             outputs:
	//                 [3]
	//             children:
	//                 (none)
	//             not children:
	//                 (none)
	//             +++ NOT 5i +++
	//             outputs:
	//                 [3]
	//             children:
	//                 (none)
	//             not children:
	//                 (none)
	// not children:
	//     (none)
}

func (n decisionTreeEntry) toStr() string {
	if n.ci {
		return fmt.Sprintf("%di", n.value)
	}
	return fmt.Sprintf("%d", n.value)
}

func (n decisionTreeNode) print(indent int) {
	indentStr := strings.Repeat("    ", indent)
	indentStr2 := strings.Repeat("    ", indent+1)
	type mapEntry struct {
		k decisionTreeEntry
		v *decisionTreeNode
	}
	toMapEntries := func(m map[decisionTreeEntry]*decisionTreeNode) (res []mapEntry) {
		for k, v := range m {
			res = append(res, mapEntry{k, v})
		}
		sort.Slice(res, func(i, j int) bool { return strings.Compare(res[i].k.toStr(), res[j].k.toStr()) < 0 })
		return
	}
	fmt.Printf("%soutputs:\n", indentStr)
	fmt.Printf("%s%+v\n", indentStr2, n.outputs)
	fmt.Printf("%schildren:\n", indentStr)
	for _, child := range toMapEntries(n.children) {
		fmt.Printf("%s+++ %s +++\n", indentStr2, child.k.toStr())
		child.v.print(indent + 1)
	}
	if len(n.children) == 0 {
		fmt.Printf("%s(none)\n", indentStr2)
	}
	fmt.Printf("%snot children:\n", indentStr)
	for _, child := range toMapEntries(n.notChildren) {
		fmt.Printf("%s+++ NOT %s +++\n", indentStr2, child.k.toStr())
		child.v.print(indent + 1)
	}
	if len(n.notChildren) == 0 {
		fmt.Printf("%s(none)\n", indentStr2)
	}
}
