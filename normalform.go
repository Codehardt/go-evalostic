package evalostic

import (
	"fmt"
	"sort"
	"strings"
)

type matchstr struct {
	not bool
	ci  bool
	str string
}

type matchpath []matchstr

func (m matchpath) String() string {
	allStrings := make([]string, len(m))
	for i, str := range m {
		allStrings[i] = str.String()
	}
	return strings.Join(allStrings, ", ")
}

func (m matchstr) String() string {
	var suffix string
	if m.ci {
		suffix = "i"
	}
	var prefix string
	if m.not {
		prefix = "NOT "
	}
	return fmt.Sprintf("%s%q%s", prefix, m.str, suffix)
}

func MatchStrings(n node) []matchpath {
	res := matchStrings(n)
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

func matchStrings(n node) []matchpath {
	switch v := n.(type) {
	case nodeAND:
		return []matchpath{append(matchStrings(v.node1)[0], matchStrings(v.node2)[0]...)}
	case nodeOR:
		return append(matchStrings(v.node1), matchStrings(v.node2)...)
	case nodeNOT:
		val := v.node.(nodeVAL)
		return []matchpath{
			{matchstr{
				not: true,
				ci:  val.caseInsensitive,
				str: val.nodeValue,
			}},
		}
	case nodeVAL:
		return []matchpath{
			{matchstr{
				not: false,
				ci:  v.caseInsensitive,
				str: v.nodeValue,
			}},
		}
	default:
		panic("unknown node type")
	}
}

func (n nodeAND) NormalForm() node {
	n.node1 = n.node1.NormalForm()
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
		}).NormalForm()
	}
	n.node2 = n.node2.NormalForm()
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
		}).NormalForm()
	}
	return n
}

func (n nodeOR) NormalForm() node {
	n.node1 = n.node1.NormalForm()
	n.node2 = n.node2.NormalForm()
	return n
}

func (n nodeNOT) NormalForm() node {
	n.node = n.node.NormalForm()
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
		}).NormalForm()
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
		}).NormalForm()
	case nodeVAL:
		return n
	case nodeNOT:
		return v.node
	default:
		panic("unknown node type")
	}
}

func (n nodeVAL) NormalForm() node {
	return n
}
