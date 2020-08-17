package ahocorasick

// AhoCorasick is an interface that returns all matching strings in a text. Use New() to initialize a new AhoCorasick interface.
type AhoCorasick interface {
	// Match returns all indices of strings that were found in the passed text
	Match(text string) []int
}

type ahoCorasick struct {
	root *node
}

// Match is the interface implementation of AhoCorasick's Match function
func (a *ahoCorasick) Match(text string) []int {
	return a.root.find(text)
}

// New builds a new AhoCorasick interface.
func New(allStrings []string) AhoCorasick {
	ac := &ahoCorasick{root: new(node)}
	ac.root.children = make(map[byte]*node)
	ac.root.fail = ac.root
	for i, s := range allStrings {
		ac.root.add(i, s)
	}
	allFailures(ac.root, ac.root, nil)
	return ac
}

type node struct {
	children map[byte]*node
	leaf     *int
	fail     *node
}

func (n *node) add(i int, s string) {
	child := n.children[s[0]]
	if child == nil {
		child = new(node)
		child.children = make(map[byte]*node)
		n.children[s[0]] = child
	}
	s = s[1:]
	if s == "" {
		child.leaf = &i
	} else {
		child.add(i, s)
	}
}

func allFailures(root *node, child *node, prefix []byte) {
	for b, child := range child.children {
		child.fail = failure(root, append(prefix, b)[1:])
		allFailures(root, child, append(prefix, b))
	}
}

func failure(root *node, suffix []byte) *node {
	curr := root
	for _, b := range suffix {
		child, ok := curr.children[b]
		if !ok {
			return failure(root, suffix[1:]) // suffix not found in trie, try a shorter suffix
		}
		curr = child
	}
	return curr
}

func (n *node) find(s string) (res []int) {
	if n.leaf != nil {
		res = append(res, *n.leaf)
	}
	if s == "" {
		return
	}
	if child, ok := n.children[s[0]]; ok {
		res = append(res, child.find(s[1:])...)
	} else if n.fail != nil {
		res = append(res, n.fail.find(s[1:])...)
	}
	return
}
