package evalostic

type decisionTreeEntry struct {
	value int
	ci    bool
}

type decisionTreeNode struct {
	children    map[decisionTreeEntry]*decisionTreeNode
	notChildren map[decisionTreeEntry]*decisionTreeNode
	outputs     []int
}

func (n *decisionTreeNode) add(path andPathIndex, output int) {
	if len(path) == 0 {
		n.outputs = append(n.outputs, output)
		return
	}
	entry := decisionTreeEntry{value: path[0].i, ci: path[0].ci}
	if path[0].not {
		if child, ok := n.notChildren[entry]; ok {
			child.add(path[1:], output)
		} else {
			child := &decisionTreeNode{
				children:    make(map[decisionTreeEntry]*decisionTreeNode),
				notChildren: make(map[decisionTreeEntry]*decisionTreeNode),
			}
			child.add(path[1:], output)
			n.notChildren[entry] = child
		}
	} else {
		if child, ok := n.children[entry]; ok {
			child.add(path[1:], output)
		} else {
			child := &decisionTreeNode{
				children:    make(map[decisionTreeEntry]*decisionTreeNode),
				notChildren: make(map[decisionTreeEntry]*decisionTreeNode),
			}
			child.add(path[1:], output)
			n.children[entry] = child
		}
	}
}

func (n *decisionTreeNode) find(searches map[decisionTreeEntry]struct{}) (res []int) {
	res = append(res, n.outputs...)
	for search := range searches {
		if child, ok := n.children[search]; ok {
			res = append(res, child.find(searches)...)
		}
	}
	for notSearch, notChild := range n.notChildren {
		if _, ok := searches[notSearch]; !ok {
			res = append(res, notChild.find(searches)...)
		}
	}
	return
}
