package evalostic

func extractStrings(n node) ([]string, bool) {
	switch v := n.(type) {
	case nodeVAL:
		return []string{v.nodeValue}, true
	case nodeAND:
		n1str, n1b := extractStrings(v.node1)
		n2str, n2b := extractStrings(v.node2)
		return append(n1str, n2str...), n1b || n2b
	case nodeOR:
		n1str, n1b := extractStrings(v.node1)
		n2str, n2b := extractStrings(v.node2)
		return append(n1str, n2str...), n1b && n2b
	case nodeNOT:
		str, nb := extractStrings(v.node)
		return str, !nb
	default:
		return nil, true
	}
}
