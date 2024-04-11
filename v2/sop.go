package evalostic

import (
	"fmt"
	"sort"
	"strings"
)

type andString struct {
	not bool
	str string
}

type andStringIndex struct {
	not bool
	i   int
}

type andPath []andString

type andPathIndex []andStringIndex

func (m andPath) String() string {
	allStrings := make([]string, len(m))
	for i, str := range m {
		allStrings[i] = str.String()
	}
	return strings.Join(allStrings, ", ")
}

func (m andString) String() string {
	var prefix string
	if m.not {
		prefix = "NOT "
	}
	return fmt.Sprintf("%s%q", prefix, m.str)
}

func getAndPaths(n node) []andPath {
	res := getUnsortedAndPaths(n)
	for _, path := range res {
		sort.Slice(path, func(i, j int) bool {
			s1, s2 := path[i], path[j]
			if s1.not && !s2.not {
				return false
			}
			if !s1.not && s2.not {
				return true
			}
			return strings.Compare(s1.str, s2.str) < 0
		})
	}
	return res
}

func getUnsortedAndPaths(n node) []andPath {
	switch v := n.(type) {
	case nodeAND:
		var res []andPath
		c1 := getUnsortedAndPaths(v.node1)
		c2 := getUnsortedAndPaths(v.node2)
		for _, andPathC1 := range c1 {
			for _, andPathC2 := range c2 {
				res = append(res, append(andPathC1, andPathC2...))
			}
		}
		return res
	case nodeOR:
		c1 := getUnsortedAndPaths(v.node1)
		c2 := getUnsortedAndPaths(v.node2)
		return append(c1, c2...)
	case nodeNOT:
		val := v.node.(nodeVAL)
		return []andPath{
			{andString{
				not: true,
				str: val.nodeValue,
			}},
		}
	case nodeVAL:
		return []andPath{
			{andString{
				not: false,
				str: v.nodeValue,
			}},
		}
	default:
		panic("unknown node type")
	}
}

func (n nodeAND) SOP() node {
	//n.node1 = n.node1.SOP()
	if or, ok := n.node1.(nodeOR); ok {
		return (nodeOR{
			twoSubNodes{
				node1: nodeAND{
					twoSubNodes{
						node1: or.node1,
						node2: n.node2,
					},
				},
				node2: nodeAND{
					twoSubNodes{
						node1: or.node2,
						node2: n.node2,
					},
				},
			},
		}).SOP()
	}
	//n.node2 = n.node2.SOP()
	if or, ok := n.node2.(nodeOR); ok {
		return (nodeOR{
			twoSubNodes{
				node1: nodeAND{
					twoSubNodes{
						node1: n.node1,
						node2: or.node1,
					},
				},
				node2: nodeAND{
					twoSubNodes{
						node1: n.node1,
						node2: or.node2,
					},
				},
			},
		}).SOP()
	}
	n.node1 = n.node1.SOP()
	n.node2 = n.node2.SOP()
	return n
}

func (n nodeOR) SOP() node {
	n.node1 = n.node1.SOP()
	n.node2 = n.node2.SOP()
	return n
}

func (n nodeNOT) SOP() node {
	//n.node = n.node.SOP()
	switch v := n.node.(type) {
	case nodeAND:
		return (nodeOR{
			twoSubNodes{
				node1: nodeNOT{
					oneSubNode{
						node: v.node1,
					},
				},
				node2: nodeNOT{
					oneSubNode{
						node: v.node2,
					},
				},
			},
		}).SOP()
	case nodeOR:
		return (nodeAND{
			twoSubNodes{
				node1: nodeNOT{
					oneSubNode{
						node: v.node1,
					},
				},
				node2: nodeNOT{
					oneSubNode{
						node: v.node2,
					},
				},
			},
		}).SOP()
	case nodeVAL:
		return n
	case nodeNOT:
		return v.node.SOP()
	default:
		panic("unknown node type")
	}
}

func (n nodeVAL) SOP() node {
	return n
}
