package evalostic

import (
	"encoding/json"
)

// ExportElasticSearchQuery exports the compiled query into an ElasticSearch query, e.g.
// `"foo" OR "baz"` will be compiled to
// {"bool":{"should":[{"wildcard":{"raw":{"case_insensitive":false,"value":"foo"}}},{"wildcard":{"raw":{"case_insensitive":false,"value":"bar"}}}]}}
func (e *Evalostic) ExportElasticSearchQuery(wildcardField string) string {
	indexToStrings := make(map[int]string)
	for k, v := range e.strings {
		indexToStrings[v] = k
	}
	query := e.exportElasticSearchQuerySub(wildcardField, indexToStrings, decisionTreeEntry{value: -1}, e.decisionTree)
	if query == nil {
		return ""
	}
	b, _ := json.Marshal(query)
	return string(b)
}

func (e *Evalostic) exportElasticSearchQuerySub(wildcardField string, indexToStrings map[int]string, entry decisionTreeEntry, node *decisionTreeNode) interface{} {

	type Map map[string]interface{}

	isLeaf := len(node.outputs) != 0
	wildcard := Map{
		"wildcard": Map{
			wildcardField: Map{
				"value":            indexToStrings[entry.value],
				"case_insensitive": entry.ci,
			},
		},
	}
	if entry.value == -1 {
		// special case: do not use root node as wildcard
		wildcard = nil
	}
	if isLeaf && wildcard != nil {
		// special case: if it's a leaf, we don't need to process the sub tree
		return wildcard
	}

	var should, shouldNot []interface{}

	for subEntry, subNode := range node.children {
		if subQuery := e.exportElasticSearchQuerySub(wildcardField, indexToStrings, subEntry, subNode); subQuery != nil {
			should = append(should, subQuery)
		}
	}
	for subEntry, subNode := range node.notChildren {
		if subQuery := e.exportElasticSearchQuerySub(wildcardField, indexToStrings, subEntry, subNode); subQuery != nil {
			shouldNot = append(shouldNot, subQuery)
		}
	}

	toQuery := func(should []interface{}, not bool) interface{} {
		if len(should) == 0 {
			return nil
		}
		var res interface{}
		if len(should) == 1 {
			res = should[0]
		} else {
			res = Map{
				"bool": Map{
					"should": should,
				},
			}
		}
		if not {
			// wrap OR conditions with a NOT
			res = Map{
				"bool": Map{
					"must_not": []interface{}{res},
				},
			}
		}
		return res
	}

	notChildQuery := toQuery(shouldNot, true)
	if notChildQuery != nil {
		should = append(should, notChildQuery)
	}
	childQuery := toQuery(should, false)
	if childQuery == nil {
		return nil
	}
	if wildcard == nil {
		return childQuery
	}
	return Map{
		"bool": Map{
			"must": []interface{}{
				wildcard,
				childQuery,
			},
		},
	}
}
