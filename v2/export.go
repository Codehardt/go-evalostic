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
	indexToStrings := make(map[int]string)
	for k, v := range e.strings {
		indexToStrings[v] = k
	}
	query := e.exportElasticSearchQuerySub(wildcardField, useMatchPhrase, indexToStrings, decisionTreeEntry{value: -1}, e.decisionTree, false)
	if query == nil {
		return make(map[string]interface{})
	}
	return query
}

var wildcardReplacer = strings.NewReplacer("\\", "\\\\", "*", "\\*", "?", "\\?")

func (e *Evalostic) exportElasticSearchQuerySub(wildcardField string, useMatchPhrase bool, indexToStrings map[int]string, entry decisionTreeEntry, node *decisionTreeNode, not bool) map[string]interface{} {
	isLeaf := len(node.outputs) != 0
	wildcard := map[string]interface{}{
		"wildcard": map[string]interface{}{
			wildcardField: map[string]interface{}{
				"value":            "*" + wildcardReplacer.Replace(indexToStrings[entry.value]) + "*",
				"case_insensitive": true,
			},
		},
	}
	if useMatchPhrase {
		wildcard = map[string]interface{}{
			"match_phrase": map[string]interface{}{
				wildcardField: indexToStrings[entry.value],
			},
		}
	}
	if not {
		wildcard = map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []interface{}{wildcard},
			},
		}
	}
	if entry.value == -1 {
		// special case: do not use root node as wildcard
		wildcard = nil
	}
	if isLeaf && wildcard != nil {
		// special case: if it's a leaf, we don't need to process the sub tree
		return wildcard
	}

	var should []map[string]interface{}

	for subEntry, subNode := range node.children {
		if subQuery := e.exportElasticSearchQuerySub(wildcardField, useMatchPhrase, indexToStrings, subEntry, subNode, false); subQuery != nil {
			should = append(should, subQuery)
		}
	}
	for subEntry, subNode := range node.notChildren {
		if subQuery := e.exportElasticSearchQuerySub(wildcardField, useMatchPhrase, indexToStrings, subEntry, subNode, true); subQuery != nil {
			should = append(should, subQuery)
		}
	}

	toQuery := func(should []map[string]interface{}) map[string]interface{} {
		if len(should) == 0 {
			return nil
		}
		var res map[string]interface{}
		if len(should) == 1 {
			res = should[0]
		} else {
			res = map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			}
		}
		return res
	}

	childQuery := toQuery(should)
	if childQuery == nil {
		return nil
	}
	if wildcard == nil {
		return childQuery
	}
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []interface{}{
				wildcard,
				childQuery,
			},
		},
	}
}
