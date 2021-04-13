package evalostic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type node interface {
	String() string
	Condition() string
	Value() string
	Children() (node, node)
	SOP() node
}

type (
	oneSubNode  struct{ node node }
	twoSubNodes struct{ node1, node2 node }
	valueNode   struct {
		nodeValue       string
		caseInsensitive bool
	}
)

func (oneSubNode) Value() string  { return "" }
func (twoSubNodes) Value() string { return "" }
func (n valueNode) Value() string { return n.nodeValue }

func (n oneSubNode) Children() (node, node)  { return n.node, nil }
func (n twoSubNodes) Children() (node, node) { return n.node1, n.node2 }
func (valueNode) Children() (node, node)     { return nil, nil }

type (
	nodeAND struct{ twoSubNodes }
	nodeOR  struct{ twoSubNodes }
	nodeNOT struct{ oneSubNode }
	nodeVAL struct{ valueNode }
)

func (n nodeAND) String() string { return fmt.Sprintf("nodeAND{%s,%s}", n.node1, n.node2) }
func (n nodeOR) String() string  { return fmt.Sprintf("nodeOR{%s,%s}", n.node1, n.node2) }
func (n nodeNOT) String() string { return fmt.Sprintf("nodeNOT{%s}", n.node) }
func (n nodeVAL) String() string { return fmt.Sprintf("nodeVAL{%q}", n.nodeValue) }

func parseCondition(s string) (node, error) {
	t, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	return parse(t)
}

func parse(tokens []token) (node, error) {
	res := make([]interface{}, len(tokens))
	for i, token := range tokens {
		res[i] = token
	}
	// identify subexpressions with parentheses
	for {
		lPos := findToken(tokens, tokenTypeLPAR)
		if lPos < 0 {
			break
		}
		if lPos >= len(tokens)-1 {
			return nil, errors.New("condition ends with an opening parentheses")
		}
		var (
			rPos   = -1
			offset = lPos + 1
		)
		for {
			if offset >= len(tokens) {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos = findToken(tokens[offset:], tokenTypeRPAR)
			if rPos < 0 {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos += offset
			if lPos+1 == rPos {
				return nil, errors.New("empty subexpression found")
			}
			anotherLPos := findToken(tokens[offset:rPos], tokenTypeLPAR)
			if anotherLPos < 0 {
				break
			}
			offset = rPos + 1
		}
		subNode, err := parse(tokens[lPos+1 : rPos])
		if err != nil {
			return nil, fmt.Errorf("could not parse subexpression: %s", err)
		}
		res = append(res[:lPos], append([]interface{}{subNode}, res[rPos+1:]...)...)
		tokens = append(tokens[:lPos], append([]token{{tokenType: tokenTypeNONE}}, tokens[rPos+1:]...)...)
	}
	for _, tokenType := range []tokenType{
		tokenTypeVAL,
		tokenTypeNOT,
		tokenTypeAND,
		tokenTypeOR,
	} {
		var tokenFound = true
		for tokenFound {
			tokenFound = false
			for i, elem := range res {
				token, _ := elem.(token)
				if token.tokenType != tokenType {
					continue
				}
				tokenFound = true
				switch tokenType {
				case tokenTypeVAL:
					matched := token.matched
					if token.caseInsensitive {
						matched = strings.ToLower(matched)
					}
					res[i] = nodeVAL{valueNode{nodeValue: matched, caseInsensitive: token.caseInsensitive}}
				default:
					if i+1 >= len(res) {
						return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
					}
					subNode1, ok := res[i+1].(node)
					if !ok {
						return nil, fmt.Errorf("parameter for %s operator is not a node (1), got: %s", tokenTypeString[tokenType], res[i+1])
					}
					res = append(res[:i+1], res[i+2:]...) // remove the (i+1)th element because it has become a sub node
					switch tokenType {
					case tokenTypeNOT:
						res[i] = nodeNOT{oneSubNode{node: subNode1}}
					default:
						if i == 0 {
							return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
						}
						subNode2, ok := res[i-1].(node)
						if !ok {
							return nil, fmt.Errorf("parameter for %s operator is not a node (2)", tokenTypeString[tokenType])
						}
						n := twoSubNodes{subNode2, subNode1}
						switch tokenType {
						case tokenTypeAND:
							res[i] = nodeAND{n}
						case tokenTypeOR:
							res[i] = nodeOR{n}
						default:
							return nil, fmt.Errorf("invalid token type: %s", tokenTypeString[tokenType])
						}
						res = append(res[:i-1], res[i:]...) // remove the (i-1)the element because it has become a sub node
					}
					break
				}
			}
		}
	}
	if len(res) != 1 {
		return nil, errors.New("parse tree must have exactly one start node")
	}
	startNode, ok := res[0].(node)
	if !ok {
		return nil, errors.New("start node is not a node")
	}
	return startNode, nil
}

func (n nodeAND) Condition() string {
	return fmt.Sprintf("(%s AND %s)", n.node1.Condition(), n.node2.Condition())
}

func (n nodeOR) Condition() string {
	return fmt.Sprintf("(%s OR %s)", n.node1.Condition(), n.node2.Condition())
}

func (n nodeVAL) Condition() string {
	if n.caseInsensitive {
		return strconv.Quote(n.nodeValue) + "i"
	}
	return strconv.Quote(n.nodeValue)
}

func (n nodeNOT) Condition() string {
	return fmt.Sprintf("NOT %s", n.node.Condition())
}
