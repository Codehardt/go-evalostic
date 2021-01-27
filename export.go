package evalostic

import (
	"encoding/json"
)

func (e *Evalostic) ExportElasticSearchQuery(wildcardField string) string {
	indexToStrings := make(map[int]string)
	for k, v := range e.strings {
		indexToStrings[v] = k
	}
	query := e.exportElasticSearchQuery(wildcardField, indexToStrings, decisionTreeEntry{value: -1}, e.decisionTree)
	if query == nil {
		return ""
	}
	b, _ := json.Marshal(query)
	return string(b)
}

func (e *Evalostic) exportElasticSearchQuery(wildcardField string, indexToStrings map[int]string, entry decisionTreeEntry, node *decisionTreeNode) interface{} {
	var must []interface{}
	var mustNot []interface{}
	for entry, node := range node.children {
		subQuery := e.exportElasticSearchQuery(wildcardField, indexToStrings, entry, node)
		if subQuery != nil {
			must = append(must, subQuery)
		}
	}
	for entry, node := range node.notChildren {
		subQuery := e.exportElasticSearchQuery(wildcardField, indexToStrings, entry, node)
		if subQuery != nil {
			mustNot = append(mustNot, subQuery)
		}
	}
	var wildcard interface{}
	if entry.value != -1 {
		wildcard = map[string]interface{}{
			"wildcard": map[string]interface{}{
				wildcardField: map[string]interface{}{
					"value":            indexToStrings[entry.value],
					"case_insensitive": entry.ci,
				},
			},
		}
	}
	if len(node.outputs) != 0 {
		return wildcard
	}
	if len(must) == 0 && len(mustNot) == 0 {
		return nil
	}
	boolRes := make(map[string]interface{})
	if len(must) != 0 {
		boolRes["should"] = must
	}
	if len(mustNot) != 0 {
		if len(mustNot) == 1 {
			boolRes["must_not"] = mustNot
		} else {
			boolRes["must_not"] = map[string]interface{}{
				"bool": map[string]interface{}{
					"should": mustNot,
				},
			}
		}
	}
	if wildcard == nil {
		return map[string]interface{}{"bool": boolRes}
	}
	res := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []interface{}{
				wildcard,
				map[string]interface{}{"bool": boolRes},
			},
		},
	}
	return res
}
