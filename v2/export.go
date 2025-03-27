package evalostic

import (
	"encoding/json"
	"strings"
)

// ExportElasticSearchQuery exports the compiled query into an ElasticSearch query, e.g.
// `"foo" OR "baz"` will be compiled to
// {"bool":{"should":[{"wildcard":{"raw":{"case_insensitive":false,"value":"*foo*"}}},{"wildcard":{"raw":{"case_insensitive":false,"value":"*bar*"}}}]}}
func (e *Evalostic) ExportElasticSearchQuery(wildcardField string, useMatchPhrase bool) string {
	b, _ := json.MarshalIndent(e.ExportElasticSearchQueryMap(wildcardField, useMatchPhrase), "", "  ")
	return string(b)
}

// ExportElasticSearchQuery exports the compiled query into an ElasticSearch query, e.g.
// `"foo" OR "baz"` will be compiled to
// {"bool":{"should":[{"wildcard":{"raw":{"case_insensitive":false,"value":"foo"}}},{"wildcard":{"raw":{"case_insensitive":false,"value":"bar"}}}]}}
func (e *Evalostic) ExportElasticSearchQueryMap(wildcardField string, useMatchPhrase bool) map[string]interface{} {
	var root node
	for _, n := range e.orig {
		if root == nil {
			root = n
		} else {
			root = nodeOR{twoSubNodes{root, n}}
		}
	}
	return nodeToElasticSearchQuery(root, useMatchPhrase)
}

func nodeToElasticSearchQuery(n node, useMatchPhrase bool) map[string]interface{} {
	switch v := n.(type) {
	case nodeVAL:
		return leafToElasticSearchQuery(v, useMatchPhrase)
	case nodeNOT:
		return notToElasticSearchQuery(v, useMatchPhrase)
	case nodeOR:
		return orToElasticSearchQuery(v, useMatchPhrase)
	case nodeAND:
		return andToElasticSearchQuery(v, useMatchPhrase)
	default:
		return nil
	}
}

var wildcardReplacer = strings.NewReplacer("\\", "\\\\", "*", "\\*", "?", "\\?")

func notToElasticSearchQuery(n nodeNOT, useMatchPhrase bool) map[string]interface{} {
	if not, ok := n.node.(nodeNOT); ok { // check for double negation
		return nodeToElasticSearchQuery(not.node, useMatchPhrase)
	}
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must_not": []map[string]interface{}{
				nodeToElasticSearchQuery(n.node, useMatchPhrase),
			},
		},
	}
}

func flattenOr(n nodeOR) []node {
	var nodes []node
	if or, ok := n.node1.(nodeOR); ok {
		nodes = append(nodes, flattenOr(or)...)
	} else {
		nodes = append(nodes, n.node1)
	}
	if or, ok := n.node2.(nodeOR); ok {
		nodes = append(nodes, flattenOr(or)...)
	} else {
		nodes = append(nodes, n.node2)
	}
	return nodes
}

func flattenAnd(n nodeAND) []node {
	var nodes []node
	if and, ok := n.node1.(nodeAND); ok {
		nodes = append(nodes, flattenAnd(and)...)
	} else {
		nodes = append(nodes, n.node1)
	}
	if and, ok := n.node2.(nodeAND); ok {
		nodes = append(nodes, flattenAnd(and)...)
	} else {
		nodes = append(nodes, n.node2)
	}
	return nodes
}

func orToElasticSearchQuery(n nodeOR, useMatchPhrase bool) map[string]interface{} {
	var should []map[string]interface{}
	for _, node := range flattenOr(n) {
		should = append(should, nodeToElasticSearchQuery(node, useMatchPhrase))
	}
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"should": should,
		},
	}
}

func andToElasticSearchQuery(n nodeAND, useMatchPhrase bool) map[string]interface{} {
	var must []map[string]interface{}
	for _, node := range flattenAnd(n) {
		must = append(must, nodeToElasticSearchQuery(node, useMatchPhrase))
	}
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must": must,
		},
	}
}

func leafToElasticSearchQuery(n nodeVAL, useMatchPhrase bool) map[string]interface{} {
	if useMatchPhrase {
		return map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"raw": n.nodeValue,
			},
		}
	}
	return map[string]interface{}{
		"wildcard": map[string]interface{}{
			"raw": map[string]interface{}{
				"value":            "*" + wildcardReplacer.Replace(n.nodeValue) + "*",
				"case_insensitive": true,
			},
		},
	}
}
